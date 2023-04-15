package cache

import (
	"strconv"
	"time"
)

func (n *NormalCache) get(key string) (string, int64, bool) {
	ret, ok := n.stringMap[key]
	now := time.Now()
	if !ok || ret.deadline.Before(now) {
		delete(n.stringMap, key)
		return "", 0, false
	}
	return ret.value, int64(ret.deadline.Sub(now) / time.Second), true
}

func (n *NormalCache) set(key, value string, expire int64) bool {
	newValue := stringValue{
		value:    value,
		deadline: time.Now().Add(time.Duration(expire) * time.Second),
	}
	cover := false
	if _, _, ok := n.get(key); ok {
		cover = true
	}
	n.stringMap[key] = newValue
	return cover
}

func (n *NormalCache) setNX(key, value string, expire int64) bool {
	if _, _, ok := n.get(key); ok {
		return false
	}
	n.set(key, value, expire)
	return true
}

func (n *NormalCache) del(key string) bool {
	if _, _, ok := n.get(key); ok {
		delete(n.stringMap, key)
		return true
	}
	return false
}

func (n *NormalCache) incr(key string) (int64, bool) {
	value, remain, ok := n.get(key)
	if !ok {
		value = "1"
		remain = permanent
	}
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	num++
	n.set(key, strconv.FormatInt(num, 10), remain)
	return num, true
}

func (n *NormalCache) decr(key string) (int64, bool) {
	value, remain, ok := n.get(key)
	if !ok {
		value = "-1"
		remain = permanent
	}
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	num--
	n.set(key, strconv.FormatInt(num, 10), remain)
	return num, true
}
