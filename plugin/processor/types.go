package processor

import (
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

type Type string

const (
	Pre   Type = "pre"
	Intra Type = "intra"
	Post  Type = "post"
)

// Types is a convenience for iterating all of the processor types.
// The order of this slice is intentional!
var Types []Type = []Type{Pre, Intra, Post}

type Processor interface {
	Init() error
	Process(<-chan types.Media, chan<- types.Media)
}

type ProcessorFunc func(<-chan types.Media, chan<- types.Media)

func (ProcessorFunc) Init() error {
	return nil
}

func (pf ProcessorFunc) Process(in <-chan types.Media, out chan<- types.Media) {
	pf(in, out)
}

func AppendFunc(ps []Processor, fs ...ProcessorFunc) []Processor {
	pfs := make([]Processor, len(fs))
	for i := range fs {
		pfs[i] = fs[i]
	}
	return append(ps, pfs...)
}

var Registry map[Type]map[string](func() Processor) = map[Type]map[string](func() Processor){
	Pre:   map[string](func() Processor){},
	Intra: map[string](func() Processor){},
	Post:  map[string](func() Processor){},
}

func Register(t Type, name string, initializer func() Processor) {
	if _, ok := Registry[t][name]; ok {
		log.Fatalf("processor registry already contains plugin named %s", name)
	}
	Registry[t][name] = initializer
}
