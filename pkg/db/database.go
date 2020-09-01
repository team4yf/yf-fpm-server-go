//Package db the db api
package db

//Database the interface of the db (default postgres)
type Database interface {
	AutoMigrate(...interface{}) error

	Find(*QueryData, interface{}) error

	Count(*BaseData, *int64) error

	FindAndCount(*QueryData, interface{}, *int64) error

	First(*QueryData, interface{}) error

	Create(*BaseData, interface{}) error

	Remove(*BaseData, *int64) error

	Updates(*BaseData, CommonMap, *int64) error

	Execute(sql string, rows *int64) error

	Raw(sql string, result interface{}) error

	Raws(sql string, iterator func() interface{}, appender func(interface{})) error

	Transaction(func(Database) error) error
}
