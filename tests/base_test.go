package tests

import (
	"os"
	"qa-packages-golang/reports"
	"qa-packages-golang/utils"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	startTime := time.Now()
	os.Remove(utils.GetValueFromConfigFile().ProjectPath + "/reports/resources/temp_report.html")
	reports.StartTime = startTime.Format(time.RFC822)
	m.Run()
	reports.TotalExecutionTime = time.Since(startTime)
	reports.TestReportCreate()
}
