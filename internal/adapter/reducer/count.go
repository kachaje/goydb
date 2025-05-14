package reducer

import (
	"github.com/kachaje/goydb/pkg/model"
)

type Count struct {
	result map[any]any
}

func NewCount() *Count {
	return &Count{
		result: make(map[any]any),
	}
}

func (r *Count) Reduce(doc *model.Document) {
	value, ok := r.result[doc.Key]

	if ok {
		r.result[doc.Key] = value.(int64) + 1
	} else {
		r.result[doc.Key] = int64(1)
	}
}

func (r *Count) Result() map[any]any {
	return r.result
}
