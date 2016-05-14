package nori

type Transport interface {
	Init(*Context) error
}
