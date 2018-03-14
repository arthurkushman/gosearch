package gosearch

type RequestHandler interface {
	GetLimit() int
	GetOffset() int
	GetPreTags() string
	GetPostTags() string
	IsHighlighted() bool
	GetHighlightedFields() []string
	GetRequestMethod() string
	GetRequestBody() map[string]string
}

