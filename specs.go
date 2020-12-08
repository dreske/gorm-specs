package go_specs

import (
	"gorm.io/gorm"
)

type Context interface {
	GetParameter(name string) ([]string, bool)
}

type Spec interface {
	Apply(db *gorm.DB, ctx Context) *gorm.DB
	And(spec Spec) Spec
	Or(spec Spec) Spec
}

func (s *callbackSpec) And(spec Spec) Spec {
	return and(s, spec)
}

func (s *callbackSpec) Or(spec Spec) Spec {
	return or(s, spec)
}

type callbackSpec struct {
	Callback func(db *gorm.DB, ctx Context) *gorm.DB
}

func (s *callbackSpec) Apply(db *gorm.DB, ctx Context) *gorm.DB {
	return s.Callback(db, ctx)
}

func SpecFunc(fn func(db *gorm.DB, ctx Context) *gorm.DB) Spec {
	return &callbackSpec{
		Callback: fn,
	}
}
