package groupguard

import (
	"reflect"
	"testing"
)

type SubTestStruct struct {
	FieldWithGroup    string `group:"group1"`
	FieldWithTwoGroup string `group:"group1, group2"`
	FieldWithoutGroup string
}

type TestStruct struct {
	SubTestStruct
	FieldWithGroup    string `group:"    group1"`
	FieldWithTwoGroup string `group:"group1, group2   "`
	FieldWithoutGroup string
}

func TestParceGroups(t *testing.T) {

	testObj := TestStruct{
		FieldWithTwoGroup: "FieldWithTwoGroup",
		FieldWithoutGroup: "FieldWithoutGroup",
	}

	objType := reflect.TypeOf(testObj)

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		groups := parseGroups(&field)

		// Test TestStruct
		switch field.Name {
		case "FieldWithGroup":
			if len(groups) != 1 || groups[0] != "group1" {
				t.Errorf("FieldWithGroup: %v", groups)
			}
		case "FieldWithTwoGroup":
			if len(groups) != 2 || groups[0] != "group1" || groups[1] != "group2" {
				t.Errorf("FieldWithTwoGroup: %v", groups)
			}
		case "FieldWithoutGroup":
			if len(groups) != 0 {
				t.Errorf("FieldWithoutGroup: %v", groups)
			}
		}

		// Test SubTestStruct
		if field.Name == "SubTestStruct" {
			for j := 0; j < field.Type.NumField(); j++ {
				subField := field.Type.Field(j)
				subGroups := parseGroups(&subField)

				switch subField.Name {
				case "FieldWithGroup":
					if len(subGroups) != 1 || subGroups[0] != "group1" {
						t.Errorf("FieldWithGroup: %v", subGroups)
					}
				case "FieldWithTwoGroup":
					if len(subGroups) != 2 || subGroups[0] != "group1" || subGroups[1] != "group2" {
						t.Errorf("FieldWithTwoGroup: %v", subGroups)
					}
				case "FieldWithoutGroup":
					if len(subGroups) != 0 {
						t.Errorf("FieldWithoutGroup: %v", subGroups)
					}
				}
			}
		}
	}
}

func TestFilter(t *testing.T) {

	testObj := TestStruct{
		SubTestStruct: SubTestStruct{
			FieldWithGroup:    "SubFieldWithGroup",
			FieldWithTwoGroup: "SubFieldWithTwoGroup",
			FieldWithoutGroup: "SubFieldWithoutGroup",
		},
		FieldWithGroup:    "FieldWithGroup",
		FieldWithTwoGroup: "FieldWithTwoGroup",
		FieldWithoutGroup: "FieldWithoutGroup",
	}

	groups := []string{"group1"}

	newObj, err := Filter(groups, testObj)

	if err != nil {
		t.Errorf("Filter: %v", err)
	}

	newObjType := reflect.TypeOf(newObj)

	for i := 0; i < newObjType.NumField(); i++ {
		field := newObjType.Field(i)

		switch field.Name {
		case "FieldWithGroup":
			if field.Type.Kind() != reflect.String {
				t.Errorf("FieldWithGroup: %v", field.Type.Kind())
			}
		case "FieldWithTwoGroup":
			if field.Type.Kind() != reflect.String {
				t.Errorf("FieldWithTwoGroup: %v", field.Type.Kind())
			}
		case "FieldWithoutGroup":
			if field.Type.Kind() != reflect.String {
				t.Errorf("FieldWithoutGroup: %v", field.Type.Kind())
			}
		case "SubTestStruct":
			for j := 0; j < field.Type.NumField(); j++ {
				subField := field.Type.Field(j)

				switch subField.Name {
				case "FieldWithGroup":
					if subField.Type.Kind() != reflect.String {
						t.Errorf("FieldWithGroup: %v", subField.Type.Kind())
					}
				case "FieldWithTwoGroup":
					if subField.Type.Kind() != reflect.String {
						t.Errorf("FieldWithTwoGroup: %v", subField.Type.Kind())
					}
				case "FieldWithoutGroup":
					if subField.Type.Kind() != reflect.String {
						t.Errorf("FieldWithoutGroup: %v", subField.Type.Kind())
					}
				}
			}
		}
	}
}

func TestCompareGroups(t *testing.T) {

	groups := []string{"group1", "group2"}

	groupsStruct := []string{"group1"}

	if !compareGroups(groups, groupsStruct) {
		t.Errorf("compareGroups: %v", groupsStruct)
	}

	groupsStruct = []string{"group1", "group2"}

	if !compareGroups(groups, groupsStruct) {
		t.Errorf("compareGroups: %v", groupsStruct)
	}

	groupsStruct = []string{"group1", "group2", "group3"}

	if !compareGroups(groups, groupsStruct) {
		t.Errorf("compareGroups: %v", groupsStruct)
	}

	groupsStruct = []string{"group3"}

	if compareGroups(groups, groupsStruct) {
		t.Errorf("compareGroups: %v", groupsStruct)
	}
}
