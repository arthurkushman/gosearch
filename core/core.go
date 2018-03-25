package core

import (
	"github.com/buger/jsonparser"
	"encoding/json"
)

const (
	QUERY = "query"
	TERM  = "term"
	// offset/limit
	OFFSET = "offset"
	LIMIT  = "limit"
	// highlight settings
	HIGHLIGHT      = "highlight"
	PreTags        = "pre_tags"
	PostTags       = "post_tags"
	HashIndexGlue  = ":"
	IdDocMatch     = "MATCHING"
	ReadBufferSize = 8094
)

type Core interface {
	ParseInput(in map[string]interface{})
	SearchPhrase()
	BuildIndex()
	SetCanonicalIndex()
	SearchById()
}

func (sf *StoreFields) SearchPhrase(fieldsValue map[string]string) {

}

/**
 * Searches document by uri routed ID
 */
func (sf *StoreFields) SearchById() {
	//SetIncrKey()
	//incrMatch = incrKey + HashIndexGlue + IdDocMatch
	//// get the document hash
	//$docSha = $this- > redisConn- > hget($incrMatch, $this- > id);
	//// get serialized data
	//$data = unserialize($this- > redisConn- > hget($this- > incrKey, $docSha));
	//if (empty($data)) {
	//throw new RequestException(Errors::REQUEST_MESSAGES[Errors::REQUEST_URI_DOC_ID_NOT_FOUND], Errors::REQUEST_URI_DOC_ID_NOT_FOUND);
	//}
	//$this- > stdFields- > setOpType(IndexInterface::RESULT_FOUND);
	//$this- > stdFields- > setOpStatus(true);
	//$source = $this- > unser($data[IndexInterface::SOURCE]);
	//$this- > stdFields- > setSource($source);
	//$this- > stdFields- > setId($data[IndexInterface::ID]);
	//$this- > stdFields- > setVersion($data[IndexInterface::VERSION]);
	//$this- > stdFields- > setTimestamp($data[IndexInterface::TIMESTAMP]);
}

//func (s *Storage) SetIncrKey() {
//	s.incrKey = $this-> Index. (empty($this-> IndexType) ? ''
//	: (self::HashIndexGlue . $this-> IndexType. ''))
//}

func ParseInput(data []byte) map[string]string {
	var fieldValueMap = make(map[string]string)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		json.Unmarshal(value, fieldValueMap)
	}, "query", "term")
	return fieldValueMap
}

func BuildIndex() {

}
