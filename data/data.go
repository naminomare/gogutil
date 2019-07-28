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

// GetEqualLastIndexOf fromから逆順にvalと最初に一致したindexを返す.
// 見つからなかったときは-1を返す
func GetEqualLastIndexOf(from []interface{}, val interface{}) int {
	for a := range from {
		i := len(from) - a - 1
		f := from[i]
		res := reflect.DeepEqual(f, val)
		if res {
			return i
		}
	}
	return -1
}
