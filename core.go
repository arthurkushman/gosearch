package gosearch

const (
	QUERY = "query"
	TERM = "term"
)

type Core interface {
	ParseInput(in map[string]string)
	SearchPhrase()

}