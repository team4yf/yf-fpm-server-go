package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/internal/db"

	_ "github.com/team4yf/yf-fpm-server-go/plugins/pg"
)

func TestPG(t *testing.T) {
	app := fpm.New()

	app.Init()

	dbclient, exists := app.GetDatabase("pg")
	assert.Equal(t, true, exists, "should true")

	err := dbclient.AutoMigrate(
		&Fake{},
	)

	assert.Nil(t, err, "should nil err")

	err = dbclient.Create(&Fake{
		Name:  "c",
		Value: 100,
	}).Error()

	assert.Nil(t, err, "should nil err")

	one := &Fake{}
	err = dbclient.Model(one).Condition("name = ?", "c").First(&one).Error()

	assert.Equal(t, 100, one.Value, "should be 100")

	rows := make([]*Fake, 0)
	dbclient.Model(one).Sorter(&db.Sorter{
		Sortby: "name",
		Asc:    "asc",
	}).Condition("name = ?", "c").Find(&rows).Error()

	assert.Equal(t, true, len(rows) > 0, "should more data")

}
