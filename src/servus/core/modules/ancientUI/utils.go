package ancientUI

import (
	"fmt"
	"os"
)

// scanner - get user input and write to string. Throw panic if scan failed.
func scanner(writeTo *string) error {
	_, err := fmt.Fscan(os.Stdin, writeTo)
	return err
}

// ArgumentExists - if program running with args (ex: program.exe arg1 -arg2 --arg3 arg4) it func parse args and return true if argument specified in params exists.
func ArgumentExists(argument string) bool {
	// > 1 because first arg is program name.
	if len(os.Args) > 1 {
		for _, element := range os.Args {
			if element == argument {
				return true
			}
		}
	}
	return false
}
