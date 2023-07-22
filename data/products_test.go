package data

import "testing"

func Test_CheckValidation(t *testing.T) {
	p := &Product{
		Name:  "Adi",
		Price: 1.00,
		SKU:   "abc-abc-abc",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
