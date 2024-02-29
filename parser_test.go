package sjson

import (
	"log"
	"sort"
	"testing"
)

func TestParseUTF8_ValidJsonWithString(t *testing.T) {
	json_str := `{"key1":"val1","key2":"val2"}`
	if json, err := ParseUTF8([]byte(json_str)); err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	} else {
		containsKeys(json, []string{"key1", "key2"}, t)
		strVal, _ := json.ObjectItems()["key1"].AsString()
		if strVal != "val1" {
			t.Errorf("Expected 'val1' but was %v", strVal)
		}
	}
}

func TestParseUTF8_ValidJsonArray(t *testing.T) {
	json_str := `[1,2,3,4]`
	if json, err := ParseUTF8([]byte(json_str)); err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	} else {
		if len(json.ArrayItems()) != 4 {
			t.Errorf("Json array length %v does not match expected length 4", len(json.ArrayItems()))
		}
	}
}

func TestParseUTF8_ValidJsonWithNumber(t *testing.T) {
	json_str := `{"key1":123,"key2":1.223}`
	if json, err := ParseUTF8([]byte(json_str)); err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	} else {
		containsKeys(json, []string{"key1", "key2"}, t)
		intVal, _ := json.ObjectItems()["key1"].AsInt()
		if intVal != 123 {
			t.Errorf("Expected '123' but was %v", intVal)
		}
		floatVal, _ := json.ObjectItems()["key2"].AsFloat64()
		if floatVal != 1.223 {
			t.Errorf("Expected float value 1.223 but was %v", floatVal)
		}
	}
}

func TestParseUTF8_ValidJsonRootFloat(t *testing.T) {
	json_str := `1.223`
	if json, err := ParseUTF8([]byte(json_str)); err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	} else {
		if v, err := json.AsFloat64(); err != nil {
			t.Errorf("Error converting json to float %v", err.Error())
		} else if v != 1.223 {
			t.Errorf("Expected float value 1.223 but was %v", v)
		}
	}
}

func TestParseUTF8_ValidMixedNested(t *testing.T) {
	json_str := `{"key1":"stringVal","key2":1000,"key3":20.11,"key4":[{"key4.1":true,"key4.2":false,"key4.3":null}],"key5":{"key5.1":1}}`
	if json, err := ParseUTF8([]byte(json_str)); err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	} else {
		containsKeys(json, []string{"key1", "key2", "key3", "key4", "key5"}, t)
	}
}

func TestParse_ValidMixedNested(t *testing.T) {
	json_str := `{"key1":"stringVal","key2":1000,"key3":20.11,"key4":[{"key4.1":true,"key4.2":false,"key4.3":null}],"key5":{"key5.1":1}}`
	if json, err := Parse(&UTF8RuneScanner{bytes: []byte(json_str)}); err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	} else {
		containsKeys(json, []string{"key1", "key2", "key3", "key4", "key5"}, t)
	}
}

func TestString(t *testing.T) {
	json_str := `{"key0":{"a":1},"key1":"stringVal","key2":1000,"key3":20.11,"key4":[{"key4.1":true,"key4.2":false,"key4.3":null}],"key5":{"key5.1":1}}`
	json, err := Parse(&UTF8RuneScanner{bytes: []byte(json_str)})
	if err != nil {
		t.Errorf("Error (%v) parsing '%v' when none expected", err.Error(), json_str)
	}
	for _, k := range json.Keys() {
		v, _ := json.Get(k)
		log.Printf("%v", v)
	}
}

func isEqual[T comparable](s1, s2 []T) bool {
	if len(s1) == len(s2) {
		for i, item := range s1 {
			if item != s2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func containsKeys(j *Json, expectedKeys []string, t *testing.T) {
	keys := j.Keys()
	sort.Strings(keys)
	if !isEqual(keys, expectedKeys) {
		t.Errorf("Json keys %v do not match expected key %v", keys, expectedKeys)
	}
}
