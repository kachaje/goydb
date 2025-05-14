package port

import "github.com/kachaje/goydb/pkg/model"

type ReducerEngines map[string]ReducerServerBuilder

type ReducerServerBuilder func(fn string) (Reducer, error)

type Reducer interface {
	Reduce(doc *model.Document)
	Result() map[any]any
}
