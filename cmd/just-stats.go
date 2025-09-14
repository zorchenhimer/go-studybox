package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/fs"
	"slices"

	"github.com/alexflint/go-arg"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/script"
)

type Arguments struct {
	BaseDir string `arg:"positional,required"`
	Output  string `arg:"positional,required"`
}

type Walker struct {
	Found []string
	CDLs  []string
}

func (w *Walker) WalkFunc(path string, info fs.DirEntry, err error) error {
	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, "_scriptData.dat") {
		w.Found = append(w.Found, path)
	}

	if strings.HasSuffix(path, "_scriptData.cdl.json") {
		w.CDLs = append(w.CDLs, path)
	}

	return nil
}

func run(args *Arguments) error {
	w := &Walker{Found: []string{}}
	err := filepath.WalkDir(args.BaseDir, w.WalkFunc)
	if err != nil {
		return err
	}

	fmt.Printf("found %d scripts\n", len(w.Found))

	stats := make(script.Stats)

	for _, file := range w.Found {
		fmt.Println(file)
		var cdl *script.CodeDataLog
		cdlname := file[:len(file)-4]+".cdl.json"
		if slices.Contains(w.CDLs, cdlname) {
			fmt.Println("", cdlname)
			cdl, err = script.CdlFromJsonFile(cdlname)
			if err != nil {
				fmt.Println(" CDL read error:", err)
				cdl = nil
			}
		}

		scr, err := script.SmartParseFile(file, 0x6000, cdl)
		if err != nil {
			//if scr != nil {
			//	for _, token := range scr.Tokens {
			//		fmt.Println(token.String(scr.Labels))
			//	}
			//}
			fmt.Println(err)
			//return err
		}

		if scr != nil {
			stats.Add(scr.Stats())
		}
	}

	outfile, err := os.Create(args.Output)
	if err != nil {
		return err
	}
	defer outfile.Close()

	_, err = stats.WriteTo(outfile)
	if err != nil {
		return err
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
