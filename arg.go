package go_specs

import (
	"fmt"
	"gorm.io/gorm"
)

type ArgFunc func(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB

type Arg struct {
	Param string
	Path  []string
	Func  ArgFunc
}

func (s *Arg) Apply(db *gorm.DB, ctx Context) *gorm.DB {
	return s.Func(db, s, ctx)
}

func (s *Arg) And(spec Spec) Spec {
	return and(s, spec)
}

func (s *Arg) Or(spec Spec) Spec {
	return or(s, spec)
}

func Equal(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB {
	if values, ok := ctx.GetParameter(spec.Param); ok {
		return db.Where(fmt.Sprintf("%s = ?", spec.Path), values[0])
	}
	return db
}

func Like(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB {
	if values, ok := ctx.GetParameter(spec.Param); ok {
		value := ""
		switch v := values[0].(type) {
		case string:
			value = v
		}

		return db.Where(fmt.Sprintf("%s LIKE ?", spec.Path), "%"+value+"%")
	}
	return db
}

func LikeIgnoreCase(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB {
	if values, ok := ctx.GetParameter(spec.Param); ok {
		value := ""
		switch v := values[0].(type) {
		case string:
			value = v
		}

		return db.Where(fmt.Sprintf("%s ILIKE ?", spec.Path), "%"+value+"%")
	}
	return db
}
