package youtube

import "what2cook/pkg/interfaces"

// *Video implements content interface
type Video struct {
	title       string
	description string
	id          string
	contentType interfaces.ContentType
}

func (v *Video) Id() string {
	return v.id
}

func (v *Video) Description() string {
	return v.description
}

func (v *Video) Title() string {
	return v.title
}

func (v *Video) Type() interfaces.ContentType {
	return v.contentType
}
