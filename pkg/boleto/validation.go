package boleto

import (
	"slices"
	"strconv"
	"time"
)

const (
	b0 byte = 0x30 + iota
	b1
	b2
	b3
	b4
	b5
	b6
	b7
	b8
	b9
)

const bdot byte = 0x2e

func hasOnlyNumbers(value []byte) bool {
	for _, char := range value {
		if char < b0 || char > b9 {
			return false
		}
	}
	return true
}

// Data base usada até 2025 (Fator(1000) == 03/07/2000)
var baseDate time.Time = time.Date(1997, 10, 7, 0, 0, 0, 0, time.UTC)

// Data base usada de 2025 em diante (Fator(1000) == 22/02/2025)
var newBaseDate time.Time = time.Date(2022, 5, 29, 0, 0, 0, 0, time.UTC)

func fatorVencimentoToDate(fator uint16) (time.Time, *Error) {
	if fator < 1000 || fator > 9999 {
		return time.Time{}, ErrFatorVencimento
	}
	var base time.Time

	// 2010-01-01
	if fator >= 4469 {
		base = baseDate
	} else {
		base = newBaseDate
	}

	return base.AddDate(0, 0, int(fator)), nil
}

func parseFatorVencimento(input []byte) (*time.Time, *Error) {
	if len(input) != 4 {
		return &time.Time{}, ErrFatorVencimento
	}

	if slices.Equal(input, []byte("0000")) {
		return nil, nil
	}

	value64, parseErr := strconv.ParseUint(string(input), 10, 16)

	if parseErr != nil {
		return &time.Time{}, ErrFatorVencimento
	}

	result, err := fatorVencimentoToDate(uint16(value64))

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func parseMoeda(input byte) (CodMoeda, *Error) {
	if input == b9 {
		return MoedaReal, nil
	}
	if input == b0 {
		return MoedaOutras, nil
	}
	return CodMoeda(""), ErrCodigoMoeda
}

func hasOnlyZeroes(input []byte) bool {
	for _, c := range input {
		if c != b0 {
			return false
		}
	}
	return true
}

func parseValor(input []byte) string {
	lenInput := len(input)

	// If input contains only zeroes, it means there is no value
	if hasOnlyZeroes(input) {
		return ""
	}

	// Creates a new byte array with 1 item extra to insert "."
	result := make([]byte, lenInput+1)

	// Copy the original and inserts the "." with the cents
	copy(result, input)
	result[lenInput-2] = bdot
	result[lenInput-1] = input[lenInput-2]
	result[lenInput] = input[lenInput-1]

	firstNonZero := 0

	for i, char := range result {
		if char != b0 {
			firstNonZero = i
			break
		}
	}

	// fixes ".99" to "0.99"
	if result[firstNonZero] == bdot {
		firstNonZero--
	}

	return string(result[firstNonZero:])
}

func parseCodigoBanco(input []byte) fieldBanco {
	return fieldBanco{
		Codigo: string(input),
		Nome:   bankCodes[string(input)],
	}
}

var segmentosMap = map[byte]fieldSegmento{
	b1: {
		Codigo: byte(CodSegmentoPrefeituras),
		Nome:   string(SegmentoPrefeituras),
	},
	b2: {
		Codigo: byte(CodSegmentoSaneamento),
		Nome:   string(SegmentoSaneamento),
	},
	b3: {
		Codigo: byte(CodSegmentoEnergiaEGas),
		Nome:   string(SegmentoEnergiaEGas),
	},
	b4: {
		Codigo: byte(CodSegmentoTelecom),
		Nome:   string(SegmentoTelecom),
	},
	b5: {
		Codigo: byte(CodSegmentoOrgaosGovernamentais),
		Nome:   string(SegmentoOrgaosGovernamentais),
	},
	b6: {
		Codigo: byte(CodSegmentoCarne),
		Nome:   string(SegmentoCarne),
	},
	b7: {
		Codigo: byte(CodSegmentoMultaTransito),
		Nome:   string(SegmentoMultaTransito),
	},
	b9: {
		Codigo: byte(CodSegmentoExclusivoBanco),
		Nome:   string(SegmentoExclusivoBanco),
	},
}

func parseSegmento(input byte) (fieldSegmento, *Error) {
	segmento, found := segmentosMap[input]

	if !found {
		return fieldSegmento{}, ErrInvalidSegmento
	}

	return segmento, nil
}
