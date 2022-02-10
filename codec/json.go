package codec

import "encoding/json"

type jsonCodec struct{}

var JSON jsonCodec

func (_ jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (_ jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
