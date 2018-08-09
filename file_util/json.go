/* json.go - load json from file, check nil  */
/*
modification history
--------------------
2014/8/1, by Weiwei, create
2015/4/9, by zhangjiyang01, mv form bfe_util to golang-lib
*/
/*
DESCRIPTION
*/
package file_util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"time"
)

// load json content from file, unmarshal to jsonObject
// check if all field is set if checkNilPointer is true
// check if any field is not pointer type if allowNoPointerField is false
func LoadJsonFile(path string, jsonObject interface{}) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, jsonObject); err != nil {
		return err
	}

	return nil
}

// dump json file
func DumpJson(jsonObject interface{}, filePath string, perm os.FileMode) error {
	buf, err := json.MarshalIndent(jsonObject, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal err %s", err)
	}

	// mkdirall dir
	dirPath := path.Dir(filePath)
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("MkdirALl err %s", err.Error())
	}

	return ioutil.WriteFile(filePath, buf, perm)
}

// dump process:
// 1, dump json object to filename.{currenttime}
// 2, bak filename to filename.bak, if filename exist
// 3, copy filename.{currenttime} to filename
func AtomicDumpJson(v interface{}, filename string) error {
	bakFile := fmt.Sprintf("%s.%s", filename, time.Now().Format("20060102150405"))

	// write to bakFile
	if err := DumpJson(v, bakFile, 0744); err != nil {
		return fmt.Errorf("write file err %s", err)
	}

	// copy to filename
	if err := AtomicCopy(bakFile, filename); err != nil {
		return err
	}

	return nil
}

// check if a struct has a nil field
// if allowNoPointerField is false, it also check if fields are all pointers
// if param object is not a struct , return nil
func CheckNilField(object interface{}, allowNoPointerField bool) error {
	v := reflect.ValueOf(object)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input is not struct")
	}

	typeOfV := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr {
			if !allowNoPointerField {
				return fmt.Errorf("%s field %s is not a pointer", typeOfV, typeOfV.Field(i).Name)
			}
			continue
		}

		if f.IsNil() {
			return fmt.Errorf("%s field %s is not set", typeOfV, typeOfV.Field(i).Name)
		}
	}
	return nil
}
