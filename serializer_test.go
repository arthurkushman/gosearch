package gosearch

import (
	"testing"
	"reflect"
)

func TestUnser(t *testing.T) {
	mapData := Unser([]byte(`{
    "sendMsg":{"user":"UserName","msg":"Trying to send a message"},
    "say":"Hello"
	}`))
	if reflect.ValueOf(mapData).Kind() != reflect.Map {
		t.Error("Error: Unserialized value is not a map")
	}
}
