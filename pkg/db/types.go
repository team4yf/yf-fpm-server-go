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

//BaseData basic data
type BaseData struct {
	AffectedRows int64
	Condition    string
	Table        string
	Model        interface{}
	Arguments    []interface{}
	Err          error
}

//QueryData query defination
type QueryData struct {
	Fields []interface{}
	*BaseData
	Pager  *Pagination
	Sorter []Sorter
}

//NewQuery set the query
func NewQuery() *QueryData {
	return &QueryData{
		BaseData: &BaseData{
			Condition: "1=1",
			Arguments: make([]interface{}, 0),
		},
		Pager: &Pagination{
			Skip:  0,
			Limit: -1,
		},
		Sorter: make([]Sorter, 0),
	}
}

func (p *BaseData) SetCondition(condition string, args ...interface{}) *BaseData {
	p.Condition = condition
	p.Arguments = args
	return p
}

func (p *QueryData) AddSorter(sorters ...Sorter) *QueryData {
	p.Sorter = append(p.Sorter, sorters...)
	return p
}

func (p *QueryData) SetPager(pager *Pagination) *QueryData {
	p.Pager = pager
	return p
}

func (p *QueryData) AddFields(args ...interface{}) *QueryData {
	p.Fields = append(p.Fields, args...)
	return p
}
func (p *BaseData) SetTable(t string) *BaseData {
	p.Table = t
	return p
}

func (p *BaseData) SetModel(t interface{}) *BaseData {
	p.Model = t
	return p
}
