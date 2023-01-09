package youtube

import ifs "what2cook/pkg/interfaces"

// *Video implements content interface
type Video struct {
	title       string
	description string
	id          string
	contentType ifs.ContentType
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

func (v *Video) Type() ifs.ContentType {
	return v.contentType
}
