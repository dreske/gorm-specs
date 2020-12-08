package go_specs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type TestModel struct {
	ID uint `gorm:"primarykey"`
}

type AgTestSuite struct {
	suite.Suite
	Db *gorm.DB
}

func (suite *AgTestSuite) SetupTest() {
	suite.Db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	suite.Db = suite.Db.Session(&gorm.Session{DryRun: true})
}

func (suite *AgTestSuite) TestEqual() {
	t := suite.T()
	db := suite.Db

	eq := &Arg{Param: "name", Path: []string{"name"}, Type: String, Func: Equal}
	db = eq.Apply(db, &TestContext{
		map[string][]string{"name": {"testname1"}},
	})
	db = db.Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE [name] = ?", db.Statement.SQL.String())
	assert.Equal(t, "testname1", db.Statement.Vars[0])
}

func (suite *AgTestSuite) TestLike() {
	t := suite.T()
	db := suite.Db

	eq := &Arg{Param: "name", Path: []string{"name"}, Func: Like}
	db = eq.Apply(db, &TestContext{
		map[string][]string{"name": {"testname1"}},
	})
	db = db.Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE [name] LIKE ?", db.Statement.SQL.String())
	assert.Equal(t, "%testname1%", db.Statement.Vars[0])
}

func (suite *AgTestSuite) TestLikeIgnoreCase() {
	t := suite.T()
	db := suite.Db

	eq := &Arg{Param: "name", Path: []string{"name"}, Func: LikeIgnoreCase}
	db = eq.Apply(db, &TestContext{
		map[string][]string{"name": {"testname1"}},
	})
	db = db.Find(&TestModel{})

	assert.Equal(t, "SELECT * FROM `test_models` WHERE [name] ILIKE ?", db.Statement.SQL.String())
	assert.Equal(t, "%testname1%", db.Statement.Vars[0])
}

func TestEqualTestSuite(t *testing.T) {
	suite.Run(t, &AgTestSuite{})
}
