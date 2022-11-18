package builder

import (
	"errors"
	"github.com/helmutkemper/util"
	"log"
	"math"
	"testing"
)

func TestIPv4Generator_int64ToIP(t *testing.T) {
	ipGenerator := IPv4Generator{}
	var a, b, c, d byte
	var dA, dB, dC, dD byte
	var tmpA, tmpB, tmpC, tmpD int64
	var dSignificativePlace int
	var cSignificativePlace int
	var bSignificativePlace int
	var aSignificativePlace int
	var overflow int

	// 10.0.0.1
	dA = 10
	dB = 0
	dC = 0
	dD = 1

	tmpA = int64(float64(dA) * math.Pow(256.0, 3.0))
	tmpB = int64(float64(dB) * math.Pow(256.0, 2.0))
	tmpC = int64(float64(dC) * math.Pow(256.0, 1.0))
	tmpD = int64(float64(dD) * math.Pow(256.0, 0.0))
	startIPAsDecimal := tmpA + tmpB + tmpC + tmpD
	aSignificativePlace = int(dA)
	bSignificativePlace = int(dB)
	cSignificativePlace = int(dC)
	dSignificativePlace = int(dD)

	// 10.0.0.255
	dA = 10
	dB = 255
	dC = 255
	dD = 255

	tmpA = int64(float64(dA) * math.Pow(256.0, 3.0))
	tmpB = int64(float64(dB) * math.Pow(256.0, 2.0))
	tmpC = int64(float64(dC) * math.Pow(256.0, 1.0))
	tmpD = int64(float64(dD) * math.Pow(256.0, 0.0))
	endIPAsDecimal := tmpA + tmpB + tmpC + tmpD

	for ipAsInt := startIPAsDecimal; ipAsInt != endIPAsDecimal; ipAsInt += 1 {
		a, b, c, d = ipGenerator.int64ToIP(ipAsInt)
		if int(a) != aSignificativePlace {
			util.TraceToLog()
			t.FailNow()
		}

		if int(b) != bSignificativePlace {
			util.TraceToLog()
			t.FailNow()
		}

		if int(c) != cSignificativePlace {
			util.TraceToLog()
			t.FailNow()
		}

		if int(d) != dSignificativePlace {
			util.TraceToLog()
			t.FailNow()
		}

		dSignificativePlace += 1
		if dSignificativePlace > 255 {
			overflow = 1
			dSignificativePlace = 0
		} else {
			overflow = 0
		}

		cSignificativePlace += overflow
		if cSignificativePlace > 255 {
			overflow = 1
			cSignificativePlace = 0
		} else {
			overflow = 0
		}

		bSignificativePlace += overflow
		if bSignificativePlace > 255 {
			overflow = 1
			bSignificativePlace = 0
		} else {
			overflow = 0
		}

		aSignificativePlace += overflow
		if aSignificativePlace > 255 {
			overflow = 1
			aSignificativePlace = 0
		} else {
			overflow = 0
		}
	}
}

