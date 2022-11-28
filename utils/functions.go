package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"qa-packages-golang/structures"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/camelcase"
	"golang.org/x/net/websocket"
)

// The method takes the url/payload and do the post request.
func MakeHttpPost(url string, request interface{}, username string, password string) ([]byte, error) {

	jsonBody, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	if username == "" && password != "" {
		req.Header.Add("Authorization", "Bearer "+password)
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return body, err
}

// The function will perform a http get request and will return the response bytes.
func MakeHttpGet(url string, username string, password string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error:" + err.Error())
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	if username == "" && password != "" {
		req.Header.Add("Authorization", "Bearer "+password)
	}

	response, err := client.Do(req)

	if err != nil {
		log.Println("Error:" + err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error:" + err.Error())
	}

	return body, err
}

// The function will perform a http put request and will return the response bytes.
func MakeHttpPut(url string, username string, password string, contentType string, body io.Reader, queryParams map[string]string) ([]byte, error) {

	var (
		req *http.Request
		err error
	)

	client := &http.Client{}

	if body == nil {
		req, err = http.NewRequest("PUT", url, nil)
	} else {
		req, err = http.NewRequest("PUT", url, body)
	}

	//log.Println("PUT request for the url " + url + " is initiated.")

	req.Header.Add("Content-Type", contentType)

	queryParam := req.URL.Query()

	for key, value := range queryParams {
		queryParam.Add(key, value)
	}

	req.URL.RawQuery = queryParam.Encode()

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	response, err := client.Do(req)
	//log.Println("Status code for the PUT request is " + strconv.Itoa(response.StatusCode))

	responseBody, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	return responseBody, err
}

// The method takes in the input and converts it to string type.
func MarshalToString(input interface{}) string {
	result, err := json.Marshal(input)
	if err != nil {
		log.Printf("%v", err)
	}
	return string(result)
}

// The method will check whether the array contains the item.
func IfArrayContains(value interface{}, array []interface{}) bool {
	for _, arrayItem := range array {
		if arrayItem == value {
			return true
		}
	}
	return false
}

// The method reads the content of file and return it in byte array.
func ReadFileContent(pathToFile string) []byte {
	bytes, _ := ioutil.ReadFile(pathToFile)
	return bytes
}

// The function will perform a http get request and will return the response object.
func GetHttpResponse(url string, username string, password string) (*http.Response, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error:" + err.Error())
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	if username == "" && password != "" {
		req.Header.Add("Authorization", "Bearer "+password)
	}

	response, err := client.Do(req)

	if err != nil {
		log.Println("Error:" + err.Error())
	}

	defer response.Body.Close()

	return response, err
}

func SwitchRole(accountID, roleName, region string) (*session.Session, error) {

	baseSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	stsSvc := sts.New(baseSess)
	sessionName := "tenant-session"
	assumedRole, err := stsSvc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::" + accountID + ":role/" + roleName),
		RoleSessionName: aws.String(sessionName),
	})

	if err != nil {
		return nil, err
	}

	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			*assumedRole.Credentials.AccessKeyId,
			*assumedRole.Credentials.SecretAccessKey,
			*assumedRole.Credentials.SessionToken),
		Region: aws.String(region),
	})
}

func GetValueFromConfigFile() structures.Config {

	var cfg structures.Config

	json.Unmarshal(ReadFileContent("../config.json"), &cfg)

	return cfg
}

func ConvertCamelcaseIntoNormalSentence(camelCase string) string {
	strs := strings.Split(camelCase, "/Test")

	strArray := camelcase.Split(strs[1])

	newStr := strArray[0]

	for index := 1; index < len(strArray); index++ {
		newStr = newStr + " " + strings.ToLower(strArray[index])
	}

	return newStr + "."
}

func WebSocketSession(url, origin string) (*websocket.Conn, error) {

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println(err)
	}
	return ws, err
}

func ReadMessageOnWebSocketConnection(conn *websocket.Conn) ([]byte, int, error) {
	var (
		msg = make([]byte, 512)
		n   int
		err error
	)
	if n, err = conn.Read(msg); err != nil {
		log.Println(err)
		return nil, 0, err
	}

	return msg, n, nil
}

func HttpPutRequestObject(urlPath, imagePath string) (*http.Request, error) {

	req, err := http.NewRequest("PUT", urlPath, bytes.NewReader(ReadFileContent(imagePath)))
	if err != nil {
		return nil, err
	}

	return req, nil
}
