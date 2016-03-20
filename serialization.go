package nori

type Serializer interface {
	Name() string
	ContentType() string
}

type JSONSerializer struct {
}

var _ Serializer = (*JSONSerializer)(nil)

func (JSONSerializer) Name() string { return "JSONSerializer" }

func (JSONSerializer) ContentType() string { return "application/json" }
