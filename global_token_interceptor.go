package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elgs/gorest"
	"github.com/elgs/gosqljson"
	"io/ioutil"
	"time"
)

func init() {
	loadACL()
	gorest.RegisterGlobalDataInterceptor(&GlobalTokenInterceptor{Id: "GlobalTokenInterceptor"})
}

var acl = make(map[string]map[string]bool)
var tokenRegistry = make(map[string]map[string]string)

func checkToken(db *sql.DB, id string, key string, context map[string]interface{}, tableId string) (bool, error) {
	if id != "" && key != "" && len(tokenRegistry[id]) > 0 && tokenRegistry[id]["TOKEN_KEY"] == key {
		context["user_token"] = tokenRegistry[id]
		return true, nil
	}
	tokenTable := context["token_table"]
	if tokenTable != nil {
		tokenTableString := tokenTable.(string)
		if len(tokenTableString) > 0 && tableId == tokenTable {
			return false, errors.New("We think you are invading the system.")
		}

		gorest.MysqlSafe(&tokenTableString)
		data, err := gosqljson.QueryDbToMap(db, false, fmt.Sprint("SELECT * FROM ", tokenTableString, " WHERE ID=? AND TOKEN_KEY=? AND STATUS=?"), id, key, "0")
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		if data != nil && len(data) == 1 {
			record := data[0]
			tokenRegistry[record["ID"]] = record
			context["user_token"] = record
			return true, nil
		}
	}
	return false, errors.New("Authentication failed.")
}

func loadACL() {
	// load acl from configuration files.
	configFile := "gorest_acl.json"
	aclConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(configFile, " not found, default policies are used.")
	}
	err = json.Unmarshal(aclConfig, &acl)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(len(acl), acl)
}

func checkACL(tableId string, op string) (bool, error) {
	if acl[tableId] != nil {
		if acl[tableId][op] {
			return true, nil
		}
	}
	return false, errors.New("Access denied.")
}

type GlobalTokenInterceptor struct {
	*gorest.DefaultDataInterceptor
	Id string
}

func (this *GlobalTokenInterceptor) BeforeCreate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	if ok, err := checkACL(resourceId, "create"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	ctn, err := checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
	if ctn && err == nil {
		if context["meta"] != nil && context["meta"].(bool) {
			userToken := context["user_token"]
			if v, ok := userToken.(map[string]string); ok {
				data["CREATOR_ID"] = v["USER_ID"]
				data["CREATE_TIME"] = time.Now()
				data["UPDATER_ID"] = v["USER_ID"]
				data["UPDATE_TIME"] = time.Now()
			}
		}
	}
	return ctn, err
}
func (this *GlobalTokenInterceptor) AfterCreate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeLoad(resourceId string, ds interface{}, context map[string]interface{}, id string) (bool, error) {
	if ok, err := checkACL(resourceId, "load"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterLoad(resourceId string, ds interface{}, context map[string]interface{}, data map[string]string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeUpdate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	if ok, err := checkACL(resourceId, "update"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	ctn, err := checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
	if ctn && err == nil {
		if context["meta"] != nil && context["meta"].(bool) {
			userToken := context["user_token"]
			if v, ok := userToken.(map[string]string); ok {
				data["UPDATER_ID"] = v["USER_ID"]
				data["UPDATE_TIME"] = time.Now()
			}
		}
	}
	return ctn, err
}
func (this *GlobalTokenInterceptor) AfterUpdate(resourceId string, ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeDuplicate(resourceId string, ds interface{}, context map[string]interface{}, id string) (bool, error) {
	if ok, err := checkACL(resourceId, "duplicate"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterDuplicate(resourceId string, ds interface{}, context map[string]interface{}, id string, newId string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeDelete(resourceId string, ds interface{}, context map[string]interface{}, id string) (bool, error) {
	if ok, err := checkACL(resourceId, "delete"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterDelete(resourceId string, ds interface{}, context map[string]interface{}, id string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeListMap(resourceId string, ds interface{}, context map[string]interface{}, filter *string, sort *string, start int64, limit int64, includeTotal bool) (bool, error) {
	if ok, err := checkACL(resourceId, "list"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterListMap(resourceId string, ds interface{}, context map[string]interface{}, data []map[string]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeListArray(resourceId string, ds interface{}, context map[string]interface{}, filter *string, sort *string, start int64, limit int64, includeTotal bool) (bool, error) {
	if ok, err := checkACL(resourceId, "list"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterListArray(resourceId string, ds interface{}, context map[string]interface{}, data [][]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeQueryMap(resourceId string, ds interface{}, context map[string]interface{}, sqlSelect *string, sqlSelectCount *string, start int64, limit int64, includeTotal bool) (bool, error) {
	if ok, err := checkACL(resourceId, "query"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterQueryMap(resourceId string, ds interface{}, context map[string]interface{}, data []map[string]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeQueryArray(resourceId string, ds interface{}, context map[string]interface{}, sqlSelect *string, sqlSelectCount *string, start int64, limit int64, includeTotal bool) (bool, error) {
	if ok, err := checkACL(resourceId, "query"); !ok {
		return false, err
	}
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterQueryArray(resourceId string, ds interface{}, context map[string]interface{}, data [][]string, total int64) error {
	return nil
}
