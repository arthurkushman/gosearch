package gosearch

import (
	"testing"
	"reflect"
	"fmt"
)

func TestSer(t *testing.T) {
	mapObj := map[string]interface{}{
		"title":"val1",
		"desc":"val2",
		"foo": []byte(``),
	}
	dataBytes := Ser(mapObj)
	fmt.Println(string(dataBytes))
}

func TestUnser(t *testing.T) {
	mapData := Unser([]byte(`{
    "sendMsg":{"user":"UserName","msg":"Trying to send a message"},
    "say":"Hello"
	}`))
	if reflect.ValueOf(mapData).Kind() != reflect.Map {
		t.Error("Error: Unserialized value is not a map")
	}
}
