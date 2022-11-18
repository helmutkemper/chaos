package iotmakerdocker

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type IPv4Generator struct {
	ipA        byte
	ipB        byte
	ipC        byte
	ipD        byte
	gatewayA   byte
	gatewayB   byte
	gatewayC   byte
	gatewayD   byte
	subnetA    byte
	subnetB    byte
	subnetC    byte
	subnetD    byte
	subnetCidr byte
	maxA       byte
	maxB       byte
	maxC       byte
	maxD       byte
	ipMinAddr  int64
	ipMaxAddr  int64
	reserved   [][4]byte
}

func (el *IPv4Generator) GetCurrentIP() (digitA, digitB, digitC, digitD byte) {
	return el.ipA, el.ipB, el.ipC, el.ipD
}

func (el *IPv4Generator) GetCurrentIPAsString() (ip string) {
	return el.String()
}

func (el *IPv4Generator) int64ToIP(ipInt64 int64) (a, b, c, d byte) {

	// 10.0.0.1 = ipA.ipB.ipC.ipD
	// = ipA*256^3 + ipB*256^2 + ipC*256^1 + ipD*256^0
	// = 10*256^3 + 0*256^2 + 0*256^1 + 1*256^0
	// = 167772161
	// = 10.0.0.1 = integer 167772161
	//
	// interger = 167772161
	// ipD = integer % 256
	// integer = integer / 256
	// ipC = integer % 256
	// integer = integer / 256
	// ipB = integer % 256
	// integer = integer / 256
	// ipA = integer % 256
	d = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	c = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	b = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	a = byte(ipInt64 % 256)

	return
}

func (el *IPv4Generator) stringIPtoByte(ip string) (a, b, c, d, cidr byte, err error) {
	var tmpInt64 int64
	var ipElementsList = strings.Split(ip, "/")
	var ipOnly = ipElementsList[0]

	if len(ipElementsList) > 1 {
		tmpInt64, err = strconv.ParseInt(ipElementsList[1], 10, 64)
		if err != nil {
			err = errors.New("IP format must be 'int.int.int.int' or 'int.int.int.int/int'")
			return
		}
		cidr = byte(tmpInt64)
	} else {
		cidr = 0
	}

	ipElementsList = strings.Split(ipOnly, ".")
	if len(ipElementsList) != 4 {
		err = errors.New("IP format must be 'int.int.int.int' or 'int.int.int.int/int'")
		return
	}

	tmpInt64, err = strconv.ParseInt(ipElementsList[0], 10, 64)
	if err != nil {
		err = errors.New("IP format must be 'int.int.int.int' or 'int.int.int.int/int'")
		return
	}
	a = byte(tmpInt64)

	tmpInt64, err = strconv.ParseInt(ipElementsList[1], 10, 64)
	if err != nil {
		err = errors.New("IP format must be 'int.int.int.int' or 'int.int.int.int/int'")
		return
	}
	b = byte(tmpInt64)

	tmpInt64, err = strconv.ParseInt(ipElementsList[2], 10, 64)
	if err != nil {
		err = errors.New("IP format must be 'int.int.int.int' or 'int.int.int.int/int'")
		return
	}
	c = byte(tmpInt64)

	tmpInt64, err = strconv.ParseInt(ipElementsList[3], 10, 64)
	if err != nil {
		err = errors.New("IP format must be 'int.int.int.int' or 'int.int.int.int/int'")
		return
	}
	d = byte(tmpInt64)

	return
}

func (el IPv4Generator) split(
	ip string,
) (
	a,
	b,
	c,
	d,
	CIDRPrefix int,
	err error,
) {

	var tmpA, tmpB, tmpC, tmpD, tmpE string
	var aInt64, bInt64, cInt64, dInt64, CIDRPrefixInt64 int64

	list := strings.Split(ip, ".")
	tmpA = list[0]
	tmpB = list[1]
	tmpC = list[2]
	tmpD = list[3]

	tmpEArr := strings.Split(tmpD, "/")
	tmpD = tmpEArr[0]

	if len(tmpEArr) == 2 {
		tmpE = tmpEArr[1]
	} else {
		tmpE = ""
	}

	aInt64, err = strconv.ParseInt(tmpA, 10, 32)
	if err != nil {
		return
	}

	bInt64, err = strconv.ParseInt(tmpB, 10, 32)
	if err != nil {
		return
	}

	cInt64, err = strconv.ParseInt(tmpC, 10, 32)
	if err != nil {
		return
	}

	dInt64, err = strconv.ParseInt(tmpD, 10, 32)
	if err != nil {
		return
	}

	if tmpE == "" {
		CIDRPrefixInt64 = 0
	} else {
		CIDRPrefixInt64, err = strconv.ParseInt(tmpE, 10, 32)
	}

	a = int(aInt64)
	b = int(bInt64)
	c = int(cInt64)
	d = int(dInt64)
	CIDRPrefix = int(CIDRPrefixInt64)

	return
}

