package boleto

var factors []byte = []byte{9, 8, 7, 6, 5, 4, 3, 2}

func mod11(input []byte, fallback byte, skip int) byte {

	hasSkip := 0
	if skip >= 0 {
		hasSkip = 1
	}

	var lenFactors = len(factors)
	// calcula a posição inicial do índice que será usado nos fatores
	var j int = lenFactors - ((len(input) - hasSkip) % lenFactors)

	var sum uint
	for i, v := range input {
		if hasSkip == 1 && i == skip {
			continue
		}
		sum += uint(v-b0) * uint(factors[j])
		j = (j + 1) % lenFactors
	}

	var dv byte = 11 - byte(sum%11)

	if dv > 9 {
		return fallback
	}

	return dv + b0
}

func mod10(input []byte) byte {
	var factor uint

	if len(input)&1 == 0 {
		factor = 1
	} else {
		factor = 2
	}

	var sum uint
	for _, v := range input {
		result := uint(v-b0) * factor
		if result < 10 {
			sum += result
		} else {
			sum += result % 10
			sum += result / 10
		}
		if factor == 1 {
			factor = 2
		} else {
			factor = 1
		}
	}

	var dv byte = (10 - byte(sum%10)) % 10
	return dv + b0
}
