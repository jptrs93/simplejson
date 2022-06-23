package sjson

import (
	"fmt"
	"strconv"
	"strings"
)

type Json struct {
	bytes  *[]byte
	start  int
	end    int
	parent *Json

	objectItems map[string]*Json // applicable when jtype is OBJECT
	arrayItems  []*Json          // applicable when jtype is ARRAY

	jtype int // one of (OBJECT, ARRAY, NULL, BOOL, NUMBER, STRING)
}

// json element types
const (
	OBJECT = iota
	ARRAY  = iota
	NULL   = iota
	BOOL   = iota
	NUMBER = iota
	STRING = iota
)

func (j *Json) IsArray() bool {
	return j.jtype == ARRAY
}

func (j *Json) IsObject() bool {
	return j.jtype == OBJECT
}

func (j *Json) IsPrimitive() bool {
	return j.jtype != ARRAY && j.jtype != OBJECT
}

func (j *Json) IsString() bool {
	return j.jtype == STRING
}

func (j *Json) IsBool() bool {
	return j.jtype == BOOL
}

func (j *Json) IsNumber() bool {
	return j.jtype == NUMBER
}

func (j *Json) IsFloat() bool {
	return j.IsNumber() && strings.ContainsRune(string((*j.bytes)[j.start:j.end]), '.')
}

func (j *Json) IsInt() bool {
	return j.IsNumber() && !j.IsFloat()
}

func (j *Json) IsNull() bool {
	return j.jtype == NULL
}

func (j *Json) Get(keyPath ...string) (*Json, error) {
	for i, key := range keyPath {
		if val, ok := j.objectItems[key]; ok {
			j = val
		} else {
			path := strings.Join(keyPath[:i+1], ".")
			return nil, fmt.Errorf("%w: %v", KeyPathError, path)
		}
	}
	return j, nil
}

// todo support mutation of json
func (j *Json) Put(name string, value *Json) (*Json, error) {
	return j, nil
}

func (j *Json) AsString() (string, error) {
	if !j.IsString() {
		return "", StringValueError
	}
	if j.parent != nil {
		return string((*j.bytes)[j.start+1 : j.end-1]), nil
	}
	return string((*j.bytes)[j.start:j.end]), nil
}

func (j *Json) AsBool() (bool, error) {
	if !j.IsBool() {
		return false, BoolValueError
	}
	return (*j.bytes)[j.start] == 't', nil
}

func (j *Json) AsInt() (int, error) {
	if !j.IsNumber() {
		return 0, NumberValueError
	}
	r, err := strconv.Atoi(string((*j.bytes)[j.start:j.end]))
	return r, err
}

func (j *Json) AsFloat64() (float64, error) {
	if !j.IsNumber() {
		return 0, NumberValueError
	}
	r, err := strconv.ParseFloat(string((*j.bytes)[j.start:j.end]), 64)
	return r, err
}

func (j *Json) GetAsInt(keyPath ...string) (int, error) {
	if val, err := j.Get(keyPath...); err != nil {
		return 0, err
	} else {
		return val.AsInt()
	}
}

func (j *Json) GetAsFloat64(keyPath ...string) (float64, error) {
	if val, err := j.Get(keyPath...); err != nil {
		return 0, err
	} else {
		return val.AsFloat64()
	}
}

func (j Json) GetAsBool(keyPath ...string) (bool, error) {
	if val, err := j.Get(keyPath...); err != nil {
		return false, err
	} else {
		return val.AsBool()
	}
}

func (j Json) GetAsString(keyPath ...string) (string, error) {
	if val, err := j.Get(keyPath...); err != nil {
		return "", err
	} else {
		return val.AsString()
	}
}

func (j *Json) ArrayItems() []*Json {
	return j.arrayItems
}

func (j *Json) ObjectItems() map[string]*Json {
	return j.objectItems
}

func (j *Json) Keys() []string {
	keys := make([]string, len(j.objectItems))
	i := 0
	for k := range j.objectItems {
		keys[i] = k
		i++
	}
	return keys
}

func (j Json) Bytes() []byte {
	return (*j.bytes)[j.start : j.end+1]
}
