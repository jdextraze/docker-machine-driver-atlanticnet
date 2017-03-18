package atlanticnet

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	DescribePlan(name, platform string) ([]Plan, error)

	ListInstances() ([]Instance, error)
	RunInstance(req RunInstanceRequest) ([]RunInstance, error)
	TerminateInstance(instanceId string) ([]TerminateInstance, error)
	DescribeInstance(instanceId string) (*InstanceDescription, error)
	RebootInstance(instanceId string, rebootType RebootType) (*RebootInstance, error)

	DescribeImage(imageId string) ([]Image, error)

	ListSshKeys() ([]SshKey, error)
}

type client struct {
	apiKey string
	secret string
	debug  bool
}

func NewClient(apiKey, secret string, debug bool) Client {
	return &client{apiKey, secret, debug}
}

func (c *client) doRequest(action string, in map[string]string, out interface{}) error {
	form := url.Values{}
	form.Add("Format", "json")
	form.Add("Version", "2010-12-30")
	form.Add("ACSAccessKeyId", c.apiKey)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	form.Add("Timestamp", timestamp)
	randomGuid := uuid.NewV4().String()
	form.Add("Rndguid", randomGuid)
	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(timestamp + randomGuid))
	signature := mac.Sum(nil)
	form.Add("Signature", string(base64.StdEncoding.EncodeToString(signature)))

	if in != nil {
		for key, value := range in {
			form.Add(key, value)
		}
	}

	encodedForm := form.Encode()
	if c.debug {
		log.Println("===>", encodedForm)
	}

	req, err := http.NewRequest(
		"POST",
		"https://cloudapi.atlantic.net/?Action="+action,
		strings.NewReader(encodedForm),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if c.debug {
		log.Println("<===", string(content))
	}

	return json.Unmarshal(content, &out)
}

type ErrorResponse struct {
	Code    string
	Message string
	Time    int
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

type BaseResponse struct {
	Error     *ErrorResponse
	Timestamp int
}
