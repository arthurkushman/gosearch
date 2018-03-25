package core

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

//func GetLimit() int {
//
//}
//
//func GetOffset() int {
//
//}
//
//func GetPreTags() string {
//
//}
//
//func GetPostTags() string {
//
//}
//
//func IsHighlighted() bool {
//
//}
//
//func GetHighlightedFields() []string {
//
//}
//
//func GetRequestMethod() string {
//
//}
//
//func GetRequestBody() map[string]string {
//
//}