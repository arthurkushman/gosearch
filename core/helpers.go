package core

import (
	"time"
	"encoding/json"
)

func (sf *StoreFields) GetJsonOutput() []byte {
	if sf.Fld.OpType == ResultCreated {
		return getJsonCreated(sf.Fld)
	}

	if sf.Fld.OpType == ResultFound {
		return getJsonFound(sf)
	}

	return []byte{}
}

/** PUT
Example:
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

func getJsonCreated(fields Fields) []byte {
	data := make(map[string]interface{})

	data[ResultCreated] = true
	data[Took] = fields.Took
	data[Index] = fields.Index
	data[Type] = fields.IndexType
	data[Id] = fields.Id
	data[Result] = fields.Result
	data[Version] = fields.Version

	jsonData, err := json.Marshal(data);
	if err != nil {
		panic(err)
	}

	return jsonData
}

type HitsItem struct {
	Index     string                 `json:"_index" bson:"_index"`
	Type      string                 `json:"_type" bson:"_type"`
	Id        uint64                 `json:"_id" bson:"_id"`
	Timestamp uint64                 `json:"_timestamp" bson:"_timestamp"`
	Source    map[string]interface{} `json:"_source" bson:"_source"`
}

type HitsObject struct {
	Total     int      `json:"total" bson:"total"`
	HitsSlice []Fields `json:"hits" bson:"hits"` // slice of Fields which contain Source object
}

type FoundOutput struct {
	Took     int64 `json:"city_id" bson:"timed_out"`
	TimedOut bool  `json:"timed_out" bson:"timed_out"`
	HitsObject     `json:"hits" bson:"hits"`
}

/** POST/GET
Example:
{
    "took": 5,
    "timed_out": false,
    "hits": {
        "total": 1,
        "hits": [
            {
                "_index": "myindex",
                "_type": "myindextype",
                "_id": "1",
                "_timestamp": "1535187435",
                "_source": {
                    "title": "Lorem ipsum is a pseudo-Latin text used in web design",
                    "text": "Lorem ipsum is a pseudo-Latin text used in web design, typography, layout, and printing in place of English to emphasise design elements over content. It's also <tag1><tag2>called placeholder</tag1></tag2> (or filler) text. It's a convenient tool for mock-ups. It helps to outline the visual elements of a document or presentation, eg typography, font, or layout. Lorem ipsum is mostly a part of a Latin text by the classical author and philosopher Cicero. Its words and letters have been changed by addition or removal, so to deliberately render its content nonsensical; it's not genuine, correct, or comprehensible Latin anymore. While lorem ipsum's still resembles classical Latin, it actually has no meaning whatsoever. As Cicero's text doesn't contain the letters K, W, or Z, alien to latin, these, and others are often inserted randomly to mimic the typographic appearence of European languages, as are digraphs not to be found in the original.",
                    "data": "2017-08-21"
                }
            }
        ]
    }
}
 */
func getJsonFound(sf *StoreFields) []byte {
	data := &FoundOutput{Took: sf.Fld.Took, TimedOut: false,
		HitsObject: HitsObject{Total: len(sf.Collection), HitsSlice: sf.Collection}}

	jsonData, err := json.Marshal(data);
	if err != nil {
		panic(err)
	}

	return jsonData
}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
