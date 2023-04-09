package service

import (
	"bytes"
	"encoding/binary"
	"hot-cache/cache"
	"hot-cache/cache/body"
	"hot-cache/cache/request"
	"hot-cache/cache/response"
	"hot-cache/common/io"
	"hot-cache/common/logs"
	"hot-cache/global"
	"net"
)

type Service struct {
	cache cache.Cache
}

func NewService() *Service {
	return &Service{
		cache: cache.NewNormalCache(),
	}
}

func (sv *Service) Handle(conn *net.TCPConn) {
	defer func() {
		_ = conn.Close()
		logs.Info(conn.RemoteAddr(), "disconnected")
	}()
	readBuffer := make([]byte, global.AppConfig.ReadBuffer)
	dataBuffer := bytes.NewBuffer(nil)

	for {
		// loop read
		cnt, err := conn.Read(readBuffer)
		if err != nil {
			io.TryToClose(conn)
			return
		}
		dataBuffer.Write(readBuffer[:cnt])

		// complete size
		if dataBuffer.Len() > body.LenByte {
			var size int64

			// read size
			err = binary.Read(bytes.NewReader(dataBuffer.Bytes()[:body.LenByte]), binary.BigEndian, &size)
			if err != nil {
				logs.Error(err)
				return
			}

			// complete message
			if total := int64(body.LenByte) + size + int64(body.CfgIntegerByte); total <= int64(dataBuffer.Len()) {
				req := request.NewReq()
				var resp *response.Response
				err = req.Setup(dataBuffer.Bytes()[:total])
				if err != nil {
					resp = response.NewResp()
					err = resp.SetupError(req.Cfg[request.CfgIDHigh], req.Cfg[request.CfgIDHigh], response.BadReq.Error())
					if err != nil {
						return
					}
				} else {
					resp = sv.cache.Solve(req)
					if resp == nil {
						resp = response.NewResp()
						err = resp.SetupError(req.Cfg[request.CfgIDHigh], req.Cfg[request.CfgIDHigh], "unknown error")
						if err != nil {
							return
						}
					}
				}
				writeBuf, err := resp.ToBytes()
				if err != nil {
					return
				}
				read := 0
				for read < len(writeBuf) {
					tmp, err := conn.Write(writeBuf[read:])
					if err != nil {
						return
					}
					read += tmp
				}
				tmp := dataBuffer.Bytes()[total:]
				dataBuffer.Reset()
				dataBuffer.Write(tmp)
			}
		}
	}
}
