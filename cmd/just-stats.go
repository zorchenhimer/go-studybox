package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/fs"

	"github.com/alexflint/go-arg"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/script"
)

type Arguments struct {
	BaseDir string `arg:"positional,required"`
	Output  string `arg:"positional,required"`
}

type Walker struct {
	Found []string
}

func (w *Walker) WalkFunc(path string, info fs.DirEntry, err error) error {
	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, "_scriptData.dat") {
		w.Found = append(w.Found, path)
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
		scr, err := script.ParseFile(file, 0x0000)
		if err != nil {
			if scr != nil {
				for _, token := range scr.Tokens {
					fmt.Println(token.String(scr.Labels))
				}
			}
			return err
		}

		stats.Add(scr.Stats())
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
