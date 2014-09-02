package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgs/gorest"
	"github.com/elgs/gosqljson"
	"os"
	"time"
)

func init() {
	tableId := "secret"
	gorest.RegisterDataInterceptor(tableId, &SecretDataInterceptor{Id: "secret"})
}

type SecretDataInterceptor struct {
	*gorest.EchoDataInterceptor
	Id string
}

var header string = `# Secrets for authentication using CHAP
# client	server	secret			IP addresses
`

func (this *SecretDataInterceptor) BeforeCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	userId := context["user_id"]
	data["CREATOR_ID"] = userId
	data["CREATE_TIME"] = time.Now()
	data["UPDATER_ID"] = userId
	data["UPDATE_TIME"] = time.Now()
	return true, nil
}
func (this *SecretDataInterceptor) AfterCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	f, err := os.OpenFile("/Users/elgs/Desktop/chap-secrets", os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	text := fmt.Sprint(data["CLIENT"], "\t", data["SERVER"], "\t", data["SECRET"], "\t", data["IP_ADDRESSES"], "\n")
	_, err = f.WriteString(text)
	return err
}
func (this *SecretDataInterceptor) BeforeLoad(ds interface{}, context map[string]interface{}, id string) (bool, error) {
	context["extra_filter"] = context["user_id"]
	return true, nil
}
func (this *SecretDataInterceptor) BeforeUpdate(ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	userId := context["user_id"]
	data["UPDATER_ID"] = userId
	data["UPDATE_TIME"] = time.Now()
	return true, nil
}
func (this *SecretDataInterceptor) AfterUpdate(ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	if db, ok := ds.(*sql.DB); ok {
		return updateFile(db)
	}
	return errors.New("Failed to access database.")
}
func (this *SecretDataInterceptor) AfterDelete(ds interface{}, context map[string]interface{}, id string) error {
	if db, ok := ds.(*sql.DB); ok {
		return updateFile(db)
	}
	return errors.New("Failed to access database.")
}
func (this *SecretDataInterceptor) BeforeListMap(ds interface{}, context map[string]interface{}, filter *string, sort *string, start int64, limit int64, includeTotal bool) (bool, error) {
	*filter += fmt.Sprint(" AND (CREATOR_ID='", context["user_id"], "')")
	if db, ok := ds.(*sql.DB); ok {
		err := loadFromFile(db)
		if err != nil {
			return true, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func loadFromFile(db *sql.DB) error {
	file, err := os.Open("/Users/elgs/Desktop/chap-secrets")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func updateFile(db *sql.DB) error {
	m, err := gosqljson.QueryDbToMap(db, false, "SELECT * FROM secret")
	f, err := os.OpenFile("/Users/elgs/Desktop/chap-secrets", os.O_WRONLY, 0600)
	defer f.Close()
	f.WriteString(header)
	for _, data := range m {
		text := fmt.Sprint(data["CLIENT"], "\t", data["SERVER"], "\t", data["SECRET"], "\t", data["IP_ADDRESSES"], "\n")
		_, err = f.WriteString(text)
	}
	return err
}
