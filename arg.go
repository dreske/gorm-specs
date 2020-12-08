package go_specs

import (
	"fmt"
	"gorm.io/gorm"
)

type ArgFunc func(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB

type ArgType uint

const (
	String ArgType = iota
	Integer
	Float
)

type Arg struct {
	Param string
	Path  []string
	Type  ArgType
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
	if values, ok := getParameter(ctx, spec.Param, spec.Type); ok {
		return db.Where(fmt.Sprintf("%s = ?", spec.Path), values[0])
	}
	return db
}

func Like(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB {
	if values, ok := getParameter(ctx, spec.Param, String); ok {
		return db.Where(fmt.Sprintf("%s LIKE ?", spec.Path), "%"+values[0].(string)+"%")
	}
	return db
}

func LikeIgnoreCase(db *gorm.DB, spec *Arg, ctx Context) *gorm.DB {
	if values, ok := getParameter(ctx, spec.Param, String); ok {
		return db.Where(fmt.Sprintf("%s ILIKE ?", spec.Path), "%"+values[0].(string)+"%")
	}
	return db
}

func getParameter(ctx Context, name string, target ArgType) ([]interface{}, bool) {
	values, ok := ctx.GetParameter(name)
	if !ok {
		return nil, false
	}

	result := make([]interface{}, len(values))
	switch target {
	case String:
		for i, val := range values {
			result[i] = val
		}
		return result, true
	default:
		return nil, false
	}
}
