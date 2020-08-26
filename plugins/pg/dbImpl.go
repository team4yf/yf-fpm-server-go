package pg

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/team4yf/yf-fpm-server-go/internal/db"
)

type queryData struct {
	condition string
	arguments []interface{}
	pager     *db.Pagination
	sorter    []*db.Sorter
	err       *error
	model     interface{}
}

func newQuery() *queryData {
	return &queryData{
		condition: "1=1",
		arguments: make([]interface{}, 0),
		pager: &db.Pagination{
			Skip:  0,
			Limit: 20,
		},
		sorter: make([]*db.Sorter, 0),
	}
}

type pgImpl struct {
	db *gorm.DB
	q  *queryData
}

//New create a new instance
func New(setting *DBSetting) db.Database {
	db := CreateDb(setting)
	return &pgImpl{
		db: db,
	}
}

//CreateDb create new instance
func CreateDb(setting *DBSetting) *gorm.DB {
	//use the config for the app
	dsn := getDbEngineDSN(setting)
	db, err := gorm.Open(setting.Engine, dsn)
	if err != nil {
		panic(err)
	}

	db.DB().SetConnMaxLifetime(time.Minute * 5)
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(500)

	db.LogMode(setting.ShowSQL)

	return db
}

// 获取数据库引擎DSN  mysql,sqlite,postgres
func getDbEngineDSN(db *DBSetting) string {
	engine := strings.ToLower(db.Engine)
	dsn := ""
	switch engine {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowNativePasswords=true",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database,
			db.Charset)
	case "postgres":
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database)
	}

	return dsn
}

//AutoMigrate migrate table from the model
func (p *pgImpl) AutoMigrate(tables ...interface{}) (err error) {
	return p.db.AutoMigrate(tables...).Error
}

func (p *pgImpl) Condition(condition string, args ...interface{}) db.Database {
	p.q.condition = condition
	p.q.arguments = args
	return p
}

func (p *pgImpl) Sorter(sorters ...*db.Sorter) db.Database {
	p.q.sorter = sorters
	return p
}

func (p *pgImpl) Pager(pager *db.Pagination) db.Database {
	p.q.pager = pager
	return p
}

func (p *pgImpl) Model(model interface{}) db.Database {
	p.q = newQuery()
	p.q.model = model
	return p
}

func (p *pgImpl) Error() (err error) {
	return *p.q.err
}

func (p *pgImpl) Find(result interface{}) db.Database {
	//TODO sort & skip & check the result point
	p.q.err = &p.db.Model(p.q.model).Where(p.q.condition, p.q.arguments).Limit(p.q.pager.Limit).Find(&result).Error

	return p
}

func (p *pgImpl) Count(total *int) db.Database {
	p.q.err = &p.db.Model(p.q.model).Where(p.q.condition, p.q.arguments).Count(&total).Error
	return p
}

func (p *pgImpl) FindAndCount(result []interface{}, total *int) db.Database {
	return p
}

func (p *pgImpl) First(result interface{}) db.Database {
	p.q.err = &p.db.Model(p.q.model).Where(p.q.condition, p.q.arguments).First(result).Error
	return p
}

func (p *pgImpl) Create(entity interface{}) db.Database {
	if p.q == nil {
		p.q = newQuery()
		p.q.model = entity
	}
	p.q.err = &p.db.Create(entity).Error
	return p
}

func (p *pgImpl) Remove(total *int) db.Database {
	p.q.err = &p.db.Where(p.q.condition, p.q.arguments).Delete(p.q.model).Error
	return p
}

func (p *pgImpl) Updates(updates db.CommonMap) db.Database {
	p.q.err = &p.db.Model(p.q.model).Where(p.q.condition, p.q.arguments).Updates(updates).Error
	return p
}

func (p *pgImpl) Execute(sql string) db.Database {
	p.q.err = &p.db.Exec(sql).Error
	return p
}

func (p *pgImpl) Raw(sql string, result interface{}) db.Database {
	raw := p.db.Raw(sql)
	if raw.Error != nil {
		p.q.err = &raw.Error
		return p
	}
	p.q.err = &raw.Scan(&result).Error

	return p
}

func (p *pgImpl) Raws(sql string, results []interface{}, iterator func(interface{}, interface{}) error) db.Database {
	raws, err := p.db.Raw(sql).Rows()
	if err != nil {
		p.q.err = &err
		return p
	}
	defer raws.Close()
	for raws.Next() {
		//TODO:
		// p.q.err = &rows.Scan().Error

	}

	return p
}
