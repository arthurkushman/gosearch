package gosearch

import (
	"github.com/buger/jsonparser"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

const (
	QUERY = "query"
	TERM = "term"
)

type Core interface {
	PerformSearch()
	ParseInput(in map[string]interface{})
	SearchPhrase()
	BuildIndex()
	SetCanonicalIndex()
}

func PerformSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index := vars["index"]
	indexType := vars["indextype"]
	indexId := vars["id"]
	if indexId != "" { // search by id

	} else { // search by doc

	}
}

func ParseInput(in map[string]interface{}) map[string]string {
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		fmt.Println(jsonparser.Get(value, "url"))
	}, "query", "term")
	return fieldValueMap
}

func BuildIndex() {

}