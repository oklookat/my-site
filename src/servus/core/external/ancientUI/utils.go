package ancientUI

import (
	"fmt"
	"os"
)

// get user input and write to string.
func scanner(writeTo *string) error {
	_, err := fmt.Fscan(os.Stdin, writeTo)
	return err
}
