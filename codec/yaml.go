package codec

import (
	"gopkg.in/yaml.v3"
)

type yamlCodec struct{}

var YAML yamlCodec

func (_ yamlCodec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (_ yamlCodec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
