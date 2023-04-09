package response

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
	OptionError    byte = iota
	OptionDefault  byte = 0xff
	OptionNotFound byte = 0xfe
)

// data
const (
	DataValue byte = iota
	DataInfo
)

// ValueByte other
const (
	ValueByte = 8
)

type Response struct {
	body.BaseBody
	CfgInteger uint64
	DataBytes  []byte
}

func NewResp() *Response {
	return &Response{
		BaseBody: body.BaseBody{
			Len:  0,
			Cfg:  make(map[byte]byte),
			Data: make(map[byte]any),
		},
		CfgInteger: 0,
		DataBytes:  make([]byte, 0),
	}
}

func (resp *Response) Valid() bool {
	// check len
	if resp.Len != int64(len(resp.DataBytes)) {
		return false
	}

	// header checksum
	tmp := uint64(resp.Len) ^ resp.CfgInteger
	if h7, h6, h5, h4, h3, h2, h1, h0 := byte(tmp>>56), byte(tmp>>48), byte(tmp>>40), byte(tmp>>32), byte(tmp>>24), byte(tmp>>16), byte(tmp>>8), byte(tmp); h7^h6^h5^h4^h3^h2^h1^h0 != 0 {
		return false
	}

	// data checksum
	check := byte(resp.CfgInteger >> (8 * CfgBodyChecksum))
	for _, b := range resp.DataBytes {
		check ^= b
	}
	if check != 0 {
		return false
	}

	return true
}

