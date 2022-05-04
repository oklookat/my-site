package cmd

import "os"

// close program by code 0 or 1.
func exit(withError bool) {
	var code = 0
	if withError {
		code = 1
	}
	os.Exit(code)
}

// check error, log, exit.
func afterCommand(err error) {
	if err != nil {
		cfg.Logger.Error(err.Error())
		exit(true)
	}
	exit(false)
}
