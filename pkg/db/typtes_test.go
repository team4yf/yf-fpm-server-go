package db

import "testing"

func TestQuery(t *testing.T) {
	q := NewQuery()
	q.SetCondition("a = ?", "bv").SetTable("crdget")

}
