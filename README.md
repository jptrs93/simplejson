# Simple JSON parser

`sjson` is a simple light weight package for parsing and manipulating JSON data of unknown structure. Once parsed, the structure of the JSON is captured in a general purpose struct allowing it to be repeatedly queried and manipulated efficiently. This contrasts with other packages aimed at working with unstructured JSON such as [jsonparser](https://github.com/buger/jsonparser) where the bytes must be parsed upon every query. The parsing is performed in a single pass of the bytes and may be done re-using a byte slice without copying.

## Example usage

```go
package main

import (
	"fmt"
	"github.com/jptrs93/sjson"
)

var data = []byte(`
{
	"person": {
		"name": {
		"first": "Joe",
		"last": "Smith",
	},
		"age": 29,
		"emails" : ["doesnotexist@nowhere.com"]
	}
}
`)

// load the json directly from some bytes
// alternatively use sjson.Parse(scanner) to read directly from any io.RuneScanner
json, err := sjson.ParseUTF8(data)

// get the json at path "person.emails"
subjson, err := json.Get("person", "emails")

subjson.IsArray() // test if the subjson is of type ARRAY

// iterate for the items in the subjson array
for _, item := range subjson.ArrayItems() {
  strVal, err := item.AsString() // get the item as a string
}

// get the value at path "person.name.first" directly as a string
firstName, err := json.GetAsString("person", "name", "first")

// get the value at path "person.age" directly as a integer
age, err := json.GetAsInt("person", "age")

```
