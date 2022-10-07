package utils

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
)

var json jsoniter.API

var m jsonpb.Marshaler

func init() {
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	m = jsonpb.Marshaler{}
}

func Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func Unmarshal(data []byte, target interface{}) error {
	return json.Unmarshal(data, &target)
}

func MarshalPb(pb proto.Message) string {
	result, _ := m.MarshalToString(pb)
	return result
}
