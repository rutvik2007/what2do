package interfaces

// typically like this:
// Id     string;
// Name   string
// Content []Content

type Creator interface {
	FetchContent(Source) ([]Content, error)
}

// update catalog
//
