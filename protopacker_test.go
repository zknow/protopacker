package protopacker

import (
	"bytes"
	"testing"
)

func Test1(t *testing.T) {
	mockData := []byte("Mock Data")
	bp := NewBytesPacker()
	bp.RegistBytesReceiver(func(result []byte) {
		if !bytes.Equal(result, mockData) {
			t.Error("Data Not Matched!")
		} else {
			t.Log("Data Matched!")
		}
	})
	bsPacked := bp.Pack(mockData)
	bp.Unpack(bsPacked)
}
