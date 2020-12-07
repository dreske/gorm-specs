package go_specs

import (
	"gorm.io/gorm"
)

type And struct {
	Specs []Spec
}

func and(specs ...Spec) Spec {
	return &And{specs}
}

func (s *And) And(spec Spec) Spec {
	return and(s, spec)
}

func (s *And) Or(spec Spec) Spec {
	return or(s, spec)
}

func (s *And) Apply(db *gorm.DB, ctx Context) *gorm.DB {
	conditions := make([]*gorm.DB, len(s.Specs))
	for i, spec := range s.Specs {
		conditions[i] = spec.Apply(db, ctx)
	}

	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db
}
