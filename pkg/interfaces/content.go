package interfaces

type ContentType int

const (
	// MUST BE UNIQUE
	YTVideoType ContentType = iota
	IGVideoType
)

type Content interface {
	Id() string
	Description() string
	Title() string
	// Do not serialize Type - it is not guaranteed to remain the same
	Type() ContentType
}

// catalog, update
