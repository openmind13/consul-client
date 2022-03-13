package servicediscovery

import (
	"bytes"
	"encoding/json"
)

type Config struct {
	Id   int
	Data string
}

func (c *Config) marshalJSON() []byte {
	buffer := bytes.Buffer{}
	json.NewEncoder(&buffer).Encode(c)
	return buffer.Bytes()
}

// func (c *Config) MarshalBytes() []byte {
// 	buffer := bytes.Buffer{}
// 	binary.Write(&buffer, binary.BigEndian, c)
// 	return buffer.Bytes()
// }
