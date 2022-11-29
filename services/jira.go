package components

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aadilpatel1989/qa-packages-golang/utils"

	"github.com/andygrunwald/go-jira"
)

var (
	tpJira     jira.BasicAuthTransport
	clientJira *jira.Client
	credjira   map[string]string
)

func init() {

	json.Unmarshal(utils.DecryptFileKMS(os.Getenv("GCP_KMS_NAME"), utils.ReadFileContent("../creds/jira.json.enc")), &credjira) // Initialise jira crdentials.

	var err error

	tpJira = jira.BasicAuthTransport{
		Username: credjira["EMAILID"],
		Password: credjira["API_KEY"],
	}
	clientJira, err = jira.NewClient(tpJira.Client(), utils.ConfigValuesJira["CODCOPS_URL"].(string))
	if err != nil {
		log.Println("Error:" + err.Error())
	}
}

//getProjectIdFromKey - get jira project ID from Key
func getProjectKeyFromId(projectId string) string {
	return strings.Replace(projectId, "-", "", 1)
}

// GetJiraProject - Get Project detail from project key.
func GetJiraProject(projectKey string) (*jira.Project, error) {

	project, _, err := clientJira.Project.Get(projectKey) //jira.Client.Project.Get()

	if err != nil {
		return nil, err
	}

	return project, err
}

//GetUserDetails - Get user details
func GetUserDetails(accountId string) *jira.User {
	user, _, err := clientJira.User.Get(accountId)

	if err != nil {
		log.Println("Error:" + err.Error())
	}
	//log.Printf("Name = %s, EmailAddress = %s", user.DisplayName, user.EmailAddress)
	return user
}

// GetAllIssues -  Get all issue of a jira project
func GetAllIssues(projectId string) []jira.Issue {
	//var key = getProjectKeyFromId(projectId)

	searchOptions := &jira.SearchOptions{
		MaxResults: 100,
	}

	issues, _, err := clientJira.Issue.Search(fmt.Sprintf("project=\"%s\"", projectId), searchOptions)

	if err != nil {
		log.Println("Error:" + err.Error())
	}

	return issues
}

// The method will take the project id and delete the respective jira project.
func DeleteJiraProject(url string, projectID string) {

	contextPath := "/project/" + projectID + "/delete"
	_, err := utils.MakeHttpPost(url+contextPath, projectID, credjira["EMAILID"], credjira["API_KEY"])

	if err != nil {
		log.Println("Jira project with " + projectID + " is not deleted due to " + err.Error())
	}
}

// The method will return the name of the current active sprint for the board specified in context url.
func GetActiveSprint() string {

	// var credjira map[string]string
	// json.Unmarshal(utils.ReadFileContent("../creds/jira.json"), &credjira)
	responseBytes, _ := utils.MakeHttpGet(utils.ConfigValuesJira["BASE_URL"].(string)+utils.ConfigValuesJira["CONTEXT_URL_ACTIVE_SPRINT"].(string), credjira["EMAILID"], credjira["API_KEY"])
	var response map[string]interface{}
	json.Unmarshal(responseBytes, &response)
	currentActiveSprintDetails := response["values"].([]interface{})
	currentActiveSprintName := currentActiveSprintDetails[0].(map[string]interface{})
	return currentActiveSprintName["name"].(string)
}

// The method will return the project lead for the given jira project
func GetProjectLeadForAProject(url string, projectID string) string {

	contextPathProject := "/project/" + projectID
	projectInfo, err := utils.MakeHttpGet(url+contextPathProject, credjira["EMAILID"], credjira["API_KEY"])
	if err != nil {
		log.Println("Error:" + err.Error())
	}

	var projectDetails map[string]interface{}

	json.Unmarshal(projectInfo, &projectDetails)

	leadInfo := projectDetails["lead"].(map[string]interface{})

	contextPathUser := "/user?accountId=" + leadInfo["accountId"].(string)

	userInfo, err := utils.MakeHttpGet(url+contextPathUser, credjira["EMAILID"], credjira["API_KEY"])

	if err != nil {
		log.Println("Error:" + err.Error())
	}

	var userDetails map[string]string

	json.Unmarshal(userInfo, &userDetails)

	return userDetails["emailAddress"]
}
