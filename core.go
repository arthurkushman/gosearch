package gosearch

const (
	QUERY = "query"
	TERM = "term"
)

type Core interface {
	ParseInput(in map[string]interface{})
	SearchPhrase()
	BuildIndex()
	SetCanonicalIndex()
}

func ParseInput(in map[string]interface{}) {
	for k, v := range in {
		if k == QUERY && v != nil {

		}
	}
}