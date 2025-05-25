package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	"github.com/igolaizola/ff7book"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/peterbourgon/ff/v3/ffyaml"
)

// Build flags
var version = ""
var commit = ""
var date = ""

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("ff7book", flag.ExitOnError)

	return &ffcli.Command{
		ShortUsage: "ff7book [flags] <subcommand>",
		FlagSet:    fs,
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
		Subcommands: []*ffcli.Command{
			newVersionCommand(),
			newGenerateCommand(),
		},
	}
}

func newVersionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "ff7book version",
		ShortHelp:  "print version",
		Exec: func(ctx context.Context, args []string) error {
			v := version
			if v == "" {
				if buildInfo, ok := debug.ReadBuildInfo(); ok {
					v = buildInfo.Main.Version
				}
			}
			if v == "" {
				v = "dev"
			}
			versionFields := []string{v}
			if commit != "" {
				versionFields = append(versionFields, commit)
			}
			if date != "" {
				versionFields = append(versionFields, date)
			}
			fmt.Println(strings.Join(versionFields, " "))
			return nil
		},
	}
}

func newGenerateCommand() *ffcli.Command {
	cmd := "generate"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var cfg ff7book.Config
	fs.BoolVar(&cfg.Debug, "debug", false, "enable debug mode")
	fs.StringVar(&cfg.Key, "key", "", "Google API key")
	fs.StringVar(&cfg.Input, "input", "", "input file")
	fs.StringVar(&cfg.Output, "output", "", "output file")
	fs.StringVar(&cfg.Model, "model", "gemini-2.5-pro-preview-05-06", "model to use for generation")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("ff7book %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ffyaml.Parser),
			ff.WithEnvVarPrefix("FF7BOOK"),
		},
		ShortHelp: fmt.Sprintf("ff7book %s command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return ff7book.Generate(ctx, &cfg)
		},
	}
}
