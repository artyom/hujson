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

	"github.com/tailscale/hujson"
)

func main() {
	log.SetFlags(0)
	var replace bool
	flag.BoolVar(&replace, "w", replace,
		"overwrite sources instead of printing (re)formatted text")
	flag.Parse()
	if err := run(replace, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}

func run(replace bool, names ...string) error {
	if len(names) == 0 {
		return errors.New("nothing to do")
	}
	for _, name := range names {
		if err := formatFile(replace, name); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func formatFile(replace bool, name string) error {
	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	val, err := hujson.Parse(b)
	if err != nil {
		return err
	}
	val.Format()
	if !replace {
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
