package service

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hot-cache/cache"
	"hot-cache/common/logs"
	"hot-cache/global"
	"net"
	"unsafe"
)

func init() {
	if cache.SizeByte > int(unsafe.Sizeof(0)) {
		panic(fmt.Sprintf("cache.SizeByte cannot > %v", unsafe.Sizeof(0)))
	}
}

func Handle(conn *net.TCPConn) {
	defer func() {
		_ = conn.Close()
	}()
	readBuffer := make([]byte, global.AppConfig.ReadBuffer)
	dataBuffer := bytes.NewBuffer(nil)

	for {
		// loop read
		cnt, err := conn.Read(readBuffer)
		if err != nil {
			return
		}
		dataBuffer.Write(readBuffer[:cnt])

		// complete size
		if dataBuffer.Len() > cache.SizeByte {
			var size int64

			// read size
			err = binary.Read(bytes.NewReader(dataBuffer.Bytes()[:cache.SizeByte]), binary.BigEndian, &size)
			if err != nil {
				logs.Warning(err)
				return
			}

			// complete message
			if total := int64(cache.SizeByte) + size + int64(cache.OptionByte); total <= int64(dataBuffer.Len()) {
				var op cache.Option
				err = binary.Read(bytes.NewReader(dataBuffer.Bytes()[cache.SizeByte:cache.SizeByte+cache.OptionByte]), binary.BigEndian, (*byte)(&op))
				if err != nil {
					logs.Warning(err)
					return
				}
				rep := global.AppCache.Solve(cache.Res{
					Op:   op,
					Data: dataBuffer.Bytes()[cache.SizeByte+cache.OptionByte : int64(cache.SizeByte+cache.OptionByte)+size],
				})

				// write rep byte
				size = int64(len(rep.Data))
				tmp := bytes.NewBuffer(nil)
				err = binary.Write(tmp, binary.BigEndian, size)
				if err != nil {
					logs.Warning(err)
					return
				}
				err = binary.Write(tmp, binary.BigEndian, byte(rep.Code))
				if err != nil {
					logs.Warning(err)
					return
				}
				tmp.Write(rep.Data)

				// send rep
				var writeIndex = 0
				for writeIndex < rep.Len() {
					cnt, err = conn.Write(tmp.Bytes()[writeIndex:])
					if err != nil {
						logs.Warning(err)
						return
					}
					writeIndex += cnt
				}
			}
		}
	}
}
