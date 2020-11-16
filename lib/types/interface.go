package types

// Content is an interface that spans workshops downh to subsub
type Content interface {
	Compare(Content) (error, []string)
	GetTitle() string
}
