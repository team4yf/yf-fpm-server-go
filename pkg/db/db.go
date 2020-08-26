//Package db the db api
package db

//Pagination pager tool
type Pagination struct {
	Limit int
	Skip  int
}

//Sorter the sort condition
type Sorter struct {
	Sortby string
	Asc    string
}

//CommonMap the common map for database
type CommonMap map[string]interface{}

//Database the interface of the db (default postgres)
type Database interface {
	AutoMigrate(tables ...interface{}) error

	Condition(condition string, args ...interface{}) Database

	Sorter(...Sorter) Database

	Pager(*Pagination) Database

	Model(model interface{}) Database

	Error() error

	Find(result interface{}) Database

	Count(total *int) Database

	FindAndCount(result interface{}, total *int) Database

	First(result interface{}) Database

	Create(entity interface{}) Database

	Remove(total *int) Database

	Updates(updates CommonMap, rows *int) Database

	Execute(sql string, rows *int) Database

	Raw(sql string, result interface{}) Database

	Raws(sql string, iterator func() interface{}, appender func(interface{})) Database

	Transaction(func(Database) error) error
}
