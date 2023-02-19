package interfaces

import w2d_util "what2cook/pkg/util"

type Source interface {
	Init() error
	GetCreator(string, w2d_util.SearchParameters) (Creator, error)
	// Fetch `n` content(s) for given creator
	FetchContent(Creator, int) ([]Content, error)
}

// init, update
