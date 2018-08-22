package core

import "time"

func (sf StoreFields) JsonOutput() {

}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//func GetJsonBody(buff []byte) map[string]string {
//
//}
//
//func GetJsonString() {
//
//}