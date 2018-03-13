package gosearch

import (
	"testing"
	"reflect"
	"fmt"
)

func TestSer(t *testing.T) {
	mapObj := map[string]string{
		"title":       "foo",
		"description": "foo bar baz",
	}
	dataBytes := Ser(mapObj)
	fmt.Println(string(dataBytes))
}

func TestUnser(t *testing.T) {
	mapData := Unser([]byte(`{"title": "foo","description": "foo bar baz"}`))
	if reflect.ValueOf(mapData).Kind() != reflect.Map {
		t.Error("Error: Unserialized value is not a map")
	}
	fmt.Println(mapData)
}
