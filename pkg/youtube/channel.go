package youtube

import (
	ifs "what2cook/pkg/interfaces"
)

// implements creator interface
type Channel struct {
	id     string
	name   string
	videos []ifs.Content
}

func (ch *Channel) Marshal() []byte {
	return []byte{}
}
