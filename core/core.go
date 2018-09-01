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
	"strings"
	"crypto/md5"
)

const (
	Query = "query"
	Term  = "term"

	// offset/limit
	Offset = "offset"
	Limit  = "limit"

	// highlight settings
	Highlight      = "highlight"
	PreTags        = "pre_tags"
	PostTags       = "post_tags"
	HashIndexGlue  = ":"
	ListIndexGlue  = "___"
	IdDocMatch     = "MATCHING"
	ReadBufferSize = 8094

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
	Took     = "took"
	TimedOut = "timed_out"
	Hits     = "hits"
	Total    = "total"
	Indices  = "indices"
	// system reserved keywords
	Index     = "_index"
	Type      = "_type"
	Source    = "_source"
	Id        = "_id"
	Timestamp = "_timestamp"
	Version   = "_version"
	All       = "_all"
	Result    = "result"
	Document  = "_document"
	Cat       = "_cat"
	Reindex   = "_reindex"

	// op results
	ResultDeleted  = "deleted"
	ResultCreated  = "created"
	ResultUpdated  = "updated"
	ResultFound    = "found"
	ResultNotFound = "not found"
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
	OpType   string
	OpStatus bool
	Result   string
	// source indexed document
	Source        map[string]interface{}
	Id            uint64 // doc id
	Index         string
	IndexType     string
	Version       uint64
	Timestamp     uint64
	IsSearch      bool
	Took          int64
	RequestSource []byte
}

type Storage struct {
	redis        redis.Conn
	IncrKey      string
	HashIndexKey string
	ListIndexKey string
	WordHashes   map[string]uint8
	DocHashes    map[string]uint8
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
 * Searches document by uri routed Id
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
			sf.Err.ErrMsg = "Doc Id not found"
			EchoError(w, HttpEror400, sf.Err)
		}

		sf.Fld.OpType = ResultFound
		sf.Fld.OpStatus = true
		sf.Fld.Source = data[Source].(map[string]interface{})
		sf.Fld.Id = data[Id].(uint64)
		sf.Fld.Version = data[Version].(uint64)
		sf.Fld.Timestamp = data[Timestamp].(uint64)
	}
}

func (sf *StoreFields) SetIncrKey() {
	var idxTypeGlue = ""
	if sf.Fld.IndexType != "" {
		idxTypeGlue = HashIndexGlue + sf.Fld.IndexType
	}
	sf.Stg.IncrKey = sf.Fld.Index + idxTypeGlue
}

func (sf *StoreFields) SetHashIndexKey() {
	var idxTypeGlue = HashIndexGlue
	if sf.Fld.IndexType != "" {
		idxTypeGlue = HashIndexGlue + sf.Fld.IndexType + HashIndexGlue
	}

	sf.Stg.HashIndexKey = sf.Fld.Index + idxTypeGlue
}

func (sf *StoreFields) SetListIndexKey() {
	var idxTypeGlue = ListIndexGlue
	if sf.Fld.IndexType != "" {
		idxTypeGlue = ListIndexGlue + sf.Fld.IndexType + ListIndexGlue
	}

	sf.Stg.ListIndexKey = sf.Fld.Index + idxTypeGlue
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

func EchoResult(w http.ResponseWriter, data []byte, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write(data)
}

func composeError(err Error) []byte {
	buff := bytes.NewBufferString("{")
	buff.WriteString("\"code\": \"" + err.ErrCode + "\"")
	buff.WriteString("\"message\": \"" + err.ErrMsg + "\"")
	buff.WriteString("}")

	return buff.Bytes()
}

func (sf *StoreFields) GetDocInfo() (reply interface{}, err error) {
	sourceStr := Ser(sf.Fld.Source)
	docSha := sha1.Sum([]byte(sourceStr))

	return sf.Stg.redis.Do(sf.Stg.IncrKey, docSha)
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
	for field, value := range sf.Fld.Source {
		words := strings.Fields(value.(string))
		for _, word := range words {
			sf.insertWord(field, word)
		}
	}
}

func (sf *StoreFields) insertWord(field, word string) {
	wordHash := md5.Sum([]byte(field + word))
	fmt.Println(sf.Stg.WordHashes[string(wordHash[:md5.Size])])

	// to avoid duplicate indexing
	if _, ok := sf.Stg.WordHashes[string(wordHash[:md5.Size])]; !ok {
		sf.Stg.WordHashes[string(wordHash[:md5.Size])] = 1;
		if sf.Fld.Id == 0 {
			sf.setIndexData()
		}
	}
}

func (sf *StoreFields) setIndexData() {

}

func (sf *StoreFields) setRequestDocument() {

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
