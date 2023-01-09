package interfaces

// typically like this:
// Id     string;
// Name   string
// Content []Content

type Creator interface {
	Marshal() []byte
	// Unmarshal([]byte)
}

// update catalog
//
