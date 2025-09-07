package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/rom"
)

type Arguments struct {
	Pack   *ArgPack   `arg:"subcommand:pack"`
	UnPack *ArgUnPack `arg:"subcommand:unpack"`
}

type ArgPack struct {
	Input string `arg:"positional,required"`
}

type ArgUnPack struct {
	Input   string `arg:"positional,required" help:".json metadata file"`
	NoAudio bool   `arg:"--no-audio" help:"Do not unpack the audio portion"`
	OutDir  string `arg:"--dir" help:"Base directory to unpack into (json file will be here)"`
}

func main() {
	args := &Arguments{}
	arg.MustParse(args)
	var err error

	switch {
	case args.Pack != nil:
		err = pack(args.Pack)
	case args.UnPack != nil:
		err = unpack(args.UnPack)
	default:
		fmt.Fprintln(os.Stderr, "Missing command")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func pack(args *ArgPack) error {
	if !strings.HasSuffix(args.Input, ".json") {
		return fmt.Errorf("Pack needs a json file as input")
	}

	//fmt.Println("-- Processing " + args.Input)
	sb, err := rom.Import(args.Input)
	if err != nil {
		return err
	}

	//outDir := filepath.Base(args.Input)
	//outDir = strings.ReplaceAll(outDir, ".json", "_output")

	//err = os.MkdirAll(outDir, 0777)
	//if err != nil {
	//	return err
	//}

	//err = sb.Export(outDir)
	//if err != nil {
	//	return err
	//}

	// TODO: put this in the json file?

	outname := args.Input[:len(args.Input)-len(".json")]+".studybox"
	fmt.Println(outname)
	err = sb.Write(outname)
	if err != nil {
		return err
	}

	return nil
}

func unpack(args *ArgUnPack) error {
	//fmt.Println("-- Processing " + file)
	if !strings.HasSuffix(args.Input, ".studybox") {
		return fmt.Errorf("Input needs to be a .studybox file.")
	}

	//outDir := filepath.Base(args.Input)
	outbase := filepath.Base(args.Input[:len(args.Input)-len(".studybox")])
	outdir := filepath.Dir(args.Input)
	if args.OutDir != "" {
		outdir = args.OutDir
	}
	outname := filepath.Join(outdir, outbase)
	fmt.Println(outname)
	//outDir = strings.ReplaceAll(outDir, ".studybox", "")

	err := os.MkdirAll(outname, 0777)
	if err != nil {
		return err
	}

	sb, err := rom.ReadFile(args.Input)
	if err != nil {
		return err
	}

	err = sb.Export(outname, !args.NoAudio)
	if err != nil {
		return err
	}

	return nil
}

//func main_old() {
//	if len(os.Args) < 3 {
//		fmt.Println("Missing command")
//		os.Exit(1)
//	}
//
//	matches := []string{}
//	for _, glob := range os.Args[2:len(os.Args)] {
//		m, err := filepath.Glob(glob)
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		matches = append(matches, m...)
//	}
//
//	if len(matches) == 0 {
//		fmt.Println("No files found!")
//	}
//
//	switch strings.ToLower(os.Args[1]) {
//	case "unpack":
//		for _, file := range matches {
//			fmt.Println("-- Processing " + file)
//			outDir := filepath.Base(file)
//			outDir = strings.ReplaceAll(outDir, ".studybox", "")
//
//			err := os.MkdirAll(outDir, 0777)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//
//			sb, err := rom.ReadFile(file)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//
//			err = sb.Export(outDir)
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	case "pack":
//		for _, file := range matches {
//			fmt.Println("-- Processing " + file)
//			sb, err := rom.Import(file)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//
//			outDir := filepath.Base(file)
//			outDir = strings.ReplaceAll(outDir, ".json", "_output")
//
//			err = os.MkdirAll(outDir, 0777)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//
//			err = sb.Export(outDir)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//
//			// TODO: put this in the json file?
//
//			err = sb.Write(outDir + ".studybox")
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//
//		}
//	}
//}
