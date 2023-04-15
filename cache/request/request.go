package request

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hot-cache/cache/body"
	mybytes "hot-cache/common/bytes"
	"hot-cache/common/compare"
	"hot-cache/common/number"
)

// cfg
const (
	CfgHeaderChecksum byte = iota
	CfgBodyChecksum
	CfgIDLow
	CfgIDHigh
	_
	_
	_
	CfgOption
)

// option
const (
	OptionGet byte = iota + 1
	OptionSet
	OptionSetNX
	OptionDel
	OptionIncr
	OptionDecr

	OptionAuth
)

// data
const (
	DataKey byte = iota
	DataValue
	DataExpire
	DataInfo
)

// other
const (
	ValueByte  = 8
	ExpireByte = 8
)

type Request struct {
	body.BaseBody
	CfgInteger uint64
	DataBytes  []byte
}

func NewReq() *Request {
	return &Request{
		BaseBody: body.BaseBody{
			Len:  0,
			Cfg:  make(map[byte]byte),
			Data: make(map[byte]any),
		},
		CfgInteger: 0,
		DataBytes:  make([]byte, 0),
	}
}

func (req *Request) Valid() bool {
	// check len
	if req.Len != int64(len(req.DataBytes)) {
		return false
	}

	// header checksum
	tmp := uint64(req.Len) ^ req.CfgInteger
	if h7, h6, h5, h4, h3, h2, h1, h0 := byte(tmp>>56), byte(tmp>>48), byte(tmp>>40), byte(tmp>>32), byte(tmp>>24), byte(tmp>>16), byte(tmp>>8), byte(tmp); h7^h6^h5^h4^h3^h2^h1^h0 != 0 {
		return false
	}

	// data checksum
	check := byte(req.CfgInteger >> (8 * CfgBodyChecksum))
	for _, b := range req.DataBytes {
		check ^= b
	}
	if check != 0 {
		return false
	}

	return true
}

func (req *Request) ToBytes() ([]byte, error) {
	// check valid
	if !req.Valid() {
		return nil, errors.New("bad req")
	}

	buf := bytes.NewBuffer(nil)

	// write to []byte
	err := binary.Write(buf, binary.BigEndian, req.Len)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, req.CfgInteger)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(req.DataBytes)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (req *Request) Setup(data []byte) (err error) {
	// backup
	tmp := *req
	defer func() {
		if err != nil {
			req.BaseBody = tmp.BaseBody
			req.CfgInteger = tmp.CfgInteger
			req.DataBytes = tmp.DataBytes
		}
	}()

	// check valid
	if len(data) < body.LenByte+body.CfgIntegerByte {
		return errors.New("bad data: len(data) < LenByte + CfgIntegerByte")
	}

	lenBuf := bytes.NewReader(data[:body.LenByte])
	cfgBuf := bytes.NewReader(data[body.LenByte : body.LenByte+body.CfgIntegerByte])
	dataBuf := bytes.NewReader(data[body.LenByte+body.CfgIntegerByte:])

	// setup from []byte
	err = binary.Read(lenBuf, binary.BigEndian, &req.Len)
	if err != nil {
		return err
	}
	err = binary.Read(cfgBuf, binary.BigEndian, &req.CfgInteger)
	if err != nil {
		return err
	}
	req.DataBytes = make([]byte, dataBuf.Size())
	_, err = dataBuf.Read(req.DataBytes)
	if err != nil {
		return err
	}

	// setup to BaseBody
	req.Cfg, err = req.parseCfgInteger()
	if err != nil {
		return err
	}
	req.Data, err = req.parseDataBytes()
	if err != nil {
		return err
	}

	// check valid
	if !req.Valid() {
		return errors.New("bad data")
	}

	return nil
}

func (req *Request) parseCfgInteger() (map[byte]byte, error) {
	ret := make(map[byte]byte)
	ret[CfgHeaderChecksum] = number.LowByteToOther[uint64, byte](req.CfgInteger >> (8 * CfgHeaderChecksum))
	ret[CfgBodyChecksum] = number.LowByteToOther[uint64, byte](req.CfgInteger >> (8 * CfgBodyChecksum))
	ret[CfgIDLow] = number.LowByteToOther[uint64, byte](req.CfgInteger >> (8 * CfgIDLow))
	ret[CfgIDHigh] = number.LowByteToOther[uint64, byte](req.CfgInteger >> (8 * CfgIDHigh))
	ret[CfgOption] = number.LowByteToOther[uint64, byte](req.CfgInteger >> (8 * CfgOption))
	if !compare.InList(ret[CfgOption], []byte{
		OptionGet,
		OptionSet,
		OptionSetNX,
		OptionDel,
		OptionIncr,
		OptionDecr,
		OptionAuth,
	}) {
		return nil, errors.New(fmt.Sprintf("%d not one of option number", ret[CfgOption]))
	}
	return ret, nil
}

func (req *Request) parseDataBytes() (map[byte]any, error) {
	ret := make(map[byte]any)
	var nextByte int64
	if len(req.DataBytes) == 0 {
		return ret, nil
	}

	// key
	var tmp []byte
	if req.Cfg[CfgOption] != OptionAuth {
		for _, b := range req.DataBytes {
			nextByte++
			if b&0x80 != 0 {
				if nextByte == int64(len(req.DataBytes)) {
					return nil, errors.New("bad red.DataBytes on key")
				}
				tmp = append(tmp, b)
			} else {
				tmp = append(tmp, b)
				keyBytes, err := mybytes.VLEUnmarshal(tmp)
				if err != nil {
					return nil, err
				}
				ret[DataKey] = string(keyBytes)
				break
			}
		}
	} else {
		nextByte = 1
	}

	// value
	if nextByte == int64(len(req.DataBytes)) {
		return ret, nil
	}
	if nextByte+ValueByte-1 >= int64(len(req.DataBytes)) {
		return nil, errors.New("bad red.DataBytes on value")
	}
	buf := bytes.NewReader(req.DataBytes[nextByte : nextByte+ValueByte])
	var valueLen int64
	err := binary.Read(buf, binary.BigEndian, &valueLen)
	if err != nil {
		return nil, err
	}
	if nextByte+ValueByte+valueLen-1 >= int64(len(req.DataBytes)) {
		return nil, errors.New("bad red.DataBytes on value")
	}
	ret[DataValue] = string(req.DataBytes[nextByte+ValueByte : nextByte+ValueByte+valueLen])
	nextByte = nextByte + ValueByte + valueLen

	// expire
	if nextByte == int64(len(req.DataBytes)) {
		return ret, nil
	}
	if nextByte+ExpireByte-1 >= int64(len(req.DataBytes)) {
		return nil, errors.New("bad red.DataBytes on value")
	}
	buf = bytes.NewReader(req.DataBytes[nextByte : nextByte+ExpireByte])
	var expire int64
	err = binary.Read(buf, binary.BigEndian, &expire)
	if err != nil {
		return nil, err
	}
	ret[DataExpire] = expire
	nextByte = nextByte + ExpireByte

	// info
	if nextByte == int64(len(req.DataBytes)) {
		return ret, nil
	}
	infoBytes, err := mybytes.VLEUnmarshal(req.DataBytes[nextByte:])
	if err != nil {
		return nil, err
	}
	if len(infoBytes) > 0 {
		ret[DataInfo] = string(infoBytes)
	}

	return ret, nil
}
