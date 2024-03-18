package boleto_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/fernandosavio/boleto-utils-go/pkg/boleto"
)

func TestOnlyNumbersValidation(t *testing.T) {
	value := strings.Repeat("0", 43) + "A"

	_, _, errs := boleto.New(value)

	hasError := func(e boleto.Error) bool {
		return boleto.ErrOnlyNumbers.Is(&e)
	}

	if !slices.ContainsFunc(errs, hasError) {
		t.Fatalf("%q should return ErrOnlyNumbers", value)
	}
}

func TestLengthValidation(t *testing.T) {
	invalidLengths := []int{1, 10, 20, 30, 40, 43, 46, 45, 49, 50}

	for _, count := range invalidLengths {
		value := strings.Repeat("9", count)

		_, _, errs := boleto.New(value)

		hasError := func(e boleto.Error) bool {
			return boleto.ErrInvalidLength.Is(&e)
		}

		if !slices.ContainsFunc(errs, hasError) {
			t.Fatalf("%q should return ErrInvalidLength", value)
		}
	}
}

func TestCodigoMoedaValidation(t *testing.T) {
	// Valores inválidos
	invalidValues := []string{
		"11111444455555555556666666666666666666666666",
		"11121444455555555556666666666666666666666666",
		"11131444455555555556666666666666666666666666",
		"11141444455555555556666666666666666666666666",
		"11151444455555555556666666666666666666666666",
		"11161444455555555556666666666666666666666666",
		"11171444455555555556666666666666666666666666",
		"11181444455555555556666666666666666666666666",
	}

	for _, value := range invalidValues {
		_, _, errs := boleto.New(value)

		hasError := func(e boleto.Error) bool {
			return boleto.ErrCodigoMoeda.Is(&e)
		}

		if !slices.ContainsFunc(errs, hasError) {
			t.Fatalf("%q should return ErrCodigoMoeda", value)
		}
	}

	// CodigoMoeda = real
	moedaReal := "11191444455555555556666666666666666666666666"
	cob, arrec, err := boleto.New(moedaReal)

	if err != nil {
		t.Fatalf("%q should be valid", moedaReal)
	}

	if arrec != nil {
		t.Fatalf("%q should be Cobranca", moedaReal)
	}

	if cob.Moeda != boleto.MoedaReal {
		t.Fatalf("CodigoMoeda should be %q", boleto.MoedaReal)
	}

	// CodigoMoeda = outras
	moedaOutras := "11105444455555555556666666666666666666666666"
	cob, arrec, err = boleto.New(moedaOutras)

	if err != nil {
		t.Fatalf("%q should be valid", moedaOutras)
	}

	if arrec != nil {
		t.Fatalf("%q should be Cobranca", moedaOutras)
	}

	if cob.Moeda != boleto.MoedaOutras {
		t.Fatalf("CodigoMoeda should be %q", boleto.MoedaOutras)
	}
}

func TestCodigoBancoValidation(t *testing.T) {
	// Valores inválidos
	invalidValues := []struct {
		value    string
		expected string
	}{
		{"11191444455555555556666666666666666666666666", "111"},
		{"99996444455555555556666666666666666666666666", "999"},
		{"12395444455555555556666666666666666666666666", "123"},
		{"66691444455555555556666666666666666666666666", "666"},
		{"00091444455555555556666666666666666666666666", "000"},
	}

	for _, tt := range invalidValues {
		cob, arr, err := boleto.New(tt.value)

		if err != nil {
			t.Fatalf("%q should be valid", tt.value)
		}

		if arr != nil {
			t.Fatalf("%q should be Cobranca", tt.value)
		}

		if string(cob.Banco.Codigo) != tt.expected {
			t.Fatalf("%q: CodigoBanco should be %q", tt.value, tt.expected)
		}
	}
}

