package networkdelay

func (e *Proxy) SetParserFunction(parser ParserInterface) {
	e.parser = parser
}
