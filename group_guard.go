// Package groupguard checks groups of a struct and create new object
package groupguard

import (
	"errors"
	"reflect"
	"strings"
)

// Filter is a function that filters the group members based on the group's rules
func Filter(groups []string, obj interface{}) (reflect.Value, error) {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	newObj := reflect.New(objType).Elem()

	if objValue.Kind() != reflect.Struct {
		return newObj, errors.New("obj must be a struct")
	}

	for i := 0; i < objType.NumField(); i++ {
		if objType.Field(i).Type.Kind() == reflect.Struct {
			newSubObj, err := Filter(groups, objValue.Field(i).Interface())
			if err != nil {
				return newObj, err
			}
			newObj.Field(i).Set(newSubObj)
			continue
		}

		addField(groups, objType, objValue, newObj, i)
	}

	return newObj, nil
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
