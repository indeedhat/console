package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/posener/complete/v2"
)

type CmdConfig struct {
	Key     string `yaml:"key"`
	Usage   string `yaml:"usage"`
	Cmd     string `yaml:"cmd"`
	WorkDir string `yaml:"workDir"`
	Prompt  bool   `yaml:"prompt"`
}

func (c CmdConfig) WorkingDir() string {
	if c.WorkDir == "" {
		if dir, err := os.Getwd(); err == nil {
			return dir
		}
		return filepath.Dir(os.Args[0])
	}

	p, err := filepath.Abs(c.WorkDir)
	if err != nil {
		return filepath.Dir(os.Args[0])
	}

	return p
}

func (c CmdConfig) Run(args []string) error {
	if c.Prompt {
		if !cliPromptBool(fmt.Sprintf("Are you sure you want to run %s?", c.Key)) {
			return nil
		}
	}

	if script, ok := validShellScript(c.Cmd); ok {
		args = append([]string{script}, args...)
	} else {

		fh, err := os.CreateTemp(c.WorkingDir(), "run*.sh")
		if err != nil {
			return err
		}
		defer os.Remove(fh.Name())

		go func() {
			time.Sleep(100 * time.Millisecond)
			os.Remove(fh.Name())
		}()

		fh.WriteString(c.Cmd)
		args = append([]string{filepath.Base(fh.Name())}, args...)
	}

	cmd := exec.Command(shell(), args...)
	cmd.Dir = c.WorkingDir()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// CommandEntries list
type CliCommandEntries []CmdConfig

// Find a command entry by its key
func (c CliCommandEntries) Find(key string) *CmdConfig {
	for _, entry := range c {
		if key == entry.Key {
			return &entry
		}
	}

	return nil
}

func (c CliCommandEntries) Completions() map[string]*complete.Command {
	completions := make(map[string]*complete.Command)

	for _, entry := range c {
		completions[entry.Key] = &complete.Command{}
	}

	return completions
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

func validShellScript(s string) (string, bool) {
	if strings.Contains(s, "\n") {
		return s, false
	}

	if !strings.HasSuffix(s, ".sh") {
		return s, false
	}

	if !filepath.IsAbs(s) {
		s = filepath.Join(consoleBinDir(), s)
	}
	_, err := os.Stat(s)
	return s, err == nil
}

func consoleBinDir() string {
	dir := filepath.Dir(os.Args[0])
	if dir, err := filepath.Abs(dir); err == nil {
		return dir
	}

	return dir
}