func (resp *Response) ToBytes() ([]byte, error) {
	// check valid
	if !resp.Valid() {
		return nil, errors.New("bad resp")
	}

	buf := bytes.NewBuffer(nil)

	// write to []byte
	err := binary.Write(buf, binary.BigEndian, resp.Len)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, resp.CfgInteger)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(resp.DataBytes)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (resp *Response) Setup(data []byte) (err error) {
	// backup
	tmp := *resp
	defer func() {
		if err != nil {
			resp.BaseBody = tmp.BaseBody
			resp.CfgInteger = tmp.CfgInteger
			resp.DataBytes = tmp.DataBytes
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
	err = binary.Read(lenBuf, binary.BigEndian, &resp.Len)
	if err != nil {
		return err
	}
	err = binary.Read(cfgBuf, binary.BigEndian, &resp.CfgInteger)
	if err != nil {
		return err
	}
	resp.DataBytes = make([]byte, dataBuf.Size())
	_, err = dataBuf.Read(resp.DataBytes)
	if err != nil {
		return err
	}

	// setup to BaseBody
	resp.Cfg, err = resp.parseCfgInteger()
	if err != nil {
		return err
	}
	resp.Data, err = resp.parseDataBytes()
	if err != nil {
		return err
	}

	// check valid
	if !resp.Valid() {
		return errors.New("bad data")
	}

	return nil
}

func (resp *Response) parseCfgInteger() (map[byte]byte, error) {
	ret := make(map[byte]byte)
	ret[CfgHeaderChecksum] = number.LowByteToOther[uint64, byte](resp.CfgInteger >> (8 * CfgHeaderChecksum))
	ret[CfgBodyChecksum] = number.LowByteToOther[uint64, byte](resp.CfgInteger >> (8 * CfgBodyChecksum))
	ret[CfgIDLow] = number.LowByteToOther[uint64, byte](resp.CfgInteger >> (8 * CfgIDLow))
	ret[CfgIDHigh] = number.LowByteToOther[uint64, byte](resp.CfgInteger >> (8 * CfgIDHigh))
	ret[CfgOption] = number.LowByteToOther[uint64, byte](resp.CfgInteger >> (8 * CfgOption))
	if !compare.InList(ret[CfgOption], []byte{
		OptionError,
		OptionDefault,
		OptionNotFound,
	}) {
		return nil, errors.New(fmt.Sprintf("%d not one of option number", ret[CfgOption]))
	}
	return ret, nil
}

func (resp *Response) parseDataBytes() (map[byte]any, error) {
	ret := make(map[byte]any)
	var nextByte int64
	if len(resp.DataBytes) == 0 {
		return ret, nil
	}

	// ans
	if nextByte+ValueByte-1 >= int64(len(resp.DataBytes)) {
		return nil, errors.New("bad red.DataBytes on value")
	}
	buf := bytes.NewReader(resp.DataBytes[nextByte : nextByte+ValueByte])
	var valueLen int64
	err := binary.Read(buf, binary.BigEndian, &valueLen)
	if err != nil {
		return nil, err
	}
	if valueLen > 0 {
		if nextByte+ValueByte+valueLen-1 >= int64(len(resp.DataBytes)) {
			return nil, errors.New("bad red.DataBytes on value")
		}
		ret[DataValue] = string(resp.DataBytes[nextByte+ValueByte : nextByte+ValueByte+valueLen])
	}
	nextByte = nextByte + ValueByte + valueLen

	// info
	if nextByte == int64(len(resp.DataBytes)) {
		return ret, nil
	}
	infoBytes, err := mybytes.VLEUnmarshal(resp.DataBytes[nextByte:])
	if err != nil {
		return nil, err
	}
	if len(infoBytes) > 0 {
		ret[DataInfo] = string(infoBytes)
	}

	return ret, nil
}

func (resp *Response) SetupValue(cfgIDHigh, cfgIDLow byte, value *string, info *string) (err error) {
	tmp := *resp
	defer func() {
		if err != nil {
			resp.BaseBody = tmp.BaseBody
			resp.CfgInteger = tmp.CfgInteger
			resp.DataBytes = tmp.DataBytes
		}
	}()

	// set value
	var valueLen int64
	buf := bytes.NewBuffer(nil)
	if value != nil {
		valueLen = int64(len(*value))
	}
	err = binary.Write(buf, binary.BigEndian, valueLen)
	if err != nil {
		return err
	}
	if value != nil {
		err = binary.Write(buf, binary.BigEndian, []byte(*value))
		if err != nil {
			return err
		}
	}

	// set info
	if info != nil {
		err = binary.Write(buf, binary.BigEndian, mybytes.VLEMarshal([]byte(*info)))
		if err != nil {
			return err
		}
	}

	// update cfg
	var dataXOR byte
	for _, b := range buf.Bytes() {
		dataXOR ^= b
	}
	resp.Len = int64(buf.Len())
	resp.CfgInteger = (uint64(OptionDefault) << (CfgOption * 8)) |
		(uint64(cfgIDHigh) << (CfgIDHigh * 8)) |
		(uint64(cfgIDLow) << (CfgIDLow * 8)) |
		(uint64(dataXOR) << (CfgBodyChecksum * 8))
	var cfgXOR byte
	cfgXOR = number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgBodyChecksum)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgIDLow)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgIDHigh)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgOption)) ^
		number.LowByteToOther[int64, byte](resp.Len>>0) ^
		number.LowByteToOther[int64, byte](resp.Len>>8) ^
		number.LowByteToOther[int64, byte](resp.Len>>16) ^
		number.LowByteToOther[int64, byte](resp.Len>>24) ^
		number.LowByteToOther[int64, byte](resp.Len>>32) ^
		number.LowByteToOther[int64, byte](resp.Len>>40) ^
		number.LowByteToOther[int64, byte](resp.Len>>48) ^
		number.LowByteToOther[int64, byte](resp.Len>>56)
	resp.CfgInteger = resp.CfgInteger ^ (uint64(cfgXOR) << (CfgHeaderChecksum * 8))
	resp.DataBytes = buf.Bytes()

	// setup to BaseBody
	resp.Cfg, err = resp.parseCfgInteger()
	if err != nil {
		return err
	}
	resp.Data, err = resp.parseDataBytes()
	if err != nil {
		return err
	}

	// check valid
	if !resp.Valid() {
		return errors.New("bad data")
	}
	return nil
}

