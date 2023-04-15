package cache

import (
	"hot-cache/cache/request"
	"hot-cache/cache/response"
	"hot-cache/common/number"
	"hot-cache/common/strings"
	"strconv"
	"sync"
	"time"
)

type Cache interface {
	Solve(req *request.Request) *response.Response
}

const (
	permanent = 31536000000
)

type stringValue struct {
	value    string
	deadline time.Time
}

type NormalCache struct {
	stringMap map[string]stringValue
	stringMu  sync.RWMutex
}

func NewNormalCache() Cache {
	return &NormalCache{
		stringMap: make(map[string]stringValue),
	}
}

func (n *NormalCache) Solve(req *request.Request) (resp *response.Response) {
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

	case request.OptionSet:
		n.stringMu.Lock()
		defer n.stringMu.Unlock()
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
		ret := "false"
		if n.set(key, value, expire) {
			ret = "true"
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, strings.Pointer(ret), nil)
		if err != nil {
			resp = nil
		}
		return

	case request.OptionSetNX:
		n.stringMu.Lock()
		defer n.stringMu.Unlock()
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
		ret := "false"
		if n.setNX(key, value, expire) {
			ret = "true"
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, strings.Pointer(ret), nil)
		if err != nil {
			resp = nil
		}
		return

	case request.OptionDel:
		n.stringMu.Lock()
		defer n.stringMu.Unlock()
		key, ok := req.Data[request.DataKey].(string)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissKey.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		ret := "false"
		if n.del(key) {
			ret = "true"
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, strings.Pointer(ret), nil)
		if err != nil {
			resp = nil
		}
		return

	case request.OptionIncr:
		n.stringMu.Lock()
		defer n.stringMu.Unlock()
		key, ok := req.Data[request.DataKey].(string)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissKey.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		var ret int64
		info := "false"
		if ret, ok = n.incr(key); ok {
			info = "true"
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, number.ToStringPtr(ret), strings.Pointer(info))
		if err != nil {
			resp = nil
		}
		return

	case request.OptionDecr:
		n.stringMu.Lock()
		defer n.stringMu.Unlock()
		key, ok := req.Data[request.DataKey].(string)
		if !ok {
			err := resp.SetupError(cfgIDHigh, cfgIDLow, response.MissKey.Error())
			if err != nil {
				resp = nil
			}
			return
		}
		var ret int64
		info := "false"
		if ret, ok = n.decr(key); ok {
			info = "true"
		}
		err := resp.SetupValue(cfgIDHigh, cfgIDLow, number.ToStringPtr(ret), strings.Pointer(info))
		if err != nil {
			resp = nil
		}
		return
	}
	return nil
}
