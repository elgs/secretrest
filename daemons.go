// daemons
package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/elgs/gorest"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func init() {
	gorest.RegisterJob("check_last", &gorest.Job{
		Cron: "0/5 * * * * *",
		MakeAction: func(dbo gorest.DataOperator) func() {
			return func() {
				command := fmt.Sprint("last -F | grep ppp")
				output, err := exec.Command("bash", "-c", command).CombinedOutput()
				if err != nil {
					fmt.Println("Failed to execute:", err, command)
					fmt.Println(string(output))
				} else {
					fmt.Println(string(output))
				}
			}
		},
	})
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
