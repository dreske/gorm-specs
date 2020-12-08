package go_specs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type OrTestSuite struct {
	suite.Suite
	Db *gorm.DB
}

func (suite *OrTestSuite) SetupTest() {
	suite.Db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	suite.Db = suite.Db.Session(&gorm.Session{DryRun: true})
}

func (suite *OrTestSuite) TestSimpleOr() {
	t := suite.T()
	db := suite.Db

	eq := &Or{[]Spec{
		&Arg{Param: "name", Path: []string{"firstname"}, Func: Equal},
		&Arg{Param: "name", Path: []string{"lastname"}, Func: Equal},
	},
	}
	db = eq.Apply(db, &TestContext{
		map[string][]string{"name": {"testname1"}},
	})
	db = db.Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE [firstname] = ? OR [lastname] = ?", db.Statement.SQL.String())
	assert.Equal(t, "testname1", db.Statement.Vars[0])
}

func (suite *OrTestSuite) TestAndOr() {
	t := suite.T()
	db := suite.Db

	eq := &And{[]Spec{
		&Arg{Param: "name", Path: []string{"first"}, Func: Equal},
		&Or{[]Spec{
			&Arg{Param: "name", Path: []string{"second"}, Func: Equal},
			&Arg{Param: "name", Path: []string{"third"}, Func: Equal},
		}},
	},
	}

	//db = db.Where("id = ?", 1).Where(db.Where("firstname = ?", "firstname").Or("lastname = ?", "lastname"))
	db = eq.Apply(db, &TestContext{
		map[string][]string{"name": {"testname1"}},
	})
	db = db.Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE [first] = ? AND ([second] = ? OR [third] = ?)", db.Statement.SQL.String())
}

func (suite *OrTestSuite) TestAndOrAnd() {
	t := suite.T()
	db := suite.Db

	eq := &And{[]Spec{
		&Arg{Param: "name", Path: []string{"first"}, Func: Equal},
		&Or{[]Spec{
			&Arg{Param: "name", Path: []string{"second"}, Func: Like},
			&And{[]Spec{
				&Arg{Param: "name", Path: []string{"third"}, Func: Equal},
				&Arg{Param: "name", Path: []string{"fourth"}, Func: Equal},
			}},
		}},
	},
	}

	//db = db.Where("id = ?", 1).Where(db.Where("firstname = ?", "firstname").Or("lastname = ?", "lastname"))
	db = eq.Apply(db, &TestContext{
		map[string][]string{"name": {"test"}, "id": {"1"}},
	})
	db = db.Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE [first] = ? AND ([second] LIKE ? OR ([third] = ? AND [fourth] = ?))", db.Statement.SQL.String())
}

func TestOrTestSuite(t *testing.T) {
	suite.Run(t, &OrTestSuite{})
}
