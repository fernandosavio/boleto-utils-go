package boleto

import (
	"testing"
	"time"
)

func TestFatorVencimentoToDate(t *testing.T) {
	testValues := []struct {
		fator       uint16
		expected    string
		expectedErr *Error
	}{
		{4469, "2010-01-01", nil},
		{4470, "2010-01-02", nil},
		{4468, "2034-08-22", nil},
		{4789, "2010-11-17", nil},
		{9999, "2025-02-21", nil},
		{1000, "2025-02-22", nil},
		{1002, "2025-02-24", nil},
		{1667, "2026-12-21", nil},
		{999, "", ErrFatorVencimento},
		{0, "", ErrFatorVencimento},
		{10000, "", ErrFatorVencimento},
	}

	for _, tt := range testValues {
		got, err := fatorVencimentoToDate(tt.fator)

		if tt.expectedErr != nil {
			if err != tt.expectedErr {
				t.Fatal("Expected error")
			}
		} else {

			if result := got.Format(time.DateOnly); result != tt.expected {
				t.Fatalf("Expected func(%d) = %q but got %q", tt.fator, tt.expected, result)
			}
		}
	}
}

func TestParseFatorVencimento(t *testing.T) {
	testValues := []struct {
		input       []byte
		expected    string
		expectedErr *Error
	}{
		{[]byte("4469"), "2010-01-01", nil},
		{[]byte("4470"), "2010-01-02", nil},
		{[]byte("4468"), "2034-08-22", nil},
		{[]byte("4789"), "2010-11-17", nil},
		{[]byte("9999"), "2025-02-21", nil},
		{[]byte("1000"), "2025-02-22", nil},
		{[]byte("1002"), "2025-02-24", nil},
		{[]byte("1667"), "2026-12-21", nil},
		{[]byte("0000"), "", nil},
		{[]byte("0999"), "", ErrFatorVencimento},
		{[]byte("10000"), "", ErrFatorVencimento},
		{[]byte("0"), "", ErrFatorVencimento},
		{[]byte("999"), "", ErrFatorVencimento},
	}

	for _, tt := range testValues {
		got, err := parseFatorVencimento(tt.input)

		if tt.expected == "" && tt.expectedErr == nil {
			if got != nil {
				t.Fatal("Expected nil")
			}
			continue
		}

		if tt.expectedErr != nil {
			if err != tt.expectedErr {
				t.Fatal("Expected error")
			}
		} else {
			if result := got.Format(time.DateOnly); result != tt.expected {
				t.Fatalf("Expected func(%d) = %q but got %q", tt.input, tt.expected, result)
			}
		}
	}
}

func TestParseMoeda(t *testing.T) {
	testValues := []struct {
		input       byte
		expected    CodMoeda
		expectedErr *Error
	}{
		{[]byte("0")[0], MoedaOutras, nil},
		{[]byte("9")[0], MoedaReal, nil},
		{[]byte("1")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("2")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("3")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("4")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("5")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("6")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("7")[0], CodMoeda(""), ErrCodigoMoeda},
		{[]byte("8")[0], CodMoeda(""), ErrCodigoMoeda},
	}

	for _, tt := range testValues {
		got, err := parseMoeda(tt.input)

		if got != tt.expected || err != tt.expectedErr {
			t.Fatalf("func(%q) = (%q, %q): got (%q, %q)", tt.input, tt.expected, tt.expectedErr, got, err)
		}
	}
}

func TestParseValor(t *testing.T) {
	testValues := []struct {
		input    []byte
		expected string
	}{
		{[]byte("9999999999"), "99999999.99"},
		{[]byte("0999999999"), "9999999.99"},
		{[]byte("0099999999"), "999999.99"},
		{[]byte("0009999999"), "99999.99"},
		{[]byte("0000999999"), "9999.99"},
		{[]byte("0000099999"), "999.99"},
		{[]byte("0000009999"), "99.99"},
		{[]byte("0000000999"), "9.99"},
		{[]byte("0000000099"), "0.99"},
		{[]byte("0000000009"), "0.09"},
		{[]byte("0000000000"), ""},
	}

	for _, tt := range testValues {
		got := parseValor(tt.input)

		if got != tt.expected {
			t.Fatalf("func(%q) expected %q got %q", tt.input, tt.expected, got)
		}
	}
}

func TestParseCodigBanco(t *testing.T) {
	testValues := []struct {
		input        []byte
		expectedName string
	}{
		{[]byte("001"), "Banco do Brasil"},
		{[]byte("237"), "Bradesco"},
		{[]byte("341"), "Itaú"},
	}

	for _, tt := range testValues {
		got := parseCodigoBanco(tt.input)

		if got.Codigo != string(tt.input) || got.Nome != tt.expectedName {
			t.Fatalf("%q should be %q, got %q", string(tt.input), tt.expectedName, got.Nome)
		}
	}
}

func TestParseSegmento(t *testing.T) {
	testValues := []struct {
		input        string
		expectedName string
		expectedErr  bool
	}{
		{"1", "Prefeituras", false},
		{"2", "Saneamento", false},
		{"3", "Energia elétrica e gás", false},
		{"4", "Telecomunicações", false},
		{"5", "Órgaos governamentais", false},
		{"6", "Carnês", false},
		{"7", "Multas de trânsito", false},
		{"9", "Uso exclusivo do banco", false},
		{"8", "", true},
		{"0", "", true},
	}

	for _, tt := range testValues {
		got, err := parseSegmento(tt.input[0])

		if tt.expectedErr {
			if !ErrInvalidSegmento.Is(err) {
				t.Fatal("Should be invalid")
			} else {
				continue
			}
		}

		if !tt.expectedErr && err != nil {
			t.Fatal("Should be valid")
		}

		if result := string(got.Codigo); result != tt.input {
			t.Fatalf("Expected Codigo %q got %q", tt.input, result)
		}

		if got.Nome != tt.expectedName {
			t.Fatalf("Expected Nome %q got %q", tt.expectedName, got.Nome)
		}
	}
}
