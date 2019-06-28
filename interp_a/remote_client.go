package interp_a

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type RemoteEvaluator struct {
	connParams ConnectionParams
	client     http.Client
}

type ConnectionParams struct {
	Address string
}

func (evaluator *RemoteEvaluator) Init(
	connParams ConnectionParams,
) {
	evaluator.connParams = connParams
	evaluator.client = http.Client{}
}

func (evaluator RemoteEvaluator) OpEvaluate(args []interface{}) ([]interface{}, error) {
	list, err := json.Marshal(args)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Create form data
	form := url.Values{}
	form.Add("list", string(list))

	// Create request object (and encode form to do so)
	req, err := http.NewRequest(
		"POST", evaluator.connParams.Address+"/call",
		strings.NewReader(
			form.Encode(),
		),
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Perform request
	resp, err := evaluator.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Read response
	// -- read response as bytes
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(responseBytes))
	}
	// -- parse response bytes as JSON
	responseList := []interface{}{}
	err = json.Unmarshal(responseBytes, &responseList)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return responseList, nil
}
