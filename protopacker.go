package protopacker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	defaultProtoHeader         = "@@Header@@"
	defaultProtoHeaderLength   = len(defaultProtoHeader)
	defaultProtoDataLengthByte = 4
)

var (
	ErrProtoheaderNotAvailable = fmt.Errorf("Not Available Header")
	ErrUndefinedReceiver       = fmt.Errorf("Not Undefined Receiver,Use 'chan []byte' or 'func([]byte)'")
)

type BytesPacker struct {
	ProtoHeader           string
	protoHeaderByteLength int
	buf                   []byte
	receivers             []interface{}
}

func NewBytesPacker() *BytesPacker {
	return &BytesPacker{
		ProtoHeader:           defaultProtoHeader,
		protoHeaderByteLength: len(defaultProtoHeader),
	}
}

// 設定自定義Protocol Header
func (bp *BytesPacker) SetProtoHeader(header string) error {
	var err error
	if len(strings.TrimSpace(header)) == 0 {
		err = ErrProtoheaderNotAvailable
	} else {
		bp.ProtoHeader = strings.TrimSpace(header)
		bp.protoHeaderByteLength = len(bp.ProtoHeader)
	}
	return err
}

// 註冊封包(bytes)接收者
func (bp *BytesPacker) RegistBytesReceiver(receivers ...interface{}) error {
	for _, rcv := range receivers {
		switch rcv.(type) {
		case chan []byte:
		case func([]byte):
		default:
			return ErrUndefinedReceiver
		}
	}
	bp.receivers = receivers
	return nil
}

// 封包
func (bp *BytesPacker) Pack(message []byte) []byte {
	return append(append([]byte(bp.ProtoHeader), intToBytes(len(message))...), message...)
}

// 解包
func (bp *BytesPacker) Unpack(pack []byte) {
	length := len(pack)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+defaultProtoHeaderLength+defaultProtoDataLengthByte {
			break
		}
		if string(pack[i:i+defaultProtoHeaderLength]) == bp.ProtoHeader {
			msgLen := bytesToInt(pack[i+defaultProtoHeaderLength : i+defaultProtoHeaderLength+defaultProtoDataLengthByte])
			if length < i+defaultProtoHeaderLength+defaultProtoDataLengthByte+msgLen {
				break
			}
			data := pack[i+defaultProtoHeaderLength+defaultProtoDataLengthByte : i+defaultProtoHeaderLength+defaultProtoDataLengthByte+msgLen]
			bp.RiseResult(data)

			i += defaultProtoHeaderLength + defaultProtoDataLengthByte + msgLen - 1
		}
	}

	if i == length {
		bp.buf = make([]byte, 0)
	}
	bp.buf = pack[i:]
}

// 觸發＆回傳合法的封包資料
func (bp *BytesPacker) RiseResult(result []byte) {
	for _, rcv := range bp.receivers {
		switch rcv.(type) {
		case chan []byte:
			rcv.(chan []byte) <- result
		case func([]byte):
			rcv.(func([]byte))(result)
		default:
		}
	}
}

func intToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func bytesToInt(b []byte) int {
	var x int32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &x)
	return int(x)
}
