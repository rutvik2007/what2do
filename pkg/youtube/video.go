package youtube

import (
	"encoding/json"
	"time"
	ifs "what2cook/pkg/interfaces"
)

// *Video implements content interface
type Video struct {
	title        string
	description  string
	id           string
	creationTime time.Time
	contentType  ifs.ContentType
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

func (v *Video) CreatedAt() time.Time {
	return v.creationTime
}

func (v *Video) MarshalJSON() ([]byte, error) {
	v_map := make(map[string]string)
	v_map["description"] = v.description
	v_map["content_type"] = "youtubeVideo"
	v_map["id"] = v.id
	v_map["creation_time"] = v.creationTime.String()

	return json.Marshal(v_map)
}
