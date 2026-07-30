package main

import "C"
import (
	"fmt"
	"os"
	"strings"
)

//export ap_main
func ap_main() {
	args := os.Args
	if len(args) > 1 && args[1] == "ap" {
		args = args[1:]
	}

	var (
		target   string
		bHeaders bool
	)

	for _, arg := range args[1:] {
		switch arg {
		case "-h":
			bHeaders = true
		case "--help":
			fmt.Printf("Usage: ap [-v] <target_url>\n\n")
			fmt.Printf("\t -h: Headers\n")
			fmt.Printf("Please provide a target url in the following format: 'https://example.com'\n")
			return
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Printf("Unknown option: %s\n", arg)
				return
			}

			if target == "" {
				target = arg
			} else {
				fmt.Printf("Unexpected argument: %s\n", arg)
				return
			}
		}
	}

	if target == "" || !strings.Contains(target, "http") {
		fmt.Println("Missing target URL")
		return
	}

	print(bHeaders)
}

func main() {}
