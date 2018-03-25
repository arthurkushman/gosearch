package core

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

type Get interface {
	SetOpType()
	SetOpStatus()
	SetSource()
	SetId()
	SetVersion()
	SetTimestamp()
	PerformSearch()
}

type Fields struct {
	OpType    string
	OpStatus  bool
	Source    string // json source
	Id        uint64 // doc id
	Index     string
	IndexType string
	incrKey   string
	Version   int
	Timestamp int
}

type Storage struct {
	IncrKey string
}

type StoreFields struct {
	Fld Fields
	Stg Storage
}

func (sf *StoreFields) PerformSearch(w http.ResponseWriter, r *http.Request) {
	route := mux.Vars(r)
	sf.Fld.Id, _ = strconv.ParseUint(route["id"], 10, 64)
	if sf.Fld.Id > 0 { // search by id, if user set route to id
		sf.SearchById()
	} else { // search by doc
		buf := make([]byte, ReadBufferSize)
		sf.Fld.Index = route["Index"]
		sf.Fld.IndexType = route["indextype"]
		n, err := r.Body.Read(buf)
		if err != nil || n > ReadBufferSize {
			panic("Error reading from input stream")
		}
		input := ParseInput(buf)
		sf.SearchPhrase(input)
	}
}
