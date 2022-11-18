package iotmakerdocker

import (
	"errors"
	"strconv"
	"strings"
)

func (el *DockerSystem) networkGetTheBiggestAddress(ipv4A, ipv4B string) (theBiggest string, err error) {

	//10.0.0.1
	//key.0: 10
	//key.3: 1
	var ipv4AStringList = strings.Split(ipv4A, ".")
	var ipv4BStringList = strings.Split(ipv4B, ".")

	var digitA, digitB int64

	for i := 0; i != 4; i += 1 {
		digitA, err = strconv.ParseInt(ipv4AStringList[i], 10, 64)
		if err != nil {
			return
		}

		digitB, err = strconv.ParseInt(ipv4BStringList[i], 10, 64)
		if err != nil {
			return
		}

		if digitA > digitB {
			return ipv4A, nil
		}

		if digitA < digitB {
			return ipv4B, nil
		}
	}

	err = errors.New("the two addresses are identical")
	return
}
