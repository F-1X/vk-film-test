package verify

import (
	"testing"
	"vk/model"
)

var testCasesActor = []struct {
	name     string
	input    string
	expected bool
}{
	{
		"valid eng",
		"Mikle Duglas",
		true,
	},
	{
		"valid ru",
		"Mikle Duglas",
		true,
	},
	{
		"not valid other language ",
		"马克尔·古格拉斯",
		false,
	},
	{
		"not valid symbols",
		"!@#$%",
		false,
	},
	
}

func TestActor(t *testing.T) {

	for _, test := range testCasesActor {
		t.Run(test.input, func(t *testing.T) {
			result := Actor(test.input)
			if result != test.expected {
				t.Errorf("For input %s, expected %t, but got %t", test.input, test.expected, result)
			}
		})
	}
}

var testCasesGender = []struct {
	name     string
	input    string
	expected bool
}{

	{
		"valid male",
		"male",
		true,
	},
	{
		"valid female",
		"female",
		true,
	},
	{
		"not valid other all",
		"any",
		false,
	},
}

func TestGender(t *testing.T) {

	for _, test := range testCasesGender {
		t.Run(test.input, func(t *testing.T) {
			result := Gender(test.input)
			if result != test.expected {
				t.Errorf("For input %s, expected %t, but got %t", test.input, test.expected, result)
			}
		})
	}
}


var testCasesCreds = []struct {
	name     string
	input    model.Credentials
	expected bool
}{

	{
		"valid",
		model.Credentials{
			Username: "Gleb",
			Password: "qwerty1234",
		},

		true,
	},
}

func TestCreds(t *testing.T) {

	for _, test := range testCasesCreds {
		t.Run(test.name, func(t *testing.T) {
			result := Creds(test.input)
			if result != test.expected {
				t.Errorf("For input %s, expected %t, but got %t", test.input, test.expected, result)
			}
		})
	}
}