func (resp *Response) SetupError(cfgIDHigh, cfgIDLow byte, info string) (err error) {
	tmp := *resp
	defer func() {
		if err != nil {
			resp.BaseBody = tmp.BaseBody
			resp.CfgInteger = tmp.CfgInteger
			resp.DataBytes = tmp.DataBytes
		}
	}()

	// set value
	var valueLen int64
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.BigEndian, valueLen)
	if err != nil {
		return err
	}

	// set info
	err = binary.Write(buf, binary.BigEndian, mybytes.VLEMarshal([]byte(info)))
	if err != nil {
		return err
	}

	// update cfg
	var dataXOR byte
	for _, b := range buf.Bytes() {
		dataXOR ^= b
	}
	resp.Len = int64(buf.Len())
	resp.CfgInteger = (uint64(OptionError) << (CfgOption * 8)) |
		(uint64(cfgIDHigh) << (CfgIDHigh * 8)) |
		(uint64(cfgIDLow) << (CfgIDLow * 8)) |
		(uint64(dataXOR) << (CfgBodyChecksum * 8))
	var cfgXOR byte
	cfgXOR = number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgBodyChecksum)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgIDLow)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgIDHigh)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgOption)) ^
		number.LowByteToOther[int64, byte](resp.Len>>0) ^
		number.LowByteToOther[int64, byte](resp.Len>>8) ^
		number.LowByteToOther[int64, byte](resp.Len>>16) ^
		number.LowByteToOther[int64, byte](resp.Len>>24) ^
		number.LowByteToOther[int64, byte](resp.Len>>32) ^
		number.LowByteToOther[int64, byte](resp.Len>>40) ^
		number.LowByteToOther[int64, byte](resp.Len>>48) ^
		number.LowByteToOther[int64, byte](resp.Len>>56)
	resp.CfgInteger = resp.CfgInteger ^ (uint64(cfgXOR) << (CfgHeaderChecksum * 8))
	resp.DataBytes = buf.Bytes()

	// setup to BaseBody
	resp.Cfg, err = resp.parseCfgInteger()
	if err != nil {
		return err
	}
	resp.Data, err = resp.parseDataBytes()
	if err != nil {
		return err
	}

	// check valid
	if !resp.Valid() {
		return errors.New("bad data")
	}
	return nil
}

func (resp *Response) SetupNotFound(cfgIDHigh, cfgIDLow byte) (err error) {
	tmp := *resp
	defer func() {
		if err != nil {
			resp.BaseBody = tmp.BaseBody
			resp.CfgInteger = tmp.CfgInteger
			resp.DataBytes = tmp.DataBytes
		}
	}()

	// set value
	var valueLen int64
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.BigEndian, valueLen)
	if err != nil {
		return err
	}

	// update cfg
	var dataXOR byte
	for _, b := range buf.Bytes() {
		dataXOR ^= b
	}
	resp.Len = int64(buf.Len())
	resp.CfgInteger = (uint64(OptionNotFound) << (CfgOption * 8)) |
		(uint64(cfgIDHigh) << (CfgIDHigh * 8)) |
		(uint64(cfgIDLow) << (CfgIDLow * 8)) |
		(uint64(dataXOR) << (CfgBodyChecksum * 8))
	var cfgXOR byte
	cfgXOR = number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgBodyChecksum)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgIDLow)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgIDHigh)) ^
		number.LowByteToOther[uint64, byte](resp.CfgInteger>>(8*CfgOption)) ^
		number.LowByteToOther[int64, byte](resp.Len>>0) ^
		number.LowByteToOther[int64, byte](resp.Len>>8) ^
		number.LowByteToOther[int64, byte](resp.Len>>16) ^
		number.LowByteToOther[int64, byte](resp.Len>>24) ^
		number.LowByteToOther[int64, byte](resp.Len>>32) ^
		number.LowByteToOther[int64, byte](resp.Len>>40) ^
		number.LowByteToOther[int64, byte](resp.Len>>48) ^
		number.LowByteToOther[int64, byte](resp.Len>>56)
	resp.CfgInteger = resp.CfgInteger ^ (uint64(cfgXOR) << (CfgHeaderChecksum * 8))
	resp.DataBytes = buf.Bytes()

	// setup to BaseBody
	resp.Cfg, err = resp.parseCfgInteger()
	if err != nil {
		return err
	}
	resp.Data, err = resp.parseDataBytes()
	if err != nil {
		return err
	}

	// check valid
	if !resp.Valid() {
		return errors.New("bad data")
	}
	return nil
}
