package cache

import (
	"hot-cache/common/set"
	"math/rand"
)

func (n *NormalCache) inSet(key string, value string) bool {
	theSet, ok := n.setMap[key]
	if !ok {
		return false
	}
	return theSet.Exist(value)
}

func (n *NormalCache) randomInSet(key string) (string, bool) {
	theSet, ok := n.setMap[key]
	if !ok {
		return "", false
	}
	length := theSet.Len()
	if length == 0 {
		return "", false
	}
	randIndex := rand.Intn(length)
	var ans string
	theSet.Range(func(i int, s string) bool {
		if i < randIndex {
			return true
		}
		ans = s
		return false
	})
	return ans, true
}

func (n *NormalCache) insert(key string, value string) bool {
	theSet, ok := n.setMap[key]
	if !ok {
		theSet = set.NewHashSet[string]()
	}
	return theSet.Insert(value)
}

func (n *NormalCache) delInSet(key string, value string) bool {
	theSet, ok := n.setMap[key]
	if !ok {
		theSet = set.NewHashSet[string]()
	}
	return theSet.Del(value)
}
