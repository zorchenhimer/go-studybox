package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/audio"
	"git.zorchenhimer.com/Zorchenhimer/go-studybox/rom"
)

type Arguments struct {
	Input  string `arg:"positional,required"`
	Output string `arg:"positional,required"`

	BitRate int `arg:"--bit-rate", default:"4790"` // value found by trial and error
}

func run(args *Arguments) error {
	sbx, err := rom.ReadFile(args.Input)
	if err != nil {
		return err
	}

	output, err := os.Create(args.Output)
	if err != nil {
		return err
	}
	defer output.Close()

	audio.BitRate = args.BitRate

	err = audio.EncodeRom(output, sbx)
	if err != nil {
		return fmt.Errorf("Encode error: %w", err)
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
