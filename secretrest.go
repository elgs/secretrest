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
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	fileBasePath, err := filepath.Abs(u.HomeDir + string(os.PathSeparator) + "files")
	if err != nil {
		fmt.Println(err)
	}
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

	r := &gorest.Gorest{
		FileBasePath: fileBasePath,
	}

	if v, ok := config["enable_http"].(bool); ok {
		r.EnableHttp = v
	}
	if v, ok := config["host_http"].(string); ok {
		r.HostHttp = v
	}
	if v, ok := config["port_http"].(float64); ok {
		r.PortHttp = uint16(v)
	}
	if v, ok := config["enable_https"].(bool); ok {
		r.EnableHttps = v
	}
	if v, ok := config["host_https"].(string); ok {
		r.HostHttps = v
	}
	if v, ok := config["port_https"].(float64); ok {
		r.PortHttps = uint16(v)
	}
	if v, ok := config["cert_file_https"].(string); ok {
		r.CertFileHttps = v
	}
	if v, ok := config["cert_root_https"].(string); ok {
		r.CertRootHttps = v
	}
	if v, ok := config["key_file_https"].(string); ok {
		r.KeyFileHttps = v
	}
	if v, ok := config["session_key"].(string); ok {
		r.SessionKey = v
	}

	gorest.RegisterDataOperator("api", dbo)
	gorest.RegisterHttpHandlers(dbo, r)
	gorest.StartDaemons(dbo)

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
