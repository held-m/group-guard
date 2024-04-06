// Package groupguard checks groups of a struct and create new object
package groupguard

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Filter is a function that filters the group members based on the group's rules
func Filter[T interface{}](groups []string, obj T) (T, error) {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	newObj := reflect.New(objType).Elem()

	if objValue.Kind() != reflect.Struct {
		return newObj.Interface().(T), errors.New("obj must be a struct")
	}

	for i := 0; i < objType.NumField(); i++ {
		if objType.Field(i).Type.Kind() == reflect.Struct {
			newSubObj, err := Filter(groups, objValue.Field(i).Interface())
			if err != nil {
				return newObj.Interface().(T), err
			}
			newObj.Field(i).Set(reflect.ValueOf(newSubObj))
			continue
		}

		addField(groups, objType, objValue, newObj, i)
	}

	fmt.Println(newObj.Interface().(T))
	return newObj.Interface().(T), nil
}

func addField(
	groups []string,
	objType reflect.Type,
	objValue reflect.Value,
	newObj reflect.Value,
	fieldIndex int,
) {
	fieldType := objType.Field(fieldIndex)
	fieldValue := objValue.Field(fieldIndex)

	groupsStruct := parseGroups(&fieldType)
	if groupsStruct == nil {
		return
	}

	if compareGroups(groups, groupsStruct) {
		newObj.Field(fieldIndex).Set(fieldValue)
	}
}

func parseGroups(field *reflect.StructField) []string {

	groupStr := strings.ReplaceAll(field.Tag.Get("group"), " ", "")
	if groupStr == "" {
		return nil
	}

	groups := strings.Split(groupStr, ",")

	return groups
}

func compareGroups(groups []string, groupsStruct []string) bool {
	for _, groupObj := range groupsStruct {
		for _, group := range groups {
			if group == groupObj {
				return true
			}
		}
	}
	return false
}
