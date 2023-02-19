package youtube

// implements creator interface
type Channel struct {
	id   string
	name string
}

func (ch *Channel) Marshal() []byte {
	return []byte{}
}
