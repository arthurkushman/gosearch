package core

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
)

func (sf *StoreFields) BuildIndex(w http.ResponseWriter, r *http.Request) {
	route := mux.Vars(r)
	sf.Fld.Index = route["Index"]
	sf.Fld.IndexType = route["indextype"]
	fmt.Println(sf.Fld)
	os.Exit(1)

	sf.SetIncrKey()
	// start indexing
	tStart := GetMillis()
	var created = false
	docInfo, _ := sf.GetDocInfo()
	if docInfo == nil { // insert
		sf.SetCanonicalIndex()
		//$this- > insert()
		sf.Fld.Result = RESULT_CREATED
		created = true
	} else { // update
		//$this- > updateDocInfo($docInfo)
		sf.Fld.Result = RESULT_UPDATED
	}
	sf.Fld.Took = GetMillis() - tStart
	sf.Fld.OpType = RESULT_CREATED
	sf.Fld.OpStatus = created
	sf.JsonOutput()
}
