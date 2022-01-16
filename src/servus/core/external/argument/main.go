package argument

import (
	"os"
	"strings"
)

// contains argument info.
type Info struct {
	// argument without value. Like: -username.
	Clean string
	// argument with value. Like: -username=1234.
	Dirty string
	// argument value. Like: 1234.
	Value *string
}

// get argument.
func Get(arg string) *Info {
	// > 1 because first arg is program name.
	var noArguments = !(len(os.Args) > 1)
	if noArguments {
		return nil
	}
	var info = &Info{}
	for _, dirtyArg := range os.Args {
		// if same args.
		info.Dirty = dirtyArg
		if dirtyArg == arg {
			info.Clean = dirtyArg
			return info
		}
		// ok, not same. But maybe dirtyArg has value?
		var value = getValue(&dirtyArg)
		if value == nil {
			// no.
			continue
		}
		// yes. Remove value from arg for compare.
		argWithoutValue := strings.Replace(dirtyArg, "="+*value, "", -1)
		// now maybe dirty arg without value and arg same?
		if argWithoutValue == arg {
			// yes.
			info.Clean = argWithoutValue
			info.Value = value
			return info
		}
	}
	return nil
}

// get value from argument. Argument must contains equals (=) sign.
func getValue(arg *string) *string {
	if arg == nil {
		// no arg.
		return nil
	}
	// split only up to the first equal sign
	// output like: ["-username", "hello=world-1234"]
	sliced := strings.SplitN(*arg, "=", 2)
	if len(sliced) < 2 {
		// arg without equals sign.
		return nil
	}
	value := sliced[1]
	return &value
}
