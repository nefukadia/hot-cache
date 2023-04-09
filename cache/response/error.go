package response

import "errors"

var (
	NotFound   = errors.New("not found")
	MissKey    = errors.New("miss key")
	MissValue  = errors.New("miss value")
	MissExpire = errors.New("miss expire")
	MissInfo   = errors.New("miss info")
	BadReq     = errors.New("bad request")
)
