package boleto

var ErrOnlyNumbers = &Error{"only_numbers"}
var ErrCodigoMoeda = &Error{"invalid_codigo_moeda"}
var ErrFatorVencimento = &Error{"invalid_fator_vencimento"}
var ErrInvalidLength = &Error{"invalid_length"}
var ErrInvalidDV = &Error{"invalid_digito_verificador"}
var ErrInvalidDAC = &Error{"invalid_dac"}
var ErrInvalidSegmento = &Error{"invalid_segmento"}
