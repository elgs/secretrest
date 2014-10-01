package main

import (
	"encoding/json"
	"fmt"
	"github.com/elgs/gorest"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
	fbp := config["file_base_path"]
	u, _ := user.Current()
	fileBasePath, _ := filepath.Abs(u.HomeDir + string(os.PathSeparator) + "files")
	if fbp != nil {
		if !strings.HasPrefix(fbp.(string), string(os.PathSeparator)) {
			fileBasePath, _ = filepath.Abs(u.HomeDir + string(os.PathSeparator) + fbp.(string))
		} else {
			fileBasePath, _ = filepath.Abs(fbp.(string))
		}
	}
	dbo := &gorest.MySqlDataOperator{
		Ds:         ds,
		DbType:     dbType,
		TokenTable: tokenTable,
	}

	gorest.RegisterDataOperator("api", dbo)

	r := &gorest.Gorest{
		EnableHttp: config["enable_http"].(bool),
		HostHttp:   config["host_http"].(string),
		PortHttp:   uint16(config["port_http"].(float64)),

		EnableHttps:   config["enable_https"].(bool),
		HostHttps:     config["host_https"].(string),
		PortHttps:     uint16(config["port_https"].(float64)),
		CertFileHttps: config["cert_file_https"].(string),
		KeyFileHttps:  config["key_file_https"].(string),

		FileBasePath: fileBasePath,
	}
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
