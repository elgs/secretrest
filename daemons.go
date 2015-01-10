// daemons
package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func init() {
	//gorest.RegisterJob("check_last", &gorest.Job{
	//	Cron: "0/5 * * * * *",
	//	MakeAction: func(dbo gorest.DataOperator) func() {
	//		return func() {
	//		}
	//	},
	//})
}

func httpRequest(url string, method string, data string, apiTokenId string, apiTokenKey string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	if apiTokenId != "" {
		req.Header.Add("api_token_id", apiTokenId)
	}
	if apiTokenKey != "" {
		req.Header.Add("api_token_key", apiTokenKey)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	defer tr.CloseIdleConnections()

	return body, err
}
