package main

import (
	"bufio"
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgs/gorest"
	"github.com/elgs/gosqljson"
	"os"
	"strings"
	"time"
)

func init() {
	tableId := "secret"
	gorest.RegisterDataInterceptor(tableId, &SecretDataInterceptor{Id: "secret"})
}

type SecretDataInterceptor struct {
	*gorest.DefaultDataInterceptor
	Id string
}

var filePath string = "/etc/ppp/chap-secrets"

//var filePath string = "/Users/elgs/Desktop/chap-secrets"

var header string = `# Secrets for authentication using CHAP
# client	server	secret			IP addresses
`

func (this *SecretDataInterceptor) BeforeCreate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	userToken := context["user_token"]
	if v, ok := userToken.(map[string]string); ok {
		data["CREATOR_ID"] = v["ID"]
		data["CREATE_TIME"] = time.Now()
		data["UPDATER_ID"] = v["ID"]
		data["UPDATE_TIME"] = time.Now()
	}
	return true, nil
}
func (this *SecretDataInterceptor) AfterCreate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	text := fmt.Sprint(data["CLIENT"], "\t", data["SERVER"], "\t", data["SECRET"], "\t", data["IP_ADDRESSES"], "\n")
	_, err = f.WriteString(text)
	return err
}
func (this *SecretDataInterceptor) BeforeLoad(resourceId string, ds interface{}, context map[string]interface{}, id string) (bool, error) {
	userToken := context["user_token"]
	if v, ok := userToken.(map[string]string); ok {
		context["extra_filter"] = fmt.Sprint("AND CREATOR_ID='", v["ID"], "'")
	} else {
		return false, errors.New("Invalid user token.")
	}

	return true, nil
}
func (this *SecretDataInterceptor) AfterLoad(resourceId string, ds interface{}, context map[string]interface{}, data map[string]string) error {
	data["secret"] = ""
	return nil
}
func (this *SecretDataInterceptor) BeforeUpdate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	userToken := context["user_token"]
	if v, ok := userToken.(map[string]string); ok {
		data["UPDATER_ID"] = v["ID"]
		data["UPDATE_TIME"] = time.Now()
	}
	return true, nil
}
func (this *SecretDataInterceptor) AfterUpdate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	if db, ok := ds.(*sql.DB); ok {
		return updateFile(db)
	}
	return errors.New("Failed to access database.")
}
func (this *SecretDataInterceptor) AfterDelete(resourceId string, ds interface{}, context map[string]interface{}, id string) error {
	if db, ok := ds.(*sql.DB); ok {
		return updateFile(db)
	}
	return errors.New("Failed to access database.")
}
func (this *SecretDataInterceptor) BeforeListMap(resourceId string, ds interface{}, context map[string]interface{}, filter *string, sort *string, start int64, limit int64, includeTotal bool) (bool, error) {
	userToken := context["user_token"]
	if v, ok := userToken.(map[string]string); ok {
		userId := v["ID"]
		gorest.MysqlSafe(&userId)
		*filter += fmt.Sprint(" AND (CREATOR_ID='", userId, "')")
	} else {
		return false, errors.New("Invalid user.")
	}

	if db, ok := ds.(*sql.DB); ok {
		err := loadFromFile(db, context)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return true, nil
}
func (this *SecretDataInterceptor) AfterListMap(resourceId string, ds interface{}, context map[string]interface{}, data []map[string]string, total int64) error {
	for _, v := range data {
		v["secret"] = ""
	}
	return nil
}

func loadFromFile(db *sql.DB, context map[string]interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	userId := ""
	userToken := context["user_token"]
	if v, ok := userToken.(map[string]string); ok {
		userId = v["ID"]
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(text, "#") {
			continue
		}
		fields := strings.Fields(text)
		if len(fields) >= 4 {
			values := []interface{}{uuid.New(), fields[0], fields[1], fields[2], fields[3], "0"}
			values = append(values, userId, time.Now(), userId, time.Now())
			_, err := gosqljson.ExecDb(db, `INSERT OR IGNORE INTO secret(ID,CLIENT,SERVER,SECRET,IP_ADDRESSES,
			STATUS,CREATOR_ID,CREATE_TIME,UPDATER_ID,UPDATE_TIME) VALUES(?,?,?,?,?,?,?,?,?,?)`, values...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func updateFile(db *sql.DB) error {
	m, err := gosqljson.QueryDbToMap(db, false, "SELECT * FROM secret")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0600)
	defer f.Close()
	f.WriteString(header)
	for _, data := range m {
		text := fmt.Sprint(data["CLIENT"], "\t", data["SERVER"], "\t", data["SECRET"], "\t", data["IP_ADDRESSES"], "\n")
		_, err = f.WriteString(text)
	}
	return err
}
