package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nat-n/piper"
)

func main() {
	cli := piper.CLIApp{
		Name:        "pipedream",
		Description: "dreamily pipes data through your tasks",
	}

	cli.RegisterFlag(piper.Flag{
		Name:        "verbose",
		Symbol:      "v",
		Description: "Verbose mode",
	})

	cli.RegisterCommand(piper.Command{
		Name:        "start",
		Description: "takes two words",
		Args:        []string{"first word", "second word"},
		Task: func(data any, flags map[string]piper.Flag, args []string) (any, error) {
			if _, verbose := flags["verbose"]; verbose {
				fmt.Println(" - creating some data for the pipeline with those two words")
			}
			new_data := (any)(args)
			return new_data, nil
		},
	})

	cli.RegisterCommand(piper.Command{
		Name:        "upper",
		Description: "uppercase all the words",
		Task: func(data any, flags map[string]piper.Flag, args []string) (any, error) {
			words := data.([]string)
			if _, verbose := flags["verbose"]; verbose {
				fmt.Println(" - uppercasing those words")
			}
			for i, word := range words {
				words[i] = strings.ToUpper(word)
			}
			return any(words), nil
		},
	})

	cli.RegisterCommand(piper.Command{
		Name:        "lower",
		Description: "lowercase all the words",
		Task: func(data any, flags map[string]piper.Flag, args []string) (any, error) {
			words := data.([]string)
			if _, verbose := flags["verbose"]; verbose {
				fmt.Println(" - lowercasing those words")
			}
			for i, word := range words {
				words[i] = strings.ToLower(word)
			}
			return any(words), nil
		},
	})

	cli.RegisterCommand(piper.Command{
		Name:        "print",
		Description: "print whatever is in the pipeline",
		Args:        []string{"times"},
		Task: func(data any, flags map[string]piper.Flag, args []string) (any, error) {
			words := data.([]string)
			times, err := strconv.ParseInt(args[0], 0, 64)
			if err != nil {
				fmt.Println("Error: invalid argument for print")
			}
			if _, verbose := flags["verbose"]; verbose {
				fmt.Println(" - gonna print the words now")
			}
			for i := 0; i < int(times); i++ {
				for _, word := range words {
					fmt.Print(word, " ")
				}
				fmt.Print("\n")
			}
			return data, nil
		},
	})

	err := cli.Run()

	if err != nil {
		fmt.Println(err)
		cli.PrintHelp()
	}
}
