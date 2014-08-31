package main

import (
	"encoding/json"
	"fmt"
	"github.com/elgs/gorest"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
)

func main() {
	input := args()
	config := parseConfig(input[0])
	if config == nil {
		return
	}
	ds := config["data_source"].(string)
	dbType := config["db_type"].(string)
	tokenTable := config["token_table"].(string)
	dbo := &gorest.MySqlDataOperator{
		Ds:         ds,
		DbType:     dbType,
		TokenTable: tokenTable,
	}
	r := &gorest.Gorest{
		EnableHttp: config["enable_http"].(bool),
		HostHttp:   config["host_http"].(string),
		PortHttp:   uint16(config["port_http"].(float64)),

		EnableHttps:   config["enable_https"].(bool),
		HostHttps:     config["host_https"].(string),
		PortHttps:     uint16(config["port_https"].(float64)),
		CertFileHttps: config["cert_file_https"].(string),
		KeyFileHttps:  config["key_file_https"].(string),

		UrlPrefix: config["url_prefix"].(string),
		Dbo:       dbo}
	r.Serve()
}

func parseConfig(configFile string) map[string]interface{} {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(configFile, "not found")
		return nil
	}
	var config map[string]interface{}
	if err := json.Unmarshal(b, &config); err != nil {
		fmt.Println("Error parsing", configFile)
		return nil
	}
	return config
}

func args() []string {
	ret := []string{}
	if len(os.Args) <= 1 {
		ret = append(ret, "gorest.json")
	} else {
		ret = os.Args[1:]
	}
	return ret
}
