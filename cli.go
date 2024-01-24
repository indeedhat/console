package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

type CmdConfig struct {
	Key     string `yaml:"key"`
	Usage   string `yaml:"usage"`
	Cmd     string `yaml:"cmd"`
	WorkDir string `yaml:"workDir"`
}

func (c CmdConfig) WorkingDir() string {
	if c.WorkDir == "" {
		return filepath.Dir(os.Args[0])
	}

	p, err := filepath.Abs(c.WorkDir)
	if err != nil {
		return filepath.Dir(os.Args[0])
	}

	return p
}

func (c CmdConfig) Run(args []string) error {
	fh, err := os.CreateTemp(c.WorkingDir(), "run*.sh")
	if err != nil {
		return err
	}
	defer os.Remove(fh.Name())

	fh.WriteString(c.Cmd)
	args = append([]string{fh.Name()}, args...)

	cmd := exec.Command(shell(), args...)
	cmd.Dir = c.WorkingDir()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			os.Exit(err.ExitCode())
		}
	}

	return nil
}

// CommandEntries list
type CliCommandEntries []CmdConfig

// Find a command entry by its key
func (ces CliCommandEntries) Find(key string) *CmdConfig {
	for _, entry := range ces {
		if key == entry.Key {
			return &entry
		}
	}

	return nil
}

// CliUsage generator for the flags lib
func CliUsage(title, description string, register CliCommandEntries) func() {
	return func() {
		var builder strings.Builder

		builder.WriteString(title)
		builder.WriteByte('\n')

		if description != "" {
			builder.WriteString(description)
			builder.WriteByte('\n')
			builder.WriteByte('\n')
		}

		binName := filepath.Base(os.Args[0])

		builder.WriteString(fmt.Sprintf("USAGE:\n    ./%s <command>\n\n", binName))
		builder.WriteString("OPTIONS:\n")
		builder.WriteString("  -h, -help\n        Display this help message\n")

		fmt.Print(builder.String())
		flag.PrintDefaults()

		fmt.Print("\nCOMMANDS:\n")

		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		rows := make([]string, 0, len(register))
		for _, cmd := range register {
			for i, line := range strings.Split(cmd.Usage, "\n") {
				key := cmd.Key
				if i != 0 {
					key = "..."
				}

				rows = append(rows, fmt.Sprintf("    %s\t%s\n", green(key), yellow(line)))
			}
		}

		tbl := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', tabwriter.StripEscape)
		for _, row := range rows {
			fmt.Fprint(tbl, row)
		}

		tbl.Flush()
	}
}

func shell() string {
	if _, err := exec.LookPath("bash"); err == nil {
		return "bash"
	}

	if _, err := exec.LookPath("zsh"); err == nil {
		return "zsh"
	}

	return "sh"
}