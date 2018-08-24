package core

import "encoding/json"

type Serialize interface {
	Ser(data map[string]*json.RawMessage) []byte
	Unser(data []byte) map[string]*json.RawMessage
}

func Ser(data map[string]interface{}) []byte {
	objData, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	return objData
}

func Unser(data []byte) map[string]interface{} {
	var objmap map[string]interface{}
	err := json.Unmarshal(data, &objmap)

	if err != nil {
		panic(err)
	}

	return objmap
}
