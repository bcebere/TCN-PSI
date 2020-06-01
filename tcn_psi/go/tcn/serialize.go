package tcn

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
)

// GetReport inteprets data as a report and returns it as a parsed structure.
func GetReport(data []byte) (*Report, uint16, error) {
	if len(data) < ReportMinLength {
		return nil, 0, errors.New("Data too short to be a valid report")
	}
	tckBytes := [32]byte{}
	copy(tckBytes[:], data[32:64])

	memoDataLen := uint8(data[69])

	return &Report{
		RVK:      ed25519.PublicKey(data[:32]),
		TCKBytes: tckBytes,
		J1:       binary.LittleEndian.Uint16(data[64:66]),
		J2:       binary.LittleEndian.Uint16(data[66:68]),
		MemoType: data[68],
		MemoData: data[70 : 70+memoDataLen],
	}, 70 + uint16(memoDataLen), nil
}

// GetSignedReport interprets data as a signed report and returns it as a
// parsed structure.
func GetSignedReport(data []byte) (*SignedReport, error) {
	if len(data) < SignedReportMinLength {
		return nil, errors.New("Data too short to be a valid signed report")
	}

	report, reportEndPos, err := GetReport(data)
	if err != nil {
		return nil, err
	}
	endPos := reportEndPos + ed25519.SignatureSize
	if int(endPos) > len(data) {
		return nil, errors.New("invalid signed report length")
	}
	sig := data[reportEndPos:endPos]

	return &SignedReport{
		Report: report,
		Sig:    sig,
	}, nil
}
