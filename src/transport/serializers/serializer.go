package serializers

type Serializer interface {
	Serialize(any) ([]byte, error)
	ContentType() string
}
