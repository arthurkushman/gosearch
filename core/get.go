package core

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
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

func (sf *StoreFields) PerformSearch(w http.ResponseWriter, r *http.Request) {
	route := mux.Vars(r)
	sf.Fld.Id, _ = strconv.ParseUint(route["id"], 10, 64)
	sf.Fld.IsSearch = true

	if sf.Fld.Id > 0 { // search by id, if user set route to id
		sf.SearchById(w)
	} else { // search by doc
		sf.Fld.Index = route["Index"]
		sf.Fld.IndexType = route["indextype"]
		sf.SetSourceDocument(r)

		fmt.Println(sf.Fld.RequestSource)
		input := ParseInput(sf.Fld.RequestSource)
		fmt.Println(input)

		sf.SearchPhrase(input)
	}

	EchoResult(w, sf.GetJsonOutput(), http.StatusOK)
}
