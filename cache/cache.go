package cache

import (
	"hot-cache/cache/request"
	"hot-cache/cache/response"
	"hot-cache/common/set"
	"hot-cache/common/strings"
	"strconv"
	"sync"
	"time"
)

type Cache interface {
	Solve(req *request.Request) *response.Response
}

const (
	permanent = -1
)

type stringValue struct {
	value    string
	deadline time.Time
}

type NormalCache struct {
	stringMap map[string]stringValue
	listMap   map[string][]string
	setMap    map[string]set.Set[string]

	stringMu sync.RWMutex
	listMu   sync.RWMutex
	setMu    sync.RWMutex
}

func NewNormalCache() Cache {
	return &NormalCache{
		stringMap: make(map[string]stringValue),
		listMap:   make(map[string][]string),
		setMap:    make(map[string]set.Set[string]),
	}
}

func (n *NormalCache) Solve(req *request.Request) (resp *response.Response) {
	// todo
	resp = response.NewResp()
	cfgIDHigh := req.Cfg[request.CfgIDHigh]
	cfgIDLow := req.Cfg[request.CfgIDLow]
	switch req.Cfg[request.CfgOption] {
	case request.OptionGet:
		n.stringMu.RLock()
		defer n.stringMu.RUnlock()
		key, ok := req.Data[request.DataKey].(string)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissKey.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		value, expire, ok := n.get(key)
		if !ok {
			err := resp.SetupNotFound(cfgIDHigh, cfgIDLow)
			if err != nil {
				resp = nil
			}
			return
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, &value, strings.Pointer(strconv.FormatInt(expire, 10)))
		if err != nil {
			resp = nil
		}
		return
	case request.OptionSet, request.OptionSetNX, request.OptionDel,
		request.OptionIncr, request.OptionDecr:
		n.stringMu.Lock()
		defer n.stringMu.Unlock()
		// todo
		key, ok := req.Data[request.DataKey].(string)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissKey.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		value, ok := req.Data[request.DataValue].(string)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissValue.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		expire, ok := req.Data[request.DataExpire].(int64)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissExpire.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		info := "false"
		if n.set(key, value, expire) {
			info = "true"
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, nil, strings.Pointer(info))
		if err != nil {
			resp = nil
		}
		return
	case request.OptionFront, request.OptionIsInList, request.OptionRandomInList:
		n.listMu.RLock()
		defer n.listMu.RUnlock()

	case request.OptionPushback, request.OptionDelInList:
		n.listMu.Lock()
		defer n.listMu.Unlock()

	case request.OptionIsInSet, request.OptionRandomInSet:
		n.setMu.RLock()
		defer n.setMu.RUnlock()

	case request.OptionInsert, request.OptionDelInSet:
		n.setMu.Lock()
		defer n.setMu.Unlock()

	case request.OptionAuth:

	}

	return nil
}
