// piper is a CLI framework that does just what I want.
// It Manages a CLI for pipeline processes, where any number of defined tasks
// can be specified in sequence.
package piper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task func(interface{}, map[string]Flag, []string) interface{}

type Flag struct {
	Name        string
	Symbol      string
	Description string
}

type Command struct {
	Name        string
	Description string
	Args        []string
	Task        Task
}

type CLIApp struct {
	Name        string
	Description string
	Flags       []Flag
	Commands    []Command
}

func (c *CLIApp) PrintHelp() {
	fmt.Print("\n* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *\n")
	fmt.Print("\n" + c.Name + " - " + c.Description + "\n\n")
	fmt.Println("Usage:")
	fmt.Println(
		"   "+c.Name+" [global options] ",
		"[command [command options] [arguments...] ...]")
	fmt.Print("\n")
	if len(c.Flags) > 0 {
		fmt.Println("Global options:")
		for _, f := range c.Flags {
			if len(f.Symbol) > 0 {
				fmt.Print("   -" + f.Symbol)
			}
			if len(f.Description) > 0 {
				fmt.Print("  " + f.Description + "\n")
			}
		}
	}
	fmt.Print("\n")
	if len(c.Commands) > 0 {
		if len(c.Commands) > 0 {
			fmt.Println("Commands:")
			for _, s := range c.Commands {
				fmt.Print("   " + s.Name + " - " + s.Description + "\n")
				if len(s.Args) > 0 {
					fmt.Print("     args: " + s.Args[0])
					for _, a := range s.Args[1:] {
						fmt.Print(", " + a)
					}
					fmt.Print("\n")
				}
				fmt.Print("\n")
			}
		}
	}
	fmt.Print("* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *\n")
}

func (c *CLIApp) RegisterCommand(cmd Command) {
	c.Commands = append(c.Commands, cmd)
}

func (c *CLIApp) RegisterFlag(flag Flag) {
	c.Flags = append(c.Flags, flag)
}

// Parses command line arguments, constructing a pipeline of tasks from the
// subcommands along the way, returning an error if any issues are
// encountered.
// Once the arguments have been interpreted it executes the pipline.
func (c *CLIApp) Run() (err error) {
	flags := make(map[string]Flag)
	pipeline := make([]func(interface{}) interface{}, 0)
	i := 1
	for i < len(os.Args) {
		// skip whitespace
		if strings.TrimSpace(os.Args[i]) == "" {
			i++
			continue
		}
		var read int
		read, err = func(args []string) (read int, err error) {
			arg := args[0]
			if arg[:1] == "-" {
				// it's a flag
				for _, f := range c.Flags {
					if f.Symbol == arg[1:] {
						flags[f.Name] = f
						read = 1
						return
					}
				}
				err = errors.New("Unknown flag: " + arg)
			} else {
				// it's a task
				for _, t := range c.Commands {
					if t.Name == arg {
						if len(args) < len(t.Args)+1 {
							err = errors.New("Insufficient arguments provided for task " +
								arg + ", expected " + strconv.Itoa(len(t.Args)))
							return
						}
						task_args := args[1 : len(t.Args)+1]
						pipeline = append(pipeline, func(data interface{}) interface{} {
							data = t.Task(data, flags, task_args)
							return data
						})
						read = len(t.Args) + 1
						return
					}
				}
				err = errors.New("Unknown task: " + arg)
			}
			return
		}(os.Args[i:])

		if err != nil {
			return
		}
		i += read
	}

	if len(pipeline) == 0 {
		c.PrintHelp()
	}

	var data interface{}
	for _, stage := range pipeline {
		data = stage(data)
	}
	return
}
