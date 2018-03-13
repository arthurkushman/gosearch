package gosearch

import "encoding/json"

type Serialize interface {
	Ser(data map[string]*json.RawMessage) []byte
	Unser(data []byte) map[string]*json.RawMessage
}

func Ser(data map[string]*json.RawMessage) []byte {
	
}

func Unser(data []byte) map[string]*json.RawMessage {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(data, &objmap)
	if err != nil {
		panic(err)
	}
	return objmap
}
