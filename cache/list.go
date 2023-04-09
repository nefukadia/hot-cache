package cache

import "math/rand"

func (n *NormalCache) front(key string) (string, bool) {
	list, ok := n.listMap[key]
	if !ok || len(list) == 0 {
		return "", false
	}
	return list[0], true
}

func (n *NormalCache) inList(key string, value string) bool {
	list, ok := n.listMap[key]
	if !ok {
		return false
	}
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func (n *NormalCache) randomInList(key string) (string, bool) {
	list, ok := n.listMap[key]
	if !ok || len(list) == 0 {
		return "", false
	}
	randIndex := rand.Intn(len(list))
	return list[randIndex], true
}

func (n *NormalCache) pushback(key string, value string) {
	list, _ := n.listMap[key]
	list = append(list, value)
	n.listMap[key] = list
}

func (n *NormalCache) delInList(key string, value string) bool {
	if n.inList(key, value) {
		list := n.listMap[key]
		n.listMap[key] = make([]string, 0)
		for _, v := range list {
			if v != value {
				n.listMap[key] = append(n.listMap[key], v)
			}
		}
		return true
	}
	return false
}
