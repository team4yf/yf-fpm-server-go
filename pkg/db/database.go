//Package db the db api
package db

//Database the interface of the db (default postgres)
type Database interface {
	AutoMigrate(...interface{}) error

	Find(QueryData, interface{}) error

	Count(QueryData, *int) error

	FindAndCount(QueryData, interface{}, *int) error

	First(QueryData, interface{}) error

	Create(BaseData, interface{}) error

	Remove(BaseData, *int) error

	Updates(BaseData, CommonMap, *int) error

	Execute(sql string, rows *int) error

	Raw(sql string, result interface{}) error

	Raws(sql string, iterator func() interface{}, appender func(interface{})) error

	Transaction(func(Database) error) error
}
