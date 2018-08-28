package core

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

func (sf *StoreFields) BuildIndex(w http.ResponseWriter, r *http.Request) {
	route := mux.Vars(r)
	sf.Fld.Index = route["index"]
	sf.Fld.IndexType = route["indextype"]
	fmt.Println(sf.Fld)

	// perform redis connection
	sf.redisConn()

	sf.SetIncrKey()

	// start indexing
	tStart := GetMillis()

	var created = false
	docInfo, _ := sf.GetDocInfo()
	fmt.Println(docInfo)

	if docInfo == nil { // insert
		sf.SetCanonicalIndex()
		//$this- > insert()
		sf.Fld.Result = ResultCreated
		created = true
	} else { // update
		//$this- > updateDocInfo($docInfo)
		sf.Fld.Result = ResultUpdated
	}

	// store specific info
	sf.Fld.Took = GetMillis() - tStart
	sf.Fld.OpType = ResultCreated
	sf.Fld.OpStatus = created

	// output json result of created index document
	EchoResult(w, sf.GetJsonOutput(), http.StatusCreated)
}