func (el IPv4Generator) cidrPrefixToDecimal(
	cidr byte,
) (
	cidrDecimal byte,
) {

	return byte(int(math.Pow(2.0, float64(int(cidr))) - 1))
}

func (el IPv4Generator) verify(
	ipA,
	ipB,
	ipC,
	ipD byte,
) (
	err error,
) {

	var ipAsInt int64

	if el.subnetA == 0 || el.gatewayA == 0 {
		err = errors.New("initialize IPv4Generator{} first")
		return
	}

	var caseA, caseB, caseC, caseD bool
	for _, ipReserved := range el.reserved {
		caseA = ipReserved[0] == ipA
		caseB = ipReserved[1] == ipB
		caseC = ipReserved[2] == ipC
		caseD = ipReserved[3] == ipD

		if caseA && caseB && caseC && caseD == true {
			err = errors.New("ip reserved")
			return
		}
	}

	// 10.0.0.1 = ipA.ipB.ipC.ipD
	// = ipA*256^3 + ipB*256^2 + ipC*256^1 + ipD*256^0
	// = 10*256^3 + 0*256^2 + 0*256^1 + 1*256^0
	// = 167772161
	// = 10.0.0.1 As integer
	tmpA := int64(float64(ipA) * math.Pow(256.0, 3.0))
	tmpB := int64(float64(ipB) * math.Pow(256.0, 2.0))
	tmpC := int64(float64(ipC) * math.Pow(256.0, 1.0))
	tmpD := int64(float64(ipD) * math.Pow(256.0, 0.0))

	ipAsInt = tmpA + tmpB + tmpC + tmpD

	if ipAsInt < el.ipMinAddr {
		err = errors.New(fmt.Sprintf("min allowed ip is %v", el.ipAsStringForError(el.gatewayA, el.gatewayB, el.gatewayC, el.gatewayD, 1)))
		return
	}

	//fixme: bug subnet está invertida
	if ipAsInt > el.ipMaxAddr {
		//err = errors.New(fmt.Sprintf("max allowed ip is %v", el.ipAsStringForError(el.maxA, el.maxB, el.maxC, el.maxD, 0)))
		//return
	}

	return
}

// ipAsStringForError (português): esta função só deve ser usada para transformar IP em mensagem de erro
//
//	a, b, c, d: são os elementos do ip
//	inc: incrementa o ip (ip = ip + inc)
func (el *IPv4Generator) ipAsStringForError(a, b, c, d, inc byte) (
	ipAsString string,
) {

	var nextDecimalPlace byte

	d += inc
	if d > 255 {
		nextDecimalPlace = 1
	} else {
		nextDecimalPlace = 0
	}

	c += nextDecimalPlace
	if c > 255 {
		nextDecimalPlace = 1
	} else {
		nextDecimalPlace = 0
	}

	b += nextDecimalPlace
	if b > 255 {
		nextDecimalPlace = 1
	} else {
		nextDecimalPlace = 0
	}

	a += nextDecimalPlace

	ipAsString = strconv.Itoa(int(a)) + "." + strconv.Itoa(int(b)) + "." + strconv.Itoa(int(c)) + "." + strconv.Itoa(int(d))
	return
}

func (el *IPv4Generator) incIP(a, b, c, d, inc byte) (
	elementA, elementB, elementC, elementD, overflow byte,
) {

	var nextDecimalPlace int
	var aInt, bInt, cInt, dInt int

	aInt = int(a)
	bInt = int(b)
	cInt = int(c)
	dInt = int(d)

	dInt += int(inc)
	if dInt > 255 {
		nextDecimalPlace = 1
	} else {
		nextDecimalPlace = 0
	}

	cInt += nextDecimalPlace
	if cInt > 255 {
		nextDecimalPlace = 1
	} else {
		nextDecimalPlace = 0
	}

	bInt += nextDecimalPlace
	if bInt > 255 {
		nextDecimalPlace = 1
	} else {
		nextDecimalPlace = 0
	}

	aInt += nextDecimalPlace
	if aInt > 255 {
		overflow = 1
	}

	elementA = byte(aInt)
	elementB = byte(bInt)
	elementC = byte(cInt)
	elementD = byte(dInt)

	return
}

