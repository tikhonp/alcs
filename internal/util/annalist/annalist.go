package annalist

import (
	"log"

	"github.com/tikhonp/alcs/internal/util/assert"
)

// Annalists were a class of writers on Roman history,
// the period of whose literary activity lasted from the time
// of the Second Punic War to that of Sulla.

type Annalist interface {
	Log(msg ...any)
	Error(err error, msg ...any)
}

type AnnalistManager interface {
	GetAnnalist(tag string) Annalist
}

type defaultAnnalist struct {
	tag   string
	debug bool
}

func (a defaultAnnalist) Log(msg ...any) {
	log.Println(a.tag, msg)
}

func (a defaultAnnalist) Error(err error, msg ...any) {
	if a.debug {
		assert.NoError(err, a.tag, msg)
	} else {
		log.Println(err.Error(), a.tag, msg)
	}
}

type defaultAnnalistManager struct {
	debug bool
}

func (dam defaultAnnalistManager) GetAnnalist(tag string) Annalist {
	return &defaultAnnalist{tag: tag, debug: dam.debug}
}

func NewDefaultAnnalist(debug bool) AnnalistManager {
	return &defaultAnnalistManager{debug: debug}
}