func TestIPv4Generator_Init(t *testing.T) {
	var err error

	g := IPv4Generator{}
	err = g.Init(
		//10.0.0.1
		10, 0, 0, 1,
		//10.0.0.0/4
		10, 0, 0, 0, 4,
	)
	if err != nil {
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA := g.gatewayA == 10
	caseB := g.gatewayB == 0
	caseC := g.gatewayC == 0
	caseD := g.gatewayD == 1

	if caseA && caseB && caseC && caseD == false {
		err = errors.New("gateway initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA = g.subnetA == 10
	caseB = g.subnetB == 0
	caseC = g.subnetC == 0
	caseD = g.subnetD == 0
	caseE := g.subnetCidr == 4

	if caseA && caseB && caseC && caseD && caseE == false {
		err = errors.New("subnet initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA = g.ipA == 10
	caseB = g.ipB == 0
	caseC = g.ipC == 0
	caseD = g.ipD == 2
	if caseA && caseB && caseC && caseD == false {
		err = errors.New("current ip initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	tmpA := int64(float64(g.gatewayA) * math.Pow(256.0, 3.0))
	tmpB := int64(float64(g.gatewayB) * math.Pow(256.0, 2.0))
	tmpC := int64(float64(g.gatewayC) * math.Pow(256.0, 1.0))
	tmpD := int64(float64(g.gatewayD) * math.Pow(256.0, 0.0))

	if g.ipMinAddr != tmpA+tmpB+tmpC+tmpD {
		err = errors.New("min ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	tmpA = int64(float64(g.subnetA) * math.Pow(256.0, 3.0))
	tmpB = int64(float64(g.subnetB) * math.Pow(256.0, 2.0))
	tmpC = int64(float64(g.subnetC) * math.Pow(256.0, 1.0))
	tmpD = int64(float64(g.subnetD) * math.Pow(256.0, 0.0))

	tmpE := int64(math.Pow(2.0, float64(int(g.subnetCidr))) - 1)

	if g.ipMaxAddr != tmpA+tmpB+tmpC+tmpD+tmpE {
		err = errors.New("max ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	var ipInt64 = g.ipMaxAddr
	var a, b, c, d byte
	d = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	c = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	b = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	a = byte(ipInt64 % 256)

	caseA = g.maxA == a
	caseB = g.maxA == b
	caseC = g.maxA == c
	caseD = g.maxA == d

	if caseA && caseB && caseC && caseD {
		err = errors.New("max ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

}

func TestIPv4Generator_Init_2(t *testing.T) {
	var err error

	g := IPv4Generator{}
	err = g.Init(
		//10.0.0.1
		10, 0, 10, 6,
		//10.0.0.0/4
		10, 0, 10, 5, 4,
	)
	if err != nil {
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA := g.gatewayA == 10
	caseB := g.gatewayB == 0
	caseC := g.gatewayC == 10
	caseD := g.gatewayD == 6

	if caseA && caseB && caseC && caseD == false {
		err = errors.New("gateway initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA = g.subnetA == 10
	caseB = g.subnetB == 0
	caseC = g.subnetC == 10
	caseD = g.subnetD == 5
	caseE := g.subnetCidr == 4

	if caseA && caseB && caseC && caseD && caseE == false {
		err = errors.New("subnet initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA = g.ipA == 10
	caseB = g.ipB == 0
	caseC = g.ipC == 10
	caseD = g.ipD == 7
	if caseA && caseB && caseC && caseD == false {
		err = errors.New("current ip initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	tmpA := int64(float64(g.gatewayA) * math.Pow(256.0, 3.0))
	tmpB := int64(float64(g.gatewayB) * math.Pow(256.0, 2.0))
	tmpC := int64(float64(g.gatewayC) * math.Pow(256.0, 1.0))
	tmpD := int64(float64(g.gatewayD) * math.Pow(256.0, 0.0))

	if g.ipMinAddr != tmpA+tmpB+tmpC+tmpD {
		err = errors.New("min ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	tmpA = int64(float64(g.subnetA) * math.Pow(256.0, 3.0))
	tmpB = int64(float64(g.subnetB) * math.Pow(256.0, 2.0))
	tmpC = int64(float64(g.subnetC) * math.Pow(256.0, 1.0))
	tmpD = int64(float64(g.subnetD) * math.Pow(256.0, 0.0))

	tmpE := int64(math.Pow(2.0, float64(int(g.subnetCidr))) - 1)

	if g.ipMaxAddr != tmpA+tmpB+tmpC+tmpD+tmpE {
		err = errors.New("max ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	var ipInt64 = g.ipMaxAddr
	var a, b, c, d byte
	d = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	c = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	b = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	a = byte(ipInt64 % 256)

	caseA = g.maxA == a
	caseB = g.maxA == b
	caseC = g.maxA == c
	caseD = g.maxA == d

	if caseA && caseB && caseC && caseD {
		err = errors.New("max ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}
}

func TestIPv4Generator_InitWithString(t *testing.T) {
	var err error

	g := IPv4Generator{}
	err = g.InitWithString(
		"10.0.0.1",
		"10.0.0.0/4",
	)
	if err != nil {
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA := g.gatewayA == 10
	caseB := g.gatewayB == 0
	caseC := g.gatewayC == 0
	caseD := g.gatewayD == 1

	if caseA && caseB && caseC && caseD == false {
		err = errors.New("gateway initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA = g.subnetA == 10
	caseB = g.subnetB == 0
	caseC = g.subnetC == 0
	caseD = g.subnetD == 0
	caseE := g.subnetCidr == 4

	if caseA && caseB && caseC && caseD && caseE == false {
		err = errors.New("subnet initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	caseA = g.ipA == 10
	caseB = g.ipB == 0
	caseC = g.ipC == 0
	caseD = g.ipD == 2
	if caseA && caseB && caseC && caseD == false {
		err = errors.New("current ip initialization error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	tmpA := int64(float64(g.gatewayA) * math.Pow(256.0, 3.0))
	tmpB := int64(float64(g.gatewayB) * math.Pow(256.0, 2.0))
	tmpC := int64(float64(g.gatewayC) * math.Pow(256.0, 1.0))
	tmpD := int64(float64(g.gatewayD) * math.Pow(256.0, 0.0))

	if g.ipMinAddr != tmpA+tmpB+tmpC+tmpD {
		err = errors.New("min ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	tmpA = int64(float64(g.subnetA) * math.Pow(256.0, 3.0))
	tmpB = int64(float64(g.subnetB) * math.Pow(256.0, 2.0))
	tmpC = int64(float64(g.subnetC) * math.Pow(256.0, 1.0))
	tmpD = int64(float64(g.subnetD) * math.Pow(256.0, 0.0))

	tmpE := int64(math.Pow(2.0, float64(int(g.subnetCidr))) - 1)

	if g.ipMaxAddr != tmpA+tmpB+tmpC+tmpD+tmpE {
		err = errors.New("max ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	var ipInt64 = g.ipMaxAddr
	var a, b, c, d byte
	d = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	c = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	b = byte(ipInt64 % 256)
	ipInt64 = ipInt64 / 256
	a = byte(ipInt64 % 256)

	caseA = g.maxA == a
	caseB = g.maxA == b
	caseC = g.maxA == c
	caseD = g.maxA == d

	if caseA && caseB && caseC && caseD {
		err = errors.New("max ip error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

}

func TestIPv4Generator_IncCurrentIP(t *testing.T) {
	var err error

	g := IPv4Generator{}
	err = g.InitWithString(
		"10.0.0.1",
		"10.0.0.0/4",
	)
	if err != nil {
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}

	for i := 0; i != 13; i += 1 {
		err = g.IncCurrentIP()
		if err != nil {
			util.TraceToLog()
			log.Printf("error: %v", err.Error())
			t.Fail()
			return
		}
		caseA := g.ipA == 10
		caseB := g.ipB == 0
		caseC := g.ipC == 0
		caseD := g.ipD == byte(3+i)

		if caseA && caseB && caseC && caseD == false {
			err = errors.New("current ip increment error")
			util.TraceToLog()
			log.Printf("error: %v", err.Error())
			t.Fail()
			return
		}
	}

	err = g.IncCurrentIP()
	if err == nil || err.Error() != "max allowed ip is 10.0.0.15" {
		err = errors.New("current ip increment error")
		util.TraceToLog()
		log.Printf("error: %v", err.Error())
		t.Fail()
		return
	}
}

func TestIPv4Generator_incIP(t *testing.T) {
	var err error
	var rA, rB, rC, rD, rOverflow byte

	var dSignificativePlace byte
	var cSignificativePlace byte
	var bSignificativePlace byte
	var aSignificativePlace byte
	var overflow byte

	g := IPv4Generator{}

	for a := 0; a != 256; a += 1 {
		for b := 0; b != 256; b += 1 {
			for c := 0; c != 256; c += 1 {
				for d := 0; d != 256; d += 1 {
					rA, rB, rC, rD, rOverflow = g.incIP(byte(a), byte(b), byte(c), byte(d), 1)

					dSignificativePlace = byte(d)
					cSignificativePlace = byte(c)
					bSignificativePlace = byte(b)
					aSignificativePlace = byte(a)

					dSignificativePlace += 1
					if dSignificativePlace > 255 {
						overflow = 1
						dSignificativePlace = 0
					} else {
						overflow = 0
					}

					cSignificativePlace += overflow
					if cSignificativePlace > 255 {
						overflow = 1
						cSignificativePlace = 0
					} else {
						overflow = 0
					}

					bSignificativePlace += overflow
					if bSignificativePlace > 255 {
						overflow = 1
						bSignificativePlace = 0
					} else {
						overflow = 0
					}

					aSignificativePlace += overflow
					if aSignificativePlace > 255 {
						overflow = 1
						aSignificativePlace = 0
					} else {
						overflow = 0
					}

					caseA := aSignificativePlace == rA
					caseB := bSignificativePlace == rB
					caseC := cSignificativePlace == rC
					caseD := dSignificativePlace == rD
					caseE := overflow == rOverflow

					if caseA && caseB && caseC && caseD && caseE == false {
						err = errors.New("ip increment error")
						util.TraceToLog()
						log.Printf("error: %v", err.Error())
						t.Fail()
						return
					}
				}
			}
		}
	}
}

/*
package main

import (
  "errors"
  "fmt"
  "net"
)

func main() {
  ip, err := incrementIP("10.0.0.6", "10.0.0.0/29")
  fmt.Printf("%v - %v\n", err, ip)
  ip, err = incrementIP("10.0.0.7", "10.0.0.0/29")
  fmt.Printf("%v - %v\n", err, ip)
}

func incrementIP(origIP, cidr string) (string, error) {
  ip := net.ParseIP(origIP)
  _, ipNet, err := net.ParseCIDR(cidr)
  if err != nil {
    return origIP, err
  }
  fmt.Printf("%+v", ipNet)
  for i := len(ip) - 1; i >= 0; i-- {
    ip[i]++
    if ip[i] != 0 {
      break
    }
  }
  if !ipNet.Contains(ip) {
    return origIP, errors.New("overflowed CIDR while incrementing IP")
  }
  return ip.String(), nil
}
*/