func (el *IPv4Generator) IncCurrentIP() (
	err error,
) {

	var a, b, c, d, overflow byte
	var inc = byte(1)

	// (português): ip não foi inicializado
	if el.gatewayA != el.ipA {
		d = el.gatewayD
		c = el.gatewayC
		b = el.gatewayB
		a = el.gatewayA
	} else {
		d = el.ipD
		c = el.ipC
		b = el.ipB
		a = el.ipA
	}

	a, b, c, d, overflow = el.incIP(a, b, c, d, inc)
	if overflow != 0 {
		err = errors.New("the ip address is greater than the theoretical maximum allowed")
		return
	}

	err = el.verify(a, b, c, d)
	if err != nil {
		return
	}

	el.ipD = d
	el.ipC = c
	el.ipB = b
	el.ipA = a

	return
}

func (el *IPv4Generator) Init(
	gatewayA,
	gatewayB,
	gatewayC,
	gatewayD,
	subnetA,
	subnetB,
	subnetC,
	subnetD,
	subnetCidr byte,
) (
	err error,
) {
	el.reserved = make([][4]byte, 0)

	var overflow byte

	el.gatewayA = gatewayA
	el.gatewayB = gatewayB
	el.gatewayC = gatewayC
	el.gatewayD = gatewayD

	el.subnetA = subnetA
	el.subnetB = subnetB
	el.subnetC = subnetC
	el.subnetD = subnetD

	el.subnetCidr = subnetCidr

	if el.subnetA == 0 || el.gatewayA == 0 {
		err = errors.New("initialize IPv4Generator{} first")
		return
	}

	// 10.0.0.1 = ipA.ipB.ipC.ipD
	// = ipA*256^3 + ipB*256^2 + ipC*256^1 + ipD*256^0
	// = 10*256^3 + 0*256^2 + 0*256^1 + 1*256^0
	// = 167772161
	// = 10.0.0.1 As integer
	tmpA := int64(float64(el.gatewayA) * math.Pow(256.0, 3.0))
	tmpB := int64(float64(el.gatewayB) * math.Pow(256.0, 2.0))
	tmpC := int64(float64(el.gatewayC) * math.Pow(256.0, 1.0))
	tmpD := int64(float64(el.gatewayD) * math.Pow(256.0, 0.0))

	el.ipMinAddr = tmpA + tmpB + tmpC + tmpD

	subnetA, subnetB, subnetC, subnetD, overflow = el.incIP(subnetA, subnetB, subnetC, subnetD, el.cidrPrefixToDecimal(subnetCidr))
	if overflow != 0 {
		err = errors.New("the ip address is greater than the maximum theoretical ip allowed")
		return
	}

	// 10.0.0.1 = ipA.ipB.ipC.ipD
	// = ipA*256^3 + ipB*256^2 + ipC*256^1 + ipD*256^0
	// = 10*256^3 + 0*256^2 + 0*256^1 + 1*256^0
	// = 167772161
	// = 10.0.0.1 As integer
	tmpA = int64(float64(subnetA) * math.Pow(256.0, 3.0))
	tmpB = int64(float64(subnetB) * math.Pow(256.0, 2.0))
	tmpC = int64(float64(subnetC) * math.Pow(256.0, 1.0))
	tmpD = int64(float64(subnetD) * math.Pow(256.0, 0.0))

	el.cidrPrefixToDecimal(el.subnetCidr)
	el.ipMaxAddr = tmpA + tmpB + tmpC + tmpD

	el.maxA, el.maxB, el.maxC, el.maxD = el.int64ToIP(el.ipMaxAddr)
	el.ipA, el.ipB, el.ipC, el.ipD, overflow = el.incIP(el.gatewayA, el.gatewayB, el.gatewayC, el.gatewayD, 1)

	if overflow != 0 {
		err = errors.New("the ip address is greater than the maximum theoretical ip allowed")
		return
	}

	return
}

func (el *IPv4Generator) InitWithString(
	gateway string,
	subnet string,
) (
	err error,
) {
	var gatewayA, gatewayB, gatewayC, gatewayD byte
	var subnetA, subnetB, subnetC, subnetD, subnetCidr byte

	gatewayA, gatewayB, gatewayC, gatewayD, _, err = el.stringIPtoByte(gateway)
	if err != nil {
		return
	}

	subnetA, subnetB, subnetC, subnetD, subnetCidr, err = el.stringIPtoByte(subnet)
	if err != nil {
		return
	}

	err = el.Init(gatewayA, gatewayB, gatewayC, gatewayD, subnetA, subnetB, subnetC, subnetD, subnetCidr)
	return
}

func (el IPv4Generator) String() (currentIP string) {
	return strconv.FormatInt(int64(el.ipA), 10) + "." + strconv.FormatInt(int64(el.ipB), 10) + "." + strconv.FormatInt(int64(el.ipC), 10) + "." + strconv.FormatInt(int64(el.ipD), 10)
}
