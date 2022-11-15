package docker

import "time"

func (e *ContainerBuilder) selectBetweenMaxAndMin(max, min time.Duration) (selected time.Duration) {
	if int64(max)-int64(min) == 0 {
		return min
	}

	randValue := e.getRandSeed().Int63n(int64(max)-int64(min)) + int64(min)
	return time.Duration(randValue)
}
