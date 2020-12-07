package go_specs

import (
	"gorm.io/gorm"
)

type Or struct {
	Specs []Spec
}

func or(specs ...Spec) Spec {
	return &Or{specs}
}

func (s *Or) Apply(db *gorm.DB, ctx Context) *gorm.DB {
	condition := db
	for i, spec := range s.Specs {
		if i == 0 {
			condition = spec.Apply(db, ctx)
		} else {
			condition = condition.Or(spec.Apply(db, ctx))
		}
	}

	return condition
}

func (s *Or) And(spec Spec) Spec {
	return and(s, spec)
}

func (s *Or) Or(spec Spec) Spec {
	return or(s, spec)
}
