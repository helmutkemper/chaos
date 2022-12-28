package networkdelay

type Proxy struct {
	min        int
	max        int
	bufferSize int

	parser ParserInterface
}
