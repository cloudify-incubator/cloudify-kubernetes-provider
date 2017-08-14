package utils

import (
	"reflect"
	"testing"
)

type testpair struct {
	input  []string
	output [2][]string
}

var tests = []testpair{
	{[]string{"cfy-go", "status"}, [2][]string{{"cfy-go", "status"}, {}}},
	{[]string{"cfy-go", "status", "-user", "admin"}, [2][]string{{"cfy-go", "status"}, {"-user", "admin"}}},
}

func TestCliArgumentsList(t *testing.T) {
	for _, pair := range tests {
		args, options := CliArgumentsList(pair.input)

		if !reflect.DeepEqual(args, pair.output[0]) {
			t.Error(
				"For", pair.input,
				"expected", pair.output[0],
				"got", args,
			)
		}

		if !reflect.DeepEqual(options, pair.output[1]) {
			t.Error(
				"For", pair.input,
				"expected", pair.output[1],
				"got", options,
			)
		}

	}
}
