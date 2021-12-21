package argument

import (
	"os"
	"strings"
)

// Get - get argument.
func Get(arg string) *string {
	// > 1 because first arg is program name.
	var noArguments = !(len(os.Args) > 1)
	if noArguments {
		return nil
	}
	for _, element := range os.Args {
		// check is arg with value.
		var value = GetValue(&element)
		var withValue = value != nil
		if withValue {
			// if it is, remove value from arg.
			elementWithoutValue := strings.Replace(element, "="+*value, "", -1)
			if elementWithoutValue == arg {
				return &element
			}
		}
		if element == arg {
			return &element
		}
	}
	return nil
}

// GetValue - get value from argument. Argument must contains equal (=) symbol. Provide arg using Get() function.
func GetValue(arg *string) *string {
	if arg == nil {
		return nil
	}
	sliced := strings.SplitN(*arg, "=", 2)
	if len(sliced) < 2 {
		return nil
	}
	value := sliced[1]
	return &value
}
