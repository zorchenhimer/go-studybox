package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"bufio"
	"slices"
	"errors"

	"github.com/alexflint/go-arg"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/script"
)

type Arguments struct {
	Input string `arg:"positional,required"`
	Output string `arg:"positional"`
	StartAddr string `arg:"--start" default:"0x6000" help:"base address for the start of the script"`
	StatsFile string `arg:"--stats" help:"file to write some statistics to"`
	LabelFile string `arg:"--labels" help:"file containing address/label pairs"`
	CDL string `arg:"--cdl" help:"CodeDataLog json file"`
	Smart bool `arg:"--smart"`

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

	var cdl *script.CodeDataLog
	if args.CDL != "" {
		//fmt.Println("  CDL:", args.CDL)
		cdl, err = script.CdlFromJsonFile(args.CDL)
		if err != nil {
			//return fmt.Errorf("CDL Parse error: %w", err)
			cdl = nil
		}
	}

	var scr *script.Script
	if args.Smart {
		scr, err = script.SmartParseFile(args.Input, args.start, cdl)
	} else {
		scr, err = script.ParseFile(args.Input, args.start, cdl)
	}

	if err != nil {
		if errors.Is(err, script.ErrEarlyEOF) || errors.Is(err, script.ErrNavigation) {
			fmt.Println(err)
		} else {
			return fmt.Errorf("Script parse error: %w", err)
		}
	}

	if args.LabelFile != "" {
		labels, err := parseLabelFile(args.LabelFile)
		if err != nil {
			return fmt.Errorf("Labels parse error: %w", err)
		}

		for _, label := range labels {
			scr.Labels[label.Address] = label
		}
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
		//fmt.Fprintln(os.Stderr, w)
		if args.Output != "" {
			fmt.Fprintln(outfile, "; "+w)
		}
	}

	fmt.Fprintf(outfile, "; Start address: $%04X\n", scr.StartAddress)
	fmt.Fprintf(outfile, "; Stack address: $%04X\n\n", scr.StackAddress)

	slices.SortFunc(scr.Tokens, func(a, b *script.Token) int {
		if a.Offset < b.Offset { return -1 }
		if a.Offset > b.Offset { return 1 }
		return 0
	})

	//for _, lbl := range scr.Labels {
	//	fmt.Println(lbl)
	//}

	for _, token := range scr.Tokens {
		fmt.Fprintln(outfile, token.String(scr.Labels))
	}

	if args.StatsFile != "" {
		statfile, err := os.Create(args.StatsFile)
		if err != nil {
			return fmt.Errorf("Unable to create stats file: %w", err)
		}
		defer statfile.Close()

		_, err = scr.Stats().WriteTo(statfile)
		if err != nil {
			return fmt.Errorf("Error writing stats: %w", err)
		}
	}

	return nil
}

func parseLabelFile(filename string) ([]*script.Label, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	labels := []*script.Label{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.ReplaceAll(line, "\t", " ")
		parts := strings.Split(line, " ")

		parts = slices.DeleteFunc(parts, func(str string) bool {
			return str == ""
		})

		if len(parts) < 2 {
			fmt.Println("Ignoring", line)
			continue
		}

		if strings.HasPrefix(parts[0], "$") {
			parts[0] = "0x"+parts[0][1:]
		}

		addr, err := strconv.ParseInt(parts[0], 0, 32)
		if err != nil {
			fmt.Printf("Address parse error for %q: %s\n", line, err)
			continue
		}

		lbl := &script.Label{
			Name: parts[1],
			Address: int(addr),
		}

		if lbl.Name == "$" {
			lbl.Name = ""
		}

		if len(parts) > 2 {
			lbl.Comment = strings.Join(parts[2:], " ")
		}

		labels = append(labels, lbl)
	}

	return labels, nil
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
