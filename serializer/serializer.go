package serializer

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// ToBytes encode anything to []byte
func ToBytes(v interface{}) []byte {
	msesage, ok := v.(proto.Message)
	if ok {
		bytes, err := proto.Marshal(msesage)
		if err != nil {
			logrus.Error("Proto消息encode失败", err)
			return []byte{}
		}
		return bytes
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		logrus.Error("JSON消息encode失败", err)

		return []byte{123, 125}
	}
	return bytes
}

// FromBytes decode []byte to interface{}
func FromBytes(bytes []byte, v interface{}) {

	message, ok := v.(proto.Message)
	if ok {
		proto.Unmarshal(bytes, message)
		return
	}

	err := json.Unmarshal(bytes, v)
	if err != nil {
		logrus.Error("消息解析失败", err)
		return
	}
}
