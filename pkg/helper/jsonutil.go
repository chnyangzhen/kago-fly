package helper

import (
	"github.com/json-iterator/go"
)

var j *JSONSerializer

func init() {
	j = NewJsonSerializer()
}

func NewJsonSerializer() *JSONSerializer {
	return &JSONSerializer{json: jsoniter.ConfigCompatibleWithStandardLibrary}
}

type JSONSerializer struct {
	json jsoniter.API
}

// JSONToBytes 将对象序列化为字节数组
func JSONToBytes(v interface{}) ([]byte, error) {
	return j.json.Marshal(v)
}

// JSONToObject 将字节数组反序列化为对象；
func JSONToObject(data []byte, out interface{}) error {
	return j.json.Unmarshal(data, out)
}
