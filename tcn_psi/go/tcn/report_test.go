package tcn_test

import (
	"testing"

	"github.com/juliangruber/go-intersect"
	"github.com/openmined/tcn-psi/tcn"
)

func TestRecomputeTCN(t *testing.T) {
	rak, err := tcn.NewReportAuthorizationKey()
	if err != nil {
		t.Error(err.Error())
	}
	tck, err := rak.InitialTCK()
	if err != nil {
		t.Error(err.Error())
	}

	tcns := []tcn.TemporaryContactNumber{}
	for idx := 0; idx < 100; idx++ {
		val, err := tck.TemporaryContactNumber()
		if err != nil {
			t.Error(err.Error())
		}
		tcns = append(tcns, *val)
		tck, err = tck.Ratchet()
		if err != nil {
			t.Error(err.Error())
		}
	}

	signedReport, err := rak.CreateSignedReport(tcn.CoEpiV1Code, []byte("symptom data"), 20, 90)
	if err != nil {
		t.Error(err.Error())
	}

	signed, err := signedReport.Verify()
	if err != nil || !signed {
		t.Error("invalid signature")
	}

	report := signedReport.Report

	recomputedTCNS, err := report.TemporaryContactNumbers()
	if err != nil {
		t.Error(err.Error())
	}

	for idx := range recomputedTCNS {
		if recomputedTCNS[idx] != tcns[idx-1] {
			t.Errorf("invalid TCN at index %v", idx)
		}
	}
}

func TestMultipleMatches(t *testing.T) {
	// Parameters.
	numReports := 10000
	tcnsPerReport := uint16(24 * 60 / 15)

	tcns := []tcn.TemporaryContactNumber{}
	reports := []*tcn.SignedReport{}

	// Generate some tcns that will be reported.
	for ridx := 0; ridx < numReports; ridx++ {
		rak, err := tcn.NewReportAuthorizationKey()
		if err != nil {
			t.Error(err.Error())
		}
		tck, err := rak.InitialTCK()
		if err != nil {
			t.Error(err.Error())
		}

		for idx := uint16(0); idx < tcnsPerReport; idx++ {
			val, err := tck.TemporaryContactNumber()
			if err != nil {
				t.Error(err.Error())
			}
			tcns = append(tcns, *val)
			tck, err = tck.Ratchet()
			if err != nil {
				t.Error(err.Error())
			}
		}
		r, err := rak.CreateSignedReport(tcn.CoEpiV1Code, []byte{}, 1, tcnsPerReport)
		if err != nil {
			t.Error(err)
		}
		reports = append(reports, r)
	}

	expectedTCNs := make([]tcn.TemporaryContactNumber, len(tcns))
	copy(expectedTCNs, tcns)

	candidateTCNs := []tcn.TemporaryContactNumber{}

	// Generate some extra tcns that will not be reported.
	rak, err := tcn.NewReportAuthorizationKey()
	if err != nil {
		t.Error(err.Error())
	}
	tck, err := rak.InitialTCK()
	if err != nil {
		t.Error(err.Error())
	}

	for idx := 0; idx < 60000; idx++ {
		val, err := tck.TemporaryContactNumber()
		if err != nil {
			t.Error(err.Error())
		}
		tcns = append(tcns, *val)
		tck, err = tck.Ratchet()
		if err != nil {
			t.Error(err.Error())
		}
	}

	// Now expand the reports into an array of candidates
	for _, report := range reports {
		if signed, err := report.Verify(); err != nil || !signed {
			t.Errorf("invalid report")
		}
		candidates, err := report.Report.TemporaryContactNumbers()
		if err != nil {
			t.Error(err.Error())
		}
		for _, v := range candidates {
			candidateTCNs = append(candidateTCNs, v)
		}
	}

	reportedTCNS := intersect.Hash(candidateTCNs, tcns)
	intersection := intersect.Hash(reportedTCNS, expectedTCNs)

	if len(intersection) != len(expectedTCNs) {
		t.Errorf("invalid number of reported TCNs: expected %v got %v", len(expectedTCNs), len(intersection))
	}

}
