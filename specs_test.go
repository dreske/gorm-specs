package go_specs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type TestContext struct {
	Params map[string][]interface{}
}

func (t *TestContext) GetParameter(name string) ([]interface{}, bool) {
	values, ok := t.Params[name]
	return values, ok
}

type SpecsTestSuite struct {
	suite.Suite
	Db *gorm.DB
}

func (suite *SpecsTestSuite) SetupTest() {
	suite.Db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	suite.Db = suite.Db.Session(&gorm.Session{DryRun: true})
}

func (suite *SpecsTestSuite) TestBaseAnd() {
	t := suite.T()
	db := suite.Db

	spec1 := SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("first = ?", "test")
	})

	spec1 = spec1.And(SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("second = ?", "test")
	}))

	db = spec1.Apply(db, nil).Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE first = ? AND second = ?", db.Statement.SQL.String())
}

func (suite *SpecsTestSuite) TestBaseOr() {
	t := suite.T()
	db := suite.Db

	spec1 := SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("first = ?", "test")
	})

	spec1 = spec1.Or(SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("second = ?", "test")
	}))

	db = spec1.Apply(db, nil).Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE first = ? OR second = ?", db.Statement.SQL.String())
}

func (suite *SpecsTestSuite) TestBaseAndOr() {
	t := suite.T()
	db := suite.Db

	spec1 := SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("first = ?", "test")
	})

	spec1 = spec1.And(SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("second = ?", "test")
	}))

	spec1 = spec1.Or(SpecFunc(func(db *gorm.DB, ctx Context) *gorm.DB {
		return db.Where("third = ?", "test")
	}))

	db = spec1.Apply(db, nil).Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE first = ? AND second = ? OR third = ?", db.Statement.SQL.String())
}

func TestSpecsTestSuite(t *testing.T) {
	suite.Run(t, &SpecsTestSuite{})
}
