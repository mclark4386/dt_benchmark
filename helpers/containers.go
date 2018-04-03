package helpers

import (
	"sort"
)

func ArrayContains(arr []interface{}, item interface{}) bool {
	for _, i := range arr {
		if i == item {
			return true
		}
	}
	return false
}

func ArrayAppend(arr []interface{}, item interface{}) []interface{} {
	arr2 := append(arr, item)
	return arr2
}

func StringArrayContains(arr []string, item string) bool {
	i := sort.SearchStrings(arr, item)
	return i < len(arr) && arr[i] == item
}
