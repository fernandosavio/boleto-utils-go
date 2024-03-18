package boleto

type CodigoBarras []byte
type LinhaDigitavel []byte

type fieldVencimento struct {
	Fator string
	Data  string
}

type Error struct {
	Code string
}

func (e *Error) Error() string {
	return e.Code
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && e.Code == t.Code
}
