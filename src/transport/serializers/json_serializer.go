package serializers

import "encoding/json"

type JsonSerializer struct {
}

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

func (j *JsonSerializer) Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

func (j *JsonSerializer) ContentType() string {
	return "application/json"
}
