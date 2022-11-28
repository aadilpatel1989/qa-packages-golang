package reports

import (
	"io/ioutil"
	"log"
	"os"
	"qa-packages-golang/utils"
	"strconv"
	"time"
)

var (
	srNo, totalPass, totalFail int = 0, 0, 0
	currentSprintName              = "STEC11 QA Report" //services.GetActiveSprint()
	TestReportName             *os.File
	TotalExecutionTime         time.Duration
	StartTime                  string
)

func reportFirstPart(projectName string, time string, passed int, failed int, skipped int, blocked int, total int) string {
	log.Println(StartTime)
	//return "<html><body><style>table, th, td {border: 1px solid black;border-collapse: collapse;}</style><center><h1>" + projectName + "</h1><h2>Execution time: " + time[0:5] + " seconds</h2><table><tr><th>Passed</th><th>Failed</th><th>Total</th></tr><tr><td>" + strconv.Itoa(passed) + "</td><td>" + strconv.Itoa(failed) + "</td><td>" + strconv.Itoa(total) + "</td></tr></table><table width=\"100%\"><tr><th>Sr No</th><th>Testcase</th><th>Input</th><th>Output</th><th>Status</th></tr>"
	return "<html><body><style>table, th, td {border: 1px solid black;border-collapse: collapse;}</style><center><h1>" + projectName + "</h1><h2>Start time: " + StartTime + "</h2><h2>Execution time: " + time + "</h2><table><tr><th>Passed</th><th>Failed</th><th>Total</th></tr><tr><td>" + strconv.Itoa(passed) + "</td><td>" + strconv.Itoa(failed) + "</td><td>" + strconv.Itoa(total) + "</td></tr></table><table width=\"100%\"><tr><th>Sr No</th><th>Testcase</th><th>Input</th><th>Output</th><th>Status</th></tr>"

}

// The method will create the the intermediate report.
func TestExecutionUpdate(srNo int, testcase string, input string, response string, status string) {

	file, err := os.OpenFile(utils.GetValueFromConfigFile().ProjectPath+"/reports/resources/temp_report.html", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	if _, err := file.WriteString("<tr><td>" + strconv.Itoa(srNo) + "</td><td>" + testcase + "</td><td>" + input + "</td><td>" + response + "</td><td>" + status + "</td></tr>"); err != nil {
		log.Println(err)
	}

}

// The method will create the test report.
func TestReportCreate() {

	TestReportName, _ = os.OpenFile(utils.GetValueFromConfigFile().ProjectPath+"/reports/resources/"+currentSprintName+" - "+time.Now().Format(time.RFC822)+".html", os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0644)

	defer TestReportName.Close()

	bytes, _ := ioutil.ReadFile(utils.GetValueFromConfigFile().ProjectPath + "/reports/resources/temp_report.html")

	if _, err := TestReportName.WriteString(reportFirstPart(currentSprintName, TotalExecutionTime.String(), totalPass, totalFail, 10, 10, totalPass+totalFail) + string(bytes) + "</table></center></body></html>"); err != nil {
		log.Println(err)
	}
}

// The method will mark the status of each testcase result
func MarkTestStatus(testStatus string, testName string, input string, output string, err error) {
	srNo++
	if testStatus == "Pass" {
		totalPass++
		TestExecutionUpdate(srNo, testName, "Not required", "Not required", "Passed")
	} else if testStatus == "Fail" && err != nil {
		totalFail++
		TestExecutionUpdate(srNo, testName, input, err.Error(), "Failed")
	} else {
		totalFail++
		TestExecutionUpdate(srNo, testName, input, output, "Failed")
	}
}
