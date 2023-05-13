package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tsuty/grace/internal"
)

func main() {
	args := internal.Args{}
	parser := flags.NewParser(&args, flags.HelpFlag|flags.PassDoubleDash)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); !ok || flagsErr.Type|flags.ErrHelp == flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return
	}

	fmt.Fprintln(os.Stderr, internal.StartServer(args))
}
