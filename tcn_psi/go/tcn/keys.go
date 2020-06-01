package tcn

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// HTCKDomainSep is the domain separator used for the domain-separated hash
// function.
var HTCKDomainSep = []byte("H_TCK")

// HTCNDomainSep
var HTCNDomainSep = []byte("H_TCN")

// TemporaryContactNumber is a pseudorandom 128-bit value broadcast to nearby
// devices over Bluetooth
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

type ReportAuthorizationKey struct {
	RAK ed25519.PrivateKey
	RVK ed25519.PublicKey
}

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

func (r *ReportAuthorizationKey) InitialTCK() (*TemporaryContactKey, error) {
	tck0, err := r.TCK0()
	if err != nil {
		return nil, err
	}
	return tck0.Ratchet()
}

func (r *ReportAuthorizationKey) TCK0() (*TemporaryContactKey, error) {
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

func (r *ReportAuthorizationKey) CreateReport(memoType uint8, memoData []uint8, j1, j2 uint16) (*Report, error) {
	if j1 == 0 {
		j1 = 1
	}
	tck, err := r.TCK0()
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

func (r *ReportAuthorizationKey) CreateSignedReport(memoType uint8, memoData []uint8, j1, j2 uint16) (*SignedReport, error) {
	report, err := r.CreateReport(memoType, memoData, j1, j2)
	if err != nil {
		return nil, err
	}
	signed, err := GenerateSignedReport(&r.RAK, report)

	return signed, err
}
