package ancientUI

import (
	"fmt"
	"os"
)

func scanner(writeTo *string) {
	_, err := fmt.Fscan(os.Stdin, writeTo)
	if err != nil {
		panic(err)
	}
}

func ReadArg(argument string) bool{
	if len(os.Args) > 1 {
		for _, element := range os.Args {
			if element == argument {
				return true
			}
		}
	}
	return false
}
