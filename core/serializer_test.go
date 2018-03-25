package core

import (
	"testing"
	"reflect"
)

func TestSer(t *testing.T) {
	mapObj := map[string]string{
		"title":       "foo",
		"description": "foo bar baz",
	}
	dataBytes := Ser(mapObj)
	if len(dataBytes) != len([]byte(`{"title":"foo","description":"foo bar baz"}`)) {
		t.Error("Error: data bytes are not equal")
	}
}

func TestUnser(t *testing.T) {
	mapData := Unser([]byte(`{"title": "foo","description": "foo bar baz"}`))
	if reflect.ValueOf(mapData).Kind() != reflect.Map {
		t.Error("Error: Unserialized value is not a map")
	}
	if mapData["title"] != "foo" {
		t.Error("Error: data is not valid")
	}
}
