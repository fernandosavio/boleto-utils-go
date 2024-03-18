package boleto

import "time"

type Cobranca struct {
	CodigoBarras   CodigoBarras
	LinhaDigitavel LinhaDigitavel
	Moeda          CodMoeda
	Banco          fieldBanco
	Vencimento     fieldVencimento
	Valor          string
	DAC            string
}

type fieldBanco struct {
	Codigo string
	Nome   string
}

func parseCodigoBarrasCobranca(input CodigoBarras) (*Cobranca, []Error) {
	fields := struct {
		banco []byte
		moeda byte
		dac   byte
		fator []byte
		valor []byte
	}{
		banco: input[:3],
		moeda: input[3],
		dac:   input[4],
		fator: input[5:9],
		valor: input[9:19],
	}
	errs := make([]Error, 0)

	// Fator vencimento
	fatorVencimento, err := parseFatorVencimento(fields.fator)
	var vencimento fieldVencimento

	if fatorVencimento == nil {
		vencimento = fieldVencimento{
			Fator: string(fields.fator),
			Data:  "",
		}
	} else {
		vencimento = fieldVencimento{
			Fator: string(fields.fator),
			Data:  fatorVencimento.Format(time.DateOnly),
		}
	}

	if err != nil {
		errs = append(errs, *err)
	}

	// Código moeda
	moeda, err := parseMoeda(fields.moeda)
	if err != nil {
		errs = append(errs, *err)
	}

	// Dígito verificador
	dv := mod11(input[:], b1, 4)

	if err != nil {
		errs = append(errs, *ErrInvalidLength)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Cobranca{
		CodigoBarras: CodigoBarras(input),
		Moeda:        moeda,
		Vencimento:   vencimento,
		Valor:        parseValor(fields.valor),
		Banco:        parseCodigoBanco(fields.banco),
		DAC:          string(dv),
	}, nil
}

func parseLinhaDigitavelCobranca(input LinhaDigitavel) (*Cobranca, []Error) {
	codBarras := linhaDigitavelToCodigoBarrasCobranca(input)
	dv1, dv2, dv3 := calculateDVs(codBarras)

	if input[9] != dv1 || input[20] != dv2 || input[31] != dv3 {
		return nil, []Error{*ErrInvalidDV}
	}

	return parseCodigoBarrasCobranca(codBarras)
}

func linhaDigitavelToCodigoBarrasCobranca(input LinhaDigitavel) CodigoBarras {
	barcode := make([]byte, CodBarrasLength)

	copy(barcode[0:4], input[0:4])
	copy(barcode[4:19], input[32:47])
	copy(barcode[19:24], input[4:9])
	copy(barcode[24:34], input[10:20])
	copy(barcode[34:44], input[21:31])

	return CodigoBarras(barcode)
}

func codigoBarrasToLinhaDigitavelCobranca(input CodigoBarras) LinhaDigitavel {
	digitableLine := make([]byte, CobrancaLinhaDigitavelLength)
	dv1, dv2, dv3 := calculateDVs(input)

	// Campo 1
	copy(digitableLine[0:4], input[0:4])
	copy(digitableLine[4:9], input[19:24])
	digitableLine[9] = dv1

	// Campo 2
	copy(digitableLine[10:20], input[24:34])
	digitableLine[20] = dv2

	// Campo 3
	copy(digitableLine[21:31], input[34:44])
	digitableLine[31] = dv3

	// DAC
	digitableLine[32] = input[4]

	// Campo 4
	copy(digitableLine[33:47], input[5:19])

	return LinhaDigitavel(digitableLine)
}

func calculateDVs(input CodigoBarras) (byte, byte, byte) {
	field1 := make([]byte, 9)
	copy(field1[0:4], input[0:4])
	copy(field1[4:9], input[19:24])

	dv1 := mod10(field1)
	dv2 := mod10(input[24:34])
	dv3 := mod10(input[34:44])

	return dv1, dv2, dv3
}
