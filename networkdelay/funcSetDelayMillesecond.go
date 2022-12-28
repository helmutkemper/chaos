package networkdelay

func (e *Proxy) SetDelayMillisecond(min, max int) {
	e.min = min
	e.max = max
}
