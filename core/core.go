package core

import (
	"github.com/buger/jsonparser"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
	"flag"
	"net/http"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
)

const (
	Query = "query"
	Term  = "term"
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
	ResultFound    = "found"
	// canonical values
	FIELDS          = "fields"
	PROPERTIES      = "properties"
	STRUCTURE       = "structure"
	ALIASES         = "aliases"
	MAPPINGS        = "mappings"
	FIELD_TYPE      = "type"
	IGNORE_ABOVE    = "ignore_above"
	DATA_SOURCE     = "source"
	DATA_DEST       = "dest"
	DATA_INDEX      = "index"
	DATA_INDEX_TYPE = "index_type"
)

// index specific constants
const (
	TOOK      = "took"
	TIMED_OUT = "timed_out"
	HITS      = "hits"
	TOTAL     = "total"
	INDICES   = "indices"
	// system reserved keywords
	INDEX     = "_index"
	TYPE      = "_type"
	SOURCE    = "_source"
	ID        = "_id"
	TIMESTAMP = "_timestamp"
	VERSION   = "_version"
	ALL       = "_all"
	RESULT    = "result"
	DOCUMENT  = "_document"
	CAT       = "_cat"
	REINDEX   = "_reindex"
	// op results
	RESULT_DELETED   = "deleted"
	RESULT_CREATED   = "created"
	RESULT_UPDATED   = "updated"
	RESULT_FOUND     = "found"
	RESULT_NOT_FOUND = "not found"
)

const (
	HttpEror400          = 400
	ErrCodeDocIdNotFound = "101"
)

type Core interface {
	ParseInput(in map[string]interface{})
	SearchPhrase()
	BuildIndex()
	SetCanonicalIndex()
	SearchById()
}

type Fields struct {
	OpType        string
	OpStatus      bool
	Result        string
	Source        string // json source
	Id            uint64 // doc id
	Index         string
	IndexType     string
	IncrKey       string
	Version       uint64
	Timestamp     uint64
	IsSearch      bool
	Took          int64
	RequestSource []byte
}

type Storage struct {
	redis   redis.Conn
	IncrKey string
}

type Error struct {
	ErrCode string
	ErrMsg  string
}

type StoreFields struct {
	Fld        Fields
	Stg        Storage
	Err        Error
	Collection []Fields // to collect multiple docs
}

func (sf *StoreFields) SearchPhrase(fieldsValue map[string]string) {

}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

var (
	pool        *redis.Pool
	redisServer = flag.String("127.0.0.1", ":6379", "")
)

func (sf *StoreFields) redisConn() {
	flag.Parse()
	pool = newPool(*redisServer)
	sf.Stg.redis = pool.Get()
}

/**
 * Searches document by uri routed ID
 */
func (sf *StoreFields) SearchById(w http.ResponseWriter) {
	sf.SetIncrKey()
	incrMatch := sf.Stg.IncrKey + HashIndexGlue + IdDocMatch
	// get the document hash
	//docSha = $this- > redisConn- > hget(incrMatch, sf.Fld.Id)
	sf.redisConn()
	defer sf.Stg.redis.Close()
	docSha, err := sf.Stg.redis.Do("hget", incrMatch, sf.Fld.Id)
	if err == nil {
		// get serialized data
		docData, _ := redis.Bytes(sf.Stg.redis.Do("hget", sf.Stg.IncrKey, docSha))
		data := Unser(docData)
		if len(data) == 0 {
			sf.Err.ErrCode = ErrCodeDocIdNotFound
			sf.Err.ErrMsg = "Doc ID not found"
			EchoError(w, HttpEror400, sf.Err)
		}

		sf.Fld.OpType = ResultFound
		sf.Fld.OpStatus = true
		sf.Fld.Source = data[SOURCE].(string)
		sf.Fld.Id = data[ID].(uint64)
		sf.Fld.Version = data[VERSION].(uint64)
		sf.Fld.Timestamp = data[TIMESTAMP].(uint64)
	}
}

func (sf *StoreFields) SetIncrKey() {
	var idxTypeGlue = ""
	if sf.Fld.IndexType != "" {
		idxTypeGlue = HashIndexGlue + sf.Fld.IndexType
	}
	sf.Stg.IncrKey = sf.Fld.Index + idxTypeGlue
}

func ParseInput(data []byte) map[string]string {
	var fieldValueMap = make(map[string]string)
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		json.Unmarshal(value, fieldValueMap)
	}, Query, Term)

	return fieldValueMap
}

func EchoError(w http.ResponseWriter, errCode int, err Error) {
	w.WriteHeader(errCode)
	buff := composeError(err)
	w.Write(buff)
}

func composeError(err Error) []byte {
	buff := bytes.NewBufferString("{")
	buff.WriteString("\"code\": \"" + err.ErrCode + "\"")
	buff.WriteString("\"message\": \"" + err.ErrMsg + "\"")
	buff.WriteString("}")
	return buff.Bytes()
}

func (sf *StoreFields) GetDocInfo() (reply interface{}, err error) {
	docSha := sha1.Sum([]byte(sf.Fld.Source))
	return sf.Stg.redis.Do(sf.Fld.IncrKey, docSha)
}

func (sf *StoreFields) SetCanonicalIndex() {
	docSha, err := sf.Stg.redis.Do("hget", sf.Fld.Index, STRUCTURE)
 	var data interface{}
	if err == nil && docSha != nil {
		data = docSha
	} else {
		//data := map[string]map[string]interface{
			//sf.Fld.Index : map[string]string{
			//	ALIASES : interface{}
			//},
		//}
	}
	fmt.Print(data)
	sf.Stg.redis.Do("hset", sf.Fld.Index, STRUCTURE, data)
}

func (sf *StoreFields) Insert() {

}

func (sf *StoreFields) SetMappings() {

}

/**
 *  Sets the only source doc from input stream
 */
func (sf *StoreFields) SetSourceDocument(r *http.Request) {
	sf.Fld.RequestSource = Ser(sf.ReadJsonBody(r))
}

/**
 * Reads json body form request
 */
func (sf *StoreFields) ReadJsonBody(r *http.Request) map[string]interface{} {
	buf := make([]byte, ReadBufferSize)
	n, err := r.Body.Read(buf)
	jsonBytes := buf[:n]

	if (err != nil && err != io.EOF) || n > ReadBufferSize {
		panic("Error reading from input stream: " + err.Error())
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		panic(err)
	}

	return data
}
