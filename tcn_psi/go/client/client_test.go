package client

import (
	"github.com/openmined/tcn-psi/server"
	"github.com/openmined/tcn-psi/tcn"
	"regexp"
	"testing"
)

func TestClientSanity(t *testing.T) {
	c, err := Create()
	if err != nil {
		t.Errorf("Failed to create a PSI client %v", err)
	}
	if c == nil {
		t.Errorf("Failed to create a PSI client: nil")
	}

	matched, _ := regexp.MatchString(`[0-9]+[.][0-9]+[.][0-9]+(-[A-Za-z0-9]+)?`, c.Version())
	if !matched {
		t.Errorf("Got invalid version %v", c.Version())
	}
}

func helperGetReports(cnt int) ([]*tcn.SignedReport, []tcn.TemporaryContactNumber, error) {
	tcnsPerReport := 10
	tcns := []tcn.TemporaryContactNumber{}
	reports := []*tcn.SignedReport{}

	for ridx := 0; ridx < cnt; ridx++ {
		rak, err := tcn.NewReportAuthorizationKey()
		if err != nil {
			return nil, nil, err
		}
		tck, err := rak.InitialTCK()
		if err != nil {
			return nil, nil, err
		}
		for idx := 0; idx < tcnsPerReport; idx++ {
			val, err := tck.TemporaryContactNumber()
			if err != nil {
				return nil, nil, err
			}
			tcns = append(tcns, *val)
			tck, err = tck.Ratchet()
			if err != nil {
				return nil, nil, err
			}
		}
		r, err := rak.CreateSignedReport(tcn.CoEpiV1Code, []byte{}, 1, uint16(tcnsPerReport))
		if err != nil {
			return nil, nil, err
		}
		if ridx%2 == 0 {
			reports = append(reports, r)
		}
	}
	return reports, tcns, nil
}
func TestClientFailure(t *testing.T) {
	c := &TCNClient{}
	_, clientItems, err := helperGetReports(10)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = c.CreateRequest(clientItems)
	if err == nil {
		t.Errorf("CreateRequest with an invalid context should fail")
	}
	_, err = c.GetIntersectionSize("dummy1", "dummy2")
	if err == nil {
		t.Errorf("GetIntersectionSize with an invalid context should fail")
	}
	c, _ = Create()
	_, err = c.GetIntersectionSize("dummy1", "dummy2")
	if err == nil {
		t.Errorf("GetIntersectionSize with invalid input should fail")
	}
}

func TestClientServer(t *testing.T) {
	client, err := Create()
	if err != nil || client == nil {
		t.Errorf("Failed to create a PSI client %v", err)
	}

	server, err := server.CreateWithNewKey()
	if err != nil || server == nil {
		t.Errorf("Failed to create a PSI server %v", err)
	}

	serverItems, clientItems, err := helperGetReports(1000)
	if err != nil {
		t.Error(err.Error())
	}
	cntClientItems := len(clientItems)

	setup, err := server.CreateSetupMessage(0.01, int64(cntClientItems), serverItems)
	if err != nil {
		t.Errorf("failed to create setup msg %v", err)
	}
	request, err := client.CreateRequest(clientItems)
	if err != nil {
		t.Errorf("failed to create request %v", err)
	}
	serverResp, err := server.ProcessRequest(request)
	if err != nil {
		t.Errorf("failed to process request %v", err)
	}
	intersectionCnt, err := client.GetIntersectionSize(setup, serverResp)
	if err != nil {
		t.Errorf("failed to compute intersection %v", err)
	}

	if int(intersectionCnt) < (cntClientItems / 2) {
		t.Errorf("Invalid intersection. expected lower bound %v. got %v", (cntClientItems / 2), intersectionCnt)
	}

	if float64(intersectionCnt) > float64(cntClientItems/2)*float64(1.1) {
		t.Errorf("Invalid intersection. expected upper bound %v. got %v", float64(cntClientItems/2)*float64(1.1), intersectionCnt)
	}
}

var result string

func benchmarkClientCreateRequest(cnt int, b *testing.B) {
	b.ReportAllocs()
	total := 0
	for n := 0; n < b.N; n++ {
		client, err := Create()
		if err != nil || client == nil {
			b.Errorf("failed to get client")
		}
		_, inputs, err := helperGetReports(cnt)
		if err != nil {
			b.Error(err.Error())
		}
		request, err := client.CreateRequest(inputs)
		if err != nil {
			b.Errorf("failed to generate request")
		}

		//ugly hack for preventing compiler optimizations
		result = request

		total += cnt
		b.ReportMetric(float64(len(request)), "RequestSize")

	}
	b.ReportMetric(float64(total), "ElementsProcessed")
}

func BenchmarkClientCreateRequest1(b *testing.B)     { benchmarkClientCreateRequest(1, b) }
func BenchmarkClientCreateRequest10(b *testing.B)    { benchmarkClientCreateRequest(10, b) }
func BenchmarkClientCreateRequest100(b *testing.B)   { benchmarkClientCreateRequest(100, b) }
func BenchmarkClientCreateRequest1000(b *testing.B)  { benchmarkClientCreateRequest(1000, b) }
func BenchmarkClientCreateRequest10000(b *testing.B) { benchmarkClientCreateRequest(10000, b) }

var dummyInt64 int64

func benchmarkClientGetIntersectionSize(cnt int, b *testing.B) {
	b.ReportAllocs()
	total := 0
	for n := 0; n < b.N; n++ {
		client, err := Create()
		if err != nil || client == nil {
			b.Errorf("failed to get client")
		}
		server, err := server.CreateWithNewKey()
		if err != nil || server == nil {
			b.Errorf("failed to get server")
		}

		fpr := 1. / (1000000)

		serverItems, clientItems, err := helperGetReports(cnt)
		if err != nil {
			b.Error(err.Error())
		}
		setup, err := server.CreateSetupMessage(fpr, int64(cnt), serverItems)
		if err != nil {
			b.Errorf("failed to create setup msg %v", err)
		}
		request, err := client.CreateRequest(clientItems)
		if err != nil {
			b.Errorf("failed to create request %v", err)
		}
		serverResp, err := server.ProcessRequest(request)
		if err != nil {
			b.Errorf("failed to process request %v", err)
		}
		intersectionCnt, err := client.GetIntersectionSize(setup, serverResp)
		if err != nil {
			b.Errorf("failed to process response %v", err)
		}
		total += cnt
		//ugly hack for preventing compiler optimizations
		dummyInt64 = intersectionCnt
	}
	b.ReportMetric(float64(total), "ElementsProcessed")
}

func BenchmarkClientGetIntersectionSize1(b *testing.B)    { benchmarkClientGetIntersectionSize(1, b) }
func BenchmarkClientGetIntersectionSize10(b *testing.B)   { benchmarkClientGetIntersectionSize(10, b) }
func BenchmarkClientGetIntersectionSize100(b *testing.B)  { benchmarkClientGetIntersectionSize(100, b) }
func BenchmarkClientGetIntersectionSize1000(b *testing.B) { benchmarkClientGetIntersectionSize(1000, b) }
func BenchmarkClientGetIntersectionSize10000(b *testing.B) {
	benchmarkClientGetIntersectionSize(10000, b)
}
