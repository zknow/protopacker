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

type ProtoPacker struct {
	ProtoHeader       string
	protoHeaderLength int
	buf               []byte
	receivers         []interface{}
}

func NewProtoPacker() *ProtoPacker {
	return &ProtoPacker{
		ProtoHeader:       defaultProtoHeader,
		protoHeaderLength: len(defaultProtoHeader),
	}
}

// 設定自定義Protocol Header
func (bp *ProtoPacker) SetProtoHeader(header string) error {
	var err error
	if len(strings.TrimSpace(header)) == 0 {
		err = ErrProtoheaderNotAvailable
	} else {
		bp.ProtoHeader = strings.TrimSpace(header)
		bp.protoHeaderLength = len(bp.ProtoHeader)
	}
	return err
}

// 註冊封包(bytes)接收者
func (bp *ProtoPacker) RegistBytesReceiver(receivers ...interface{}) error {
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
func (bp *ProtoPacker) Pack(message []byte) []byte {
	return append(append([]byte(bp.ProtoHeader), intToBytes(len(message))...), message...)
}

// 解包
func (bp *ProtoPacker) Unpack(pack []byte) {
	length := len(pack)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+bp.protoHeaderLength+defaultProtoDataLengthByte {
			break
		}
		if string(pack[i:i+bp.protoHeaderLength]) == bp.ProtoHeader {
			msgLen := bytesToInt(pack[i+bp.protoHeaderLength : i+bp.protoHeaderLength+defaultProtoDataLengthByte])
			if length < i+bp.protoHeaderLength+defaultProtoDataLengthByte+msgLen {
				break
			}
			data := pack[i+bp.protoHeaderLength+defaultProtoDataLengthByte : i+bp.protoHeaderLength+defaultProtoDataLengthByte+msgLen]
			bp.riseResult(data)

			i += bp.protoHeaderLength + defaultProtoDataLengthByte + msgLen - 1
		}
	}

	if i == length {
		bp.buf = make([]byte, 0)
	}
	bp.buf = pack[i:]
}

// 觸發＆回傳合法的封包資料
func (bp *ProtoPacker) riseResult(result []byte) {
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
