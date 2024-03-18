package boleto

type Arrecadacao struct {
	CodigoBarras   CodigoBarras
	LinhaDigitavel LinhaDigitavel
	Segmento       fieldSegmento
}

type fieldSegmento struct {
	Codigo byte
	Nome   string
}

type fieldConvenio struct {
	Codigo string
	Nome   string
}

func parseCodigoBarrasArrecadacao(input CodigoBarras) (*Arrecadacao, []Error) {
	fields := struct {
		tipoValor  byte
		segmento   byte
		dac        byte
		valor      []byte
		convenio   []byte
		campoLivre []byte
	}{
		segmento:   input[1],
		tipoValor:  input[2],
		dac:        input[3],
		valor:      input[4:15],
		convenio:   input[15:19],
		campoLivre: input[19:44],
	}

	if fields.segmento == b6 {
		fields.convenio = input[15:23]
		fields.campoLivre = input[23:44]
	}

	errs := make([]Error, 0)

	// Segmento
	segmento, err := parseSegmento(fields.segmento)

	if err != nil {
		errs = append(errs, *err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Arrecadacao{
		Segmento: segmento,
	}, nil
}

func parseLinhaDigitavelArrecadacao(input LinhaDigitavel) (*Arrecadacao, []Error) {
	return &Arrecadacao{}, nil
}

func linhaDigitavelToCodigoBarrasArrecadacao(input []byte) []byte {
	result := make([]byte, 44, 44)
	return result
}
