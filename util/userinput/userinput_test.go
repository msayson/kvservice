package userinput

import (
	"testing"
)

func TestIsLegalCommand(t *testing.T) {
	testCases := []struct {
		input            string
		expectedLegality bool
	}{
		{"get(Hello123)", true},
		{"get(my_var_101)", true},
		{"get()", false},
		{"get(hi", false},
		{"get(", false},
		{"get", false},
		{"get(hello)1", false},
		{"set(Hello_123,MyVal123)", true},
		{"set(Hello_123,MyVal123,123)", false},
		{"set(Hello_123,)", false},
		{"set(Hello_123)", false},
		{"set()", false},
		{"set(hi,42", false},
		{"set(", false},
		{"set", false},
		{"set(hello)1", false},
		{"set(Hello_123,MyVal123)123", false},
		{"testset(Hello_123,MyVal123,NewVal123)", true},
		{"testset(Hello_123,a,b,c)", false},
		{"testset(Hello_123,a)", false},
		{"testset(Hello_123,)", false},
		{"testset(Hello_123)", false},
	}
	for _, test := range testCases {
		input := test.input
		expected := test.expectedLegality
		if IsLegalCommand(input) != expected {
			t.Errorf("IsLegalCommand(%s) returned %b, expected %b",
				input, !expected, expected)
		}
	}
}

func TestParseCommand(t *testing.T) {
	testCases := []struct {
		input string
		cmd   string
		args  []string
	}{
		{"get(Hello123)", "get", []string{"Hello123"}},
		{"set(Hello123,MyVal)", "set", []string{"Hello123", "MyVal"}},
		{"testset(Hello123,OldVal,NewVal)", "testset", []string{"Hello123", "OldVal", "NewVal"}},
		{" get(Hello123)   ", "get", []string{"Hello123"}}, //trims whitespace
	}
	for _, test := range testCases {
		input := test.input
		parsedCmd, err := ParseCommand(input)
		if err != nil {
			t.Errorf("ParseCommand(%s) returned unexpected error: %s", input, err.Error())
		}
		if parsedCmd.Command != test.cmd {
			t.Errorf("Expected ParseCommand(%s) to yield a \"%s\" command, instead received: %s", input, test.cmd, parsedCmd.Command)
		}
		if len(parsedCmd.Args) != len(test.args) {
			t.Errorf("Expected ParseCommand(%s) to yield args %s, instead received: %s", input, test.args, parsedCmd.Args)
		}
	}
}
