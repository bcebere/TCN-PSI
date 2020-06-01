package tcn

import (
	"crypto/ed25519"
	"encoding/binary"
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
