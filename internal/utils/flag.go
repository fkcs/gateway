package utils

import (
	"fmt"
	"os"
)

func FlagMustBePresent(program, key string, p *string) {
	if p == nil || *p == "" {
		fmt.Printf("command line argument '%s' is missing: run '%s --help' for more information\n", key, program)
		os.Exit(1)
	}
}
