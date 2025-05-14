package reducer

import (
	"github.com/kachaje/goydb/pkg/model"
)

type Sum struct {
	result map[any]any
}

func NewSum() *Sum {
	return &Sum{
		result: make(map[any]any),
	}
}

// Reduce will sum up using int64 if integer values are used,
// and switch to float64 as soon as decimal values are used.
func (r *Sum) Reduce(doc *model.Document) {
	value, ok := r.result[doc.Key]

	if ok {
		if cur, ok := value.(int64); ok {
			if add, ok := doc.Value.(int64); ok {
				r.result[doc.Key] = cur + add
			}
			if add, ok := doc.Value.(float64); ok {
				// switch to decimal value
				r.result[doc.Key] = float64(cur) + add
			}
		}
		if cur, ok := value.(float64); ok {
			if add, ok := doc.Value.(int64); ok {
				r.result[doc.Key] = cur + float64(add)
			}
			if add, ok := doc.Value.(float64); ok {
				r.result[doc.Key] = cur + add
			}
		}
	} else {
		r.result[doc.Key] = doc.Value
	}
}

func (r *Sum) Result() map[any]any {
	return r.result
}
