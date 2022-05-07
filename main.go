// Command hujson reformats [JWCC] files according to opinionated formatting
//
// It is a shallow wrapper of github.com/tailscale/hujson library.
//
// [JWCC]: https://nigeltao.github.io/blog/2021/json-with-commas-comments.html
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/artyom/hujson/internal/diff"
	"github.com/tailscale/hujson"
)

func main() {
	log.SetFlags(0)
	var args runArgs
	flag.BoolVar(&args.replace, "w", args.replace,
		"overwrite sources instead of printing formatted text")
	flag.BoolVar(&args.diff, "d", args.diff,
		"display diffs instead of printing formatted text")
	flag.Parse()
	if err := run(args, flag.Args()...); err != nil {
		if errors.Is(err, errNonEmptyDiff) {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}

type runArgs struct {
	replace, diff bool
}

func run(args runArgs, names ...string) error {
	if len(names) == 0 {
		return errors.New("nothing to do")
	}
	var hasDiffs bool
	for _, name := range names {
		if err := formatFile(args, name); err != nil {
			if errors.Is(err, errNonEmptyDiff) {
				hasDiffs = true
				continue
			}
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	if hasDiffs {
		return errNonEmptyDiff
	}
	return nil
}

func formatFile(args runArgs, name string) error {
	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	val, err := hujson.Parse(b)
	if err != nil {
		return err
	}
	val.Format()
	if args.diff {
		if res := diff.Diff(name, b, name+".new", []byte(val.String())); res != nil {
			fmt.Printf("%s", res)
			return errNonEmptyDiff
		}
		return nil
	}
	if !args.replace {
		fmt.Print(val.String())
		return nil
	}
	return os.WriteFile(name, []byte(val.String()), 0666)
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Usage: hujson [flags] file1.jwcc [file2.jwcc]...\n")
		flag.PrintDefaults()
	}
}

var errNonEmptyDiff = errors.New("non-empty diff")
