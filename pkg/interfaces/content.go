package interfaces

type ContentType int

const (
	// Content Types MUST BE UNIQUE
	VideoType ContentType = 0
)

type Content interface {
	Id() string
	Description() string
	Title() string
	Type() ContentType
}

// catalog, update
