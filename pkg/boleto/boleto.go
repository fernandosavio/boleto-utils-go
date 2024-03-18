package boleto

type CodMoeda string

const (
	MoedaReal   CodMoeda = "real"
	MoedaOutras CodMoeda = "outras"
)

type Segmento string
type CodigoSegmento byte

const (
	SegmentoPrefeituras          Segmento = Segmento("Prefeituras")
	SegmentoSaneamento           Segmento = Segmento("Saneamento")
	SegmentoEnergiaEGas          Segmento = Segmento("Energia elétrica e gás")
	SegmentoTelecom              Segmento = Segmento("Telecomunicações")
	SegmentoOrgaosGovernamentais Segmento = Segmento("Órgaos governamentais")
	SegmentoCarne                Segmento = Segmento("Carnês")
	SegmentoMultaTransito        Segmento = Segmento("Multas de trânsito")
	SegmentoExclusivoBanco       Segmento = Segmento("Uso exclusivo do banco")
)

const (
	CodSegmentoPrefeituras          CodigoSegmento = CodigoSegmento(b1)
	CodSegmentoSaneamento           CodigoSegmento = CodigoSegmento(b2)
	CodSegmentoEnergiaEGas          CodigoSegmento = CodigoSegmento(b3)
	CodSegmentoTelecom              CodigoSegmento = CodigoSegmento(b4)
	CodSegmentoOrgaosGovernamentais CodigoSegmento = CodigoSegmento(b5)
	CodSegmentoCarne                CodigoSegmento = CodigoSegmento(b6)
	CodSegmentoMultaTransito        CodigoSegmento = CodigoSegmento(b7)
	CodSegmentoExclusivoBanco       CodigoSegmento = CodigoSegmento(b9)
)

const (
	CodBarrasLength                 int = 44
	CobrancaLinhaDigitavelLength    int = 47
	ArrecadacaoLinhaDigitavelLength int = 48
)

func New(input string) (*Cobranca, *Arrecadacao, []Error) {
	return FromBytes([]byte(input))
}

func FromBytes(input []byte) (*Cobranca, *Arrecadacao, []Error) {
	if !hasOnlyNumbers(input) {
		return nil, nil, []Error{*ErrOnlyNumbers}
	}

	switch len(input) {
	case CodBarrasLength:
		codBarra := CodigoBarras(input)
		if input[0] == b8 {
			arr, errs := parseCodigoBarrasArrecadacao(codBarra)
			return nil, arr, errs
		} else {
			cob, errs := parseCodigoBarrasCobranca(codBarra)
			return cob, nil, errs
		}
	case CobrancaLinhaDigitavelLength:
		cob, errs := parseLinhaDigitavelCobranca(input)
		return cob, nil, errs
	case ArrecadacaoLinhaDigitavelLength:
		arr, errs := parseLinhaDigitavelArrecadacao(input)
		return nil, arr, errs
	default:
		return nil, nil, []Error{*ErrInvalidLength}
	}
}