func TestVencimentoValidation(t *testing.T) {
	// Valores válidos
	testValues := []struct {
		value  string
		factor string
		date   string
	}{
		{
			value:  "11196000055555555556666666666666666666666666",
			factor: "0000",
			date:   "",
		},
		{
			value:  "11199100055555555556666666666666666666666666",
			factor: "1000",
			date:   "2025-02-22",
		},
		{
			value:  "11191100255555555556666666666666666666666666",
			factor: "1002",
			date:   "2025-02-24",
		},
		{
			value:  "11196166755555555556666666666666666666666666",
			factor: "1667",
			date:   "2026-12-21",
		},
		{
			value:  "11198478955555555556666666666666666666666666",
			factor: "4789",
			date:   "2010-11-17",
		},
		{
			value:  "11193999955555555556666666666666666666666666",
			factor: "9999",
			date:   "2025-02-21",
		},
		{
			value:  "75696903800002500001434301033723400014933001",
			factor: "9038",
			date:   "2022-07-06",
		},
		{
			value:  "00191667900002434790000002656973019362470618",
			factor: "6679",
			date:   "2016-01-20",
		},
		{
			value:  "00195586200000773520000002464206011816073018",
			factor: "5862",
			date:   "2013-10-25",
		},
		{
			value:  "75592896700003787000003389850761252543475984",
			factor: "8967",
			date:   "2022-04-26",
		},
		{
			value:  "23791672000003249052028269705944177105205220",
			factor: "6720",
			date:   "2016-03-01",
		},
		{
			value:  "23791672000003097902028060007024617500249000",
			factor: "6720",
			date:   "2016-03-01",
		},
	}

	for _, tt := range testValues {
		cob, arr, errs := boleto.New(tt.value)

		if len(errs) > 0 {
			t.Fatalf("%q should be valid: %v", tt.value, errs)
		}

		if arr != nil {
			t.Fatalf("%q should be Cobranca", tt.value)
		}

		if string(cob.Vencimento.Fator) != tt.factor {
			t.Fatalf("%q: Vencimento.Fator should be %q", tt.value, tt.factor)
		}

		if cob.Vencimento.Data != tt.date {
			t.Fatalf("%q: Vencimento.Data should be %q", tt.value, tt.date)
		}
	}

	invalidValues := []string{
		"11196000155555555556666666666666666666666666",
		"11196099955555555556666666666666666666666666",
	}

	for _, value := range invalidValues {
		_, _, errs := boleto.New(value)

		hasError := func(e boleto.Error) bool {
			return boleto.ErrFatorVencimento.Is(&e)
		}

		if !slices.ContainsFunc(errs, hasError) {
			t.Fatalf("%q should return ErrFatorVencimento", value)
		}
	}
}

func TestValorValidation(t *testing.T) {
	// Valores válidos
	testValues := []struct {
		value    string
		expected string
	}{
		{
			value:    "11191444455555555556666666666666666666666666",
			expected: "55555555.55",
		},
		{
			value:    "11196444499999999996666666666666666666666666",
			expected: "99999999.99",
		},
		{
			value:    "11196444400000099996666666666666666666666666",
			expected: "99.99",
		},
		{
			value:    "11193444400000000006666666666666666666666666",
			expected: "",
		},
	}

	for _, tt := range testValues {
		cob, arr, err := boleto.New(tt.value)

		if err != nil {
			t.Fatalf("%q should be valid", tt.value)
		}

		if arr != nil {
			t.Fatalf("%q should be Cobranca", tt.value)
		}

		if cob.Valor != tt.expected {
			t.Fatalf("%q: Valor should be %q", tt.value, tt.expected)
		}
	}
}

func TestDigitoVerificadorValidation(t *testing.T) {
	// Valores válidos
	testValues := []struct {
		value    string
		expected string
	}{
		{
			value:    "11191444455555555556666666666666666666666666",
			expected: "1",
		},
		{
			value:    "10499898100000214032006561000100040099726390",
			expected: "9",
		},
		{
			value:    "75696903800002500001434301033723400014933001",
			expected: "6",
		},
		{
			value:    "00191667900002434790000002656973019362470618",
			expected: "1",
		},
		{
			value:    "00195586200000773520000002464206011816073018",
			expected: "5",
		},
		{
			value:    "75592896700003787000003389850761252543475984",
			expected: "2",
		},
		{
			value:    "23791672000003249052028269705944177105205220",
			expected: "1",
		},
		{
			value:    "23791672000003097902028060007024617500249000",
			expected: "1",
		},
		{
			value:    "11191100255555555556666666666666666666666666",
			expected: "1",
		},
	}

	for _, tt := range testValues {
		cob, arr, err := boleto.New(tt.value)

		if err != nil {
			t.Fatalf("%q should be valid", tt.value)
		}

		if arr != nil {
			t.Fatalf("%q should be Cobranca", tt.value)
		}

		if cob.DAC != tt.expected {
			t.Fatalf("%q: DAC should be %q", tt.value, tt.expected)
		}
	}
}
