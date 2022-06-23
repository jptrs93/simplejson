# Simple JSON parser

`sjson` is a simple light weight package for parsing and manipulating JSON data of unknown structure. Once parsed, the structure of the JSON is captured in a general purpose struct allowing it to be repeatedly queried and manipulated efficiently. This contrasts with other packages aimed at working with unstructured JSON such as [jsonparser](https://github.com/buger/jsonparser) where the bytes must be parsed upon every query. The parsing is performed in a single pass of the bytes and may be done re-using a byte slice without copying.

## Example usage

```go
import "github.com/jptrs93/sjson"

data := []byte(`
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

func main() {
  # load the json directly from some bytes
  # alternatively use Parse(scanner) to read directly from any io.RuneScanner
  json = ParseUTF8(data)
  
  
  # get the item at path "person.name.first" as a string
  firstName, err := json.GetAsString("person", "name", "first")
  
  
}
