package tcn

import (
	"crypto/ed25519"
	"encoding/binary"
)

// Describes the intended type of the contents of a memo field.
const (
	// The CoEpi symptom self-report format, version 1
	CoEpiV1Code = 0x0
	// The CovidWatch test data format, version 1
	CovidWatchV1Code = 0x1
	// ITOMemoCode is the code that marks a report as an ito report in the
	// memo.
	ITOMemoCode = 0x2
	// ReportMinLength is the minimum length of a TCN report (with memo data
	// of length 0) in bytes.
	ReportMinLength = 70
)

// Report represents a report as described in the TCN protocol:
// https://github.com/TCNCoalition/TCN#reporting
type Report struct {
	RVK      ed25519.PublicKey
	TCKBytes [32]byte
	J1       uint16
	J2       uint16
	MemoType uint8
	MemoData []uint8
}

//TemporaryContactNumbers returns a slice over all temporary contact numbers included in the report.
func (r Report) TemporaryContactNumbers() (map[uint16]TemporaryContactNumber, error) {
	tck0 := TemporaryContactKey{
		Index:    r.J1 - 1,
		RVK:      r.RVK,
		TCKBytes: r.TCKBytes,
	}
	//generate tck_{J1}
	tck, err := tck0.Ratchet()
	if err != nil {
		return nil, err
	}

	result := map[uint16]TemporaryContactNumber{}
	for idx := r.J1; idx < r.J2; idx++ {
		val, err := tck.TemporaryContactNumber()
		if err != nil {
			return nil, err
		}
		result[idx] = *val
		tck, err = tck.Ratchet()
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Bytes converts r to a concatenated byte array represention.
func (r *Report) Bytes() ([]byte, error) {
	var data []byte
	data = append(data, r.RVK...)
	data = append(data, r.TCKBytes[:]...)

	j1Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(j1Bytes, r.J1)
	j2Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(j2Bytes, r.J2)
	data = append(data, j1Bytes...)
	data = append(data, j2Bytes...)

	// Memo
	data = append(data, r.MemoType)
	data = append(data, uint8(len(r.MemoData)))
	data = append(data, r.MemoData...)

	return data, nil
}
