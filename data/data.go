package data

import (
	"reflect"
)

// GetEqualIndexOf fromからvalと一致したindexを返す.
// 見つからなかったときは-1を返す
func GetEqualIndexOf(from []interface{}, val interface{}) int {
	for i, f := range from {
		res := reflect.DeepEqual(f, val)
		if res {
			return i
		}
	}
	return -1
}
