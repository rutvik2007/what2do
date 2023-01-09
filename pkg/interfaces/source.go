package interfaces

import w2d_util "what2cook/pkg/util"

type Source interface {
	Init() error
	GetCreator(string, w2d_util.SearchParameters) (Creator, error)
}

// init, update
