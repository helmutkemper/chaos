package docker

import (
	"errors"
	"fmt"
	"github.com/helmutkemper/util"
	"strconv"
	"strings"
)

// incIpV4Address
//
// English:
//
//	Receives an IP address in the form of a string and increments it.
//
//	 Input:
//	   ip:  only the ip address. e.g.: 10.0.0.1
//	   inc: number of increments
//
//	 Output:
//	   next: the next ip address. e.g.: 10.0.0.2
//	   err: standard error object
//
// Note:
//
//   - This function does not take into account the network configuration, use it with care!
//
// Português:
//
//	Recebe um endereço IP na forma de string e incrementa o mesmo.
//
//	 Entrada:
//	   ip:  apenas o endereço ip. ex.: 10.0.0.1
//	   inc: quantidade de incrementos
//
//	 Saída:
//	   next: o próximo endereço ip. ex.: 10.0.0.2
//	   err: objeto de erro padrão
//
// Nota:
//
//	Esta função não considera a configuração da rede, use com cuidado!
func (e *ContainerBuilder) incIpV4Address(ip string, inc int64) (next string, err error) {

	// está na rede padrão do docker
	if ip == "0.0.0.0" {
		next = "0.0.0.0"
		return
	}

	var digitList []string
	digitList = strings.Split(ip, "/")
	digitList = strings.Split(digitList[0], ".")

	var digitA, digitB, digitC, digitD, overflow int64
	digitA, err = strconv.ParseInt(digitList[0], 10, 64)
	if err != nil {
		util.TraceToLog()
		return
	}

	digitB, err = strconv.ParseInt(digitList[1], 10, 64)
	if err != nil {
		util.TraceToLog()
		return
	}

	digitC, err = strconv.ParseInt(digitList[2], 10, 64)
	if err != nil {
		util.TraceToLog()
		return
	}

	digitD, err = strconv.ParseInt(digitList[3], 10, 64)
	if err != nil {
		util.TraceToLog()
		return
	}

	digitD += inc
	if digitD > 255 {
		digitD = 0
		overflow = 1
	} else {
		overflow = 0
	}

	digitC += overflow
	if digitC > 255 {
		digitC = 0
		overflow = 1
	} else {
		overflow = 0
	}

	digitB += overflow
	if digitB > 255 {
		digitB = 0
		overflow = 1
	} else {
		overflow = 0
	}

	digitA += overflow
	if digitA > 255 {
		digitA = 0
		overflow = 1
	} else {
		overflow = 0
	}

	if overflow != 0 {
		util.TraceToLog()
		err = errors.New("ip overflow")
		return
	}

	next = fmt.Sprintf("%v.%v.%v.%v", digitA, digitB, digitC, digitD)
	return
}
