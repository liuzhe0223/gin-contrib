package logger

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/sony/sonyflake"
)

const (
	ReqIdHeaderFeild = "X-ReqId"
)

var (
	sf *sonyflake.Sonyflake
)

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func NewReqId() string {
	var buf = make([]byte, binary.MaxVarintLen64)
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	binary.PutUvarint(buf, id)
	return base64.URLEncoding.EncodeToString(buf)
}
