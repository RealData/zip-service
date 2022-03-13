package main

import (
	"testing"
)

func TestFlagsError(t *testing.T) {

	var tests = []struct {
		scenario string
		args     []string
		err      error
	}{
		{"Compress/Extract flags are not set", []string{"zip-service"}, ceFlagNotSetErr},
		{"Compress and Extract flags are set", []string{"zip-service", "-c", "-e"}, ceFlagsSetErr},
		{"ZIP file flag not provided", []string{"zip-service", "-c"}, zipFileNotProvidedErr},
		{"Dir flag not provided", []string{"zip-service", "-c", "-f", "file.zip"}, dirNotProvidedErr},
		{"Number of threads set to 0", []string{"zip-service", "-c", "-f", "file.zip", "-d", "dir", "-n", "0"}, numThreadsErr},
		{"Method is not implemented", []string{"zip-service", "-c", "-f", "file.zip", "-d", "dir", "-n", "1", "-m", "zip"}, methodNotImplementedErr},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			err := parseFlags(test.args)
			if err != test.err {
				t.Errorf("Error should be '%s'", test.err)
			}
		})
	}

}
