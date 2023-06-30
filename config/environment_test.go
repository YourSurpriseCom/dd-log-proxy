package config

import (
	"os"
	"testing"
)

type addValidateEnvVar struct {
	arg1, envValue, expected string
}
type validateEnvVariable struct {
	arg1, expected string
}

func Test_validateEnvVar(t *testing.T) {

	var addValidateEnvVarTests = []addValidateEnvVar{
		{"DEBUG_LEVEL", "", "info"},
		{"DEBUG_LEVEL", "debug", "debug"},
	}
	for _, test := range addValidateEnvVarTests {
		t.Setenv(test.arg1, test.envValue)
		validateEnvVar(test.arg1)
		output := os.Getenv(test.arg1)
		if output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}

}

func TestValidate(t *testing.T) {
	t.Setenv("DEBUG_LEVEL", "")
	t.Setenv("DD_SITE", "testing")
	t.Setenv("DD_API_KEY", "api-key")
	Validate()

	var validateEnvValidateTests = []validateEnvVariable{
		{"DEBUG_LEVEL", "info"},
		{"BATCH_SIZE", "50"},
		{"BATCH_WAIT_IN_SECONDS", "5"},
		{"PORT", "1053"},
	}

	for _, test := range validateEnvValidateTests {
		output := os.Getenv(test.arg1)
		if output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}

}
