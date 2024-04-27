package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/alexflint/go-arg"
	"git.zorchenhimer.com/Zorchenhimer/go-studybox/script"
)

type Arguments struct {
	Input string `arg:"positional,required"`
	Output string `arg:"positional"`
	StartAddr string `arg:"--start" default:"0x6000" help:"base address for the start of the script"`
	start int
}

func run(args *Arguments) error {
	if args.StartAddr == "" {
		return fmt.Errorf("start address cannot be empty")
	}

	if strings.HasPrefix(args.StartAddr, "$") {
		args.StartAddr = "0x"+args.StartAddr[1:]
	}

	val, err := strconv.ParseInt(args.StartAddr, 0, 32)
	if err != nil {
		return fmt.Errorf("invalid start address %q: %w", args.StartAddr, err)
	}

	args.start = int(val)

	scr, err := script.ParseFile(args.Input, args.start)
	if err != nil {
		return err
	}

	outfile := os.Stdout
	if args.Output != "" {
		outfile, err = os.Create(args.Output)
		if err != nil {
			return fmt.Errorf("unable to create output file: %w", err)
		}
		defer outfile.Close()
	}

	for _, w := range scr.Warnings {
		fmt.Fprintln(os.Stderr, w)
		if args.Output != "" {
			fmt.Fprintln(outfile, "; "+w)
		}
	}

	fmt.Fprintf(outfile, "; Start address: $%04X\n", scr.StartAddress)
	fmt.Fprintf(outfile, "; Stack address: $%04X\n\n", scr.StackAddress)

	for _, token := range scr.Tokens {
		fmt.Fprintln(outfile, token)
	}

	return nil
}

func main() {
	args := &Arguments{}
	arg.MustParse(args)

	err := run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
