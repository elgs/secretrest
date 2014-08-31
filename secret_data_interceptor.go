package main

import (
	"fmt"
	"github.com/elgs/gorest"
)

func init() {
	tableId := "secret"
	gorest.RegisterDataInterceptor(tableId, &MyDataInterceptor{Id: "secret"})
}

type SecretDataInterceptor struct {
	*gorest.EchoDataInterceptor
	Id string
}

func (this *SecretDataInterceptor) BeforeCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	fmt.Println("Here I'm in BeforeCreate")
	//if db, ok := ds.(*sql.DB); ok {
	//	_ = db
	//}
	return true, nil
}
func (this *SecretDataInterceptor) AfterCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	fmt.Println("Here I'm in AfterCreate")
	//if db, ok := ds.(*sql.DB); ok {
	//	_ = db
	//}
	return nil
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
