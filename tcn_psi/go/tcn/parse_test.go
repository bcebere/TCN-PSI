package tcn_test

import (
	"testing"

	"github.com/openmined/tcn-psi/tcn"
	"github.com/stretchr/testify/assert"
)

func TestGetReport(t *testing.T) {
	rak, err := tcn.NewReportAuthorizationKey()
	if err != nil {
		t.Error(err.Error())
	}
	report, err := rak.CreateReport(tcn.CoEpiV1Code, []byte("symptom data"), 0, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	rb, err := report.Bytes()
	if err != nil {
		t.Error(err.Error())
		return
	}

	retReport, _, err := tcn.GetReport(rb)
	assert.NoError(t, err)
	assert.EqualValues(t, report, retReport)
}

func TestGetSignedReport(t *testing.T) {
	rak, err := tcn.NewReportAuthorizationKey()
	if err != nil {
		t.Error(err.Error())
	}
	signedReport, err := rak.CreateSignedReport(tcn.CoEpiV1Code, []byte("symptom data"), 0, 4)
	if err != nil {
		t.Error(err.Error())
		return
	}
	srb, err := signedReport.Bytes()
	if err != nil {
		t.Error(err.Error())
		return
	}

	retSignedReport, err := tcn.GetSignedReport(srb)

	assert.NoError(t, err)
	assert.EqualValues(t, signedReport, retSignedReport)
}
