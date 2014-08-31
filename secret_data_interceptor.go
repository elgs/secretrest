package main

import (
	"fmt"
	"github.com/elgs/gorest"
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

func (this *SecretDataInterceptor) BeforeCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	userId := context["user_id"]
	data["CREATOR_ID"] = userId
	data["CREATE_TIME"] = time.Now()
	data["UPDATER_ID"] = userId
	data["UPDATE_TIME"] = time.Now()
	return true, nil
}
func (this *SecretDataInterceptor) AfterCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	//if db, ok := ds.(*sql.DB); ok {
	//	_ = db
	//}
	f, err := os.OpenFile("/Users/elgs/Desktop/chap-secrets", os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	text := fmt.Sprint(data["CLIENT"], "\t", data["SERVER"], "\t", data["SECRET"], "\t", data["IP_ADDRESSES"], "\n")
	_, err = f.WriteString(text)
	return err
}
func (this *SecretDataInterceptor) BeforeLoad(ds interface{}, context map[string]interface{}, id string) (bool, error) {
	fmt.Println("Here I'm in BeforeLoad")
	return true, nil
}
func (this *SecretDataInterceptor) AfterLoad(ds interface{}, context map[string]interface{}, data map[string]string) error {
	fmt.Println("Here I'm in AfterLoad")
	return nil
}
func (this *SecretDataInterceptor) BeforeUpdate(ds interface{}, context map[string]interface{}, oldData map[string]interface{}, data map[string]interface{}) (bool, error) {
	fmt.Println("Here I'm in BeforeUpdate")
	return true, nil
}
func (this *SecretDataInterceptor) AfterUpdate(ds interface{}, context map[string]interface{}, oldData map[string]interface{}, data map[string]interface{}) error {
	fmt.Println("Here I'm in AfterUpdate")
	return nil
}
func (this *SecretDataInterceptor) AfterDelete(ds interface{}, context map[string]interface{}, id string) error {
	fmt.Println("Here I'm in AfterDelete")
	return nil
}
func (this *SecretDataInterceptor) BeforeListMap(ds interface{}, context map[string]interface{}, filter *string, sort *string, start int64, limit int64, includeTotal bool) (bool, error) {

	return true, nil
}
func (this *SecretDataInterceptor) AfterListMap(ds interface{}, context map[string]interface{}, data []map[string]string, total int64) error {
	fmt.Println("Here I'm in AfterListMap")
	return nil
}
