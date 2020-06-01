package tcn

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
