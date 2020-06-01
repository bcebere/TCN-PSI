package tcn

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// HTCKDomainSep is the domain separator used for the TCK domain-separated hash
// function.
var HTCKDomainSep = []byte("H_TCK")

// HTCNDomainSep is the domain separator for the TCN domain-separated hash function.
var HTCNDomainSep = []byte("H_TCN")

// TemporaryContactNumber is a pseudorandom 128-bit value broadcast to nearby
// devices over Bluetooth.
type TemporaryContactNumber [16]uint8

// TemporaryContactKey is a ratcheting key used to derive temporary contact
// numbers.
type TemporaryContactKey struct {
	Index    uint16
	RVK      ed25519.PublicKey
	TCKBytes [32]byte
}

// Ratchet the key forward, producing a new key for a new temporary
// contact number.
func (tck *TemporaryContactKey) Ratchet() (*TemporaryContactKey, error) {
	nextHash := sha256.New()
	if _, err := nextHash.Write([]byte(HTCKDomainSep)); err != nil {
		fmt.Printf("Failed to write tck domain separator: %s\n", err.Error())
		return nil, err
	}
	if _, err := nextHash.Write(tck.RVK); err != nil {
		fmt.Printf("Failed to write rvk: %s\n", err.Error())
		return nil, err
	}
	if _, err := nextHash.Write(tck.TCKBytes[:]); err != nil {
		fmt.Printf("Failed to write tck bytes: %s\n", err.Error())
		return nil, err
	}

	if tck.Index == math.MaxUint16 {
		return nil, errors.New("rak should be rotated")
	}

	newTCKBytes := [32]byte{}
	copy(newTCKBytes[:32], nextHash.Sum(nil))

	return &TemporaryContactKey{
		Index:    tck.Index + 1,
		RVK:      tck.RVK,
		TCKBytes: newTCKBytes,
	}, nil
}

//TemporaryContactNumber computes the temporary contact number derived from this key.
func (tck *TemporaryContactKey) TemporaryContactNumber() (*TemporaryContactNumber, error) {
	nextHash := sha256.New()
	if _, err := nextHash.Write(HTCNDomainSep); err != nil {
		fmt.Printf("Failed to write TCN domain separator: %s\n", err.Error())
		return nil, err
	}
	indexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(indexBytes, tck.Index)
	if _, err := nextHash.Write(indexBytes); err != nil {
		fmt.Printf("Failed to write index bytes: %s\n", err.Error())
		return nil, err
	}
	if _, err := nextHash.Write(tck.TCKBytes[:]); err != nil {
		fmt.Printf("Failed to write tck bytes: %s\n", err.Error())
		return nil, err
	}

	tcnBytes := [16]byte{}
	copy(tcnBytes[:16], nextHash.Sum(nil))

	result := TemporaryContactNumber(tcnBytes)
	return &result, nil
}

//ReportAuthorizationKey authorizes publication of a report of potential exposure.
type ReportAuthorizationKey struct {
	RAK ed25519.PrivateKey
	RVK ed25519.PublicKey
}

//NewReportAuthorizationKey initialize a new report authorization key from a random number generator.
func NewReportAuthorizationKey() (*ReportAuthorizationKey, error) {
	rvk, rak, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &ReportAuthorizationKey{
		RAK: rak,
		RVK: rvk,
	}, nil
}

//InitialTCK computes the initial temporary contact key derived from this report authorization key.
//Note: this function returns `tck_1`, the first temporary contact key that can be
//used to generate tcks.
func (r *ReportAuthorizationKey) InitialTCK() (*TemporaryContactKey, error) {
	tck0, err := r.tck0()
	if err != nil {
		return nil, err
	}
	return tck0.Ratchet()
}

func (r *ReportAuthorizationKey) tck0() (*TemporaryContactKey, error) {
	tck0Hash := sha256.New()
	if _, err := tck0Hash.Write([]byte(HTCKDomainSep)); err != nil {
		fmt.Printf("Failed to write tck domain separator: %s\n", err.Error())
		return nil, err
	}
	if _, err := tck0Hash.Write(r.RAK); err != nil {
		fmt.Printf("Failed to write rak: %s\n", err.Error())
		return nil, err
	}

	tck0Bytes := [32]byte{}
	copy(tck0Bytes[:32], tck0Hash.Sum(nil))

	return &TemporaryContactKey{
		Index:    0,
		RVK:      r.RVK,
		TCKBytes: tck0Bytes,
	}, nil
}

//CreateReport creates a report of potential exposure.
//
//# Inputs
//
//- `memoType`, `memoData`: the type and data for the report's memo field.
//- `j_1 > 0`: the ratchet index of the first temporary contact number in the report.
//- `j_2`: the ratchet index of the last temporary contact number other users should check.
//
//# Notes
//
//Creating a report reveals *all* temporary contact numbers subsequent to
//`j_1`, not just up to `j_2`, which is included for convenience.
//
//The `memo_data` must be less than 256 bytes long.
//
//Reports are unlinkable from each other **only up to the memo field**. In
//other words, adding the same high-entropy data to the memo fields of
//multiple reports will cause them to be linkable.
func (r *ReportAuthorizationKey) CreateReport(memoType uint8, memoData []uint8, j1, j2 uint16) (*Report, error) {
	if j1 == 0 {
		j1 = 1
	}
	tck, err := r.tck0()
	if err != nil {
		return nil, err
	}
	for idx := uint16(0); idx < j1-1; idx++ {
		tck, err = tck.Ratchet()
		if err != nil {
			return nil, err
		}
	}

	return &Report{
		RVK:      r.RVK,
		TCKBytes: tck.TCKBytes,
		J1:       j1,
		J2:       j2,
		MemoType: memoType,
		MemoData: memoData,
	}, nil
}

//CreateSignedReport creates a signed exposure report, whose source integrity can be verified to produce a `Report`.
func (r *ReportAuthorizationKey) CreateSignedReport(memoType uint8, memoData []uint8, j1, j2 uint16) (*SignedReport, error) {
	report, err := r.CreateReport(memoType, memoData, j1, j2)
	if err != nil {
		return nil, err
	}
	signed, err := GenerateSignedReport(&r.RAK, report)

	return signed, err
}
