package p2p

const (
	DORKS = iota
	URLS
	INJECTABLES
	PROXIES
)

type List struct {
	Data []string
	Type int
}