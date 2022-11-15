package docker

import "strconv"

// SizeToString
//
// Português:
//
//	Formata um valor inteiro em uma string contendo o valor em Bytes, KBytes, MBytes e GBytes.
//
//	 Entrada:
//	   value: valor inteiro representando um tamanho de memória
//
//	 Saída:
//	   size: string contendo o valor em Bytes, KBytes, MBytes e GBytes
//
// English:
//
//	Format an integer value into a string containing the value in Bytes, KBytes, MBytes and GBytes.
//
//	Input:
//	  value: integer value representing a memory size
//
//	Output:
//	  size: string containing the value in Bytes, KBytes, MBytes and GBytes
func (e *ContainerBuilder) SizeToString(value int64) (size string) {

	if value == -1 {
		return "0 B"
	}

	if value > KGigaByte {
		return strconv.FormatFloat(float64(value)/KGigaByte, 'f', 1, 64) + " GB"
	}

	if value > KMegaByte {
		return strconv.FormatFloat(float64(value)/KMegaByte, 'f', 1, 64) + " MB"
	}

	if value > KKiloByte {
		return strconv.FormatFloat(float64(value)/KKiloByte, 'f', 1, 64) + " KB"
	}

	return strconv.FormatInt(value, 10) + " B"
}
