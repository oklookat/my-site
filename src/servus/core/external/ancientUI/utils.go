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
