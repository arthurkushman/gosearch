package core

import (
	"time"
	"encoding/json"
)

/** PUT
{
    "created": true,
    "took": 40,
    "_index": "myindex",
    "_type": "myindextype",
    "_id": 1,
    "result": "created",
    "_version": 1
}
 */

func (sf *StoreFields) GetJsonOutput() []byte {
	if sf.Fld.OpType == ResultCreated {
		return getJsonCreated(sf.Fld)
	}

	if sf.Fld.OpType == ResultFound {
		return getJsonFound(sf.Fld)
	}

	return []byte{}
}

func getJsonCreated(fields Fields) []byte {
	data := make(map[string]interface{})

	data[ResultCreated] = true
	data[Took] = fields.Took
	data[Index] = fields.Index
	data[Type] = fields.Index
	data[Id] = fields.Id
	data[Result] = fields.Result
	data[Version] = fields.Version

	jsonData, err := json.Marshal(data);
	if err != nil {
		panic(err)
	}

	return jsonData
}

func getJsonFound(fields Fields) []byte {
	data := make(map[string]interface{})

	data[Took] = fields.Took
	data[TimedOut] = fields.Took

	jsonData, err := json.Marshal(data);
	if err != nil {
		panic(err)
	}

	return jsonData
}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
