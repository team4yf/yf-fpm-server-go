package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/yf-fpm-server-go/plugins/pg"
)

type countBody struct {
	C float64 `json:"c"`
	B float64 `json:"b"`
}

func TestPG(t *testing.T) {
	app := fpm.New()

	app.Init()

	dbclient, exists := app.GetDatabase("pg")
	assert.Equal(t, true, exists, "should true")

	//Test AutoMigrate
	err := dbclient.AutoMigrate(
		&Fake{},
	)

	assert.Nil(t, err, "should nil err")

	//Test Create
	err = dbclient.Create(&Fake{
		Name:  "c",
		Value: 100,
	}).Error()

	assert.Nil(t, err, "should nil err")

	//Test First, TODO: Sort, Skip, Limit
	one := &Fake{}
	err = dbclient.Model(one).Condition("name = ?", "c").First(&one).Error()

	assert.Equal(t, 100, one.Value, "should be 100")

	//Test Count
	total := 0
	err = dbclient.Model(Fake{}).Condition("name = ?", "c").Count(&total).Error()
	assert.Nil(t, err, "should nil err")
	assert.Equal(t, true, total > 0, "should gt 0")

	//Test Remove
	rows := 0
	err = dbclient.Model(Fake{}).Condition("name = ?", "c").Remove(&rows).Error()
	assert.Nil(t, err, "should nil err")
	assert.Equal(t, true, rows > 0, "should gt 0")

	//Test Execute
	err = dbclient.Execute(`delete from fake where id = 11`, &rows).Error()
	assert.Nil(t, err, "should not error")
	assert.Equal(t, true, rows >= 0, "should gt 0")

	//Test Raw
	raw := &countBody{}
	err = dbclient.Raw(`select count(1) as c, 1 as b from fake where id < 10`, raw).Error()
	assert.Nil(t, err, "should not error")
	assert.Equal(t, true, raw.C >= 0, "should gt 0")
	assert.Equal(t, true, raw.B == 1, "should eq 1")

	//Test Raws
	raws := make([]countBody, 0)
	err = dbclient.Raws(`select id as c, 1 as b from fake where id < 10`, &raws, func() interface{} {
		return &countBody{}
	}).Error()
	assert.Nil(t, err, "should not error")
	assert.Equal(t, true, raw.C >= 0, "should gt 0")
	assert.Equal(t, true, raw.B == 1, "should eq 1")

	// rows := make([]*Fake, 0)
	// dbclient.Model(one).Sorter(db.Sorter{
	// 	Sortby: "name",
	// 	Asc:    "asc",
	// }).Condition("name = ?", "c").Find(&rows).Error()

	// assert.Equal(t, true, len(rows) > 0, "should more data")

}
