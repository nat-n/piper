piper
=====

A simple and to the point CLI framework in go. The point being to allow
pipelining of commands like I needed that one time. Also generates help.

Example Usage
=============

### Creating a command line interface

```go
package main

import (
  "fmt"
  "github.com/nat-n/piper"
  "strconv"
  "strings"
)

func main() {

  // create our CLI App
  cli := piper.CLIApp{
    Name:        "pipedream",
    Description: "dreamily pipes data through your tasks",
  }

  // configure it to accept a global flag that will be visible to all tasks
  cli.RegisterFlag(piper.Flag{
    Name:        "verbose",
    Symbol:      "v",
    Description: "Verbose mode",
  })

  // Configure some commands. Each task must have a name, and a Task configured.
  // The Description is for the sake of help generation, and a number of
  // mandatory positional arguments may also be requested.
  // The task should return a reference to a data structure which the next
  // task in the pipeline will receive.

  cli.RegisterCommand(piper.Command{
    Name:        "start",
    Description: "takes two words",
    Args:        []string{"first word", "second word"},
    Task: func(data any, flags map[string]piper.Flag, args []string) (any, error) {
      // Check for global flags like so
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

  // Make it run
  err := cli.Run()

  // Be nice to the user when they screw up
  if err != nil {
    fmt.Println(err)
    cli.PrintHelp()
  }
}
```

### Invoking your command line interface:

```console
$ pipedream -v start Hello 'world!' upper print 3 lower print 1
 - creating some data for the pipeline with those two words
 - uppercasing those words
 - gonna print the words now
HELLO WORLD!
HELLO WORLD!
HELLO WORLD!
 - lowercasing those words
 - gonna print the words now
hello world!

$ pipedream

* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

pipedream - dreamily pipes data through your tasks

Usage:
   pipedream [global options] [command [arguments...] ...]

Global options:
   -v  Verbose mode

Commands:
   start - takes two words
     args: first word, second word

   upper - uppercase all the words

   lower - lowercase all the words

   print - print whatever is in the pipeline
     args: times

* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
```

License
=======

MIT. go nuts.
