package reducer

import (
	"github.com/kachaje/goydb/pkg/model"
)

type None struct {
	result map[any]any
}

func NewNone() *Count {
	return &Count{
		result: make(map[any]any),
	}
}

func (r *None) Reduce(doc *model.Document) {
	r.result[doc.ID] = doc
}

func (r *None) Result() map[any]any {
	return r.result
}
