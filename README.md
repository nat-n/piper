piper
=====

A simple and to the point CLI framework in go. The point being to allow
pipelining of commands like I needed that one time. Also generates help.

Example Usage
=============

### Creating a comand line interface

```go
package main

import (
  "fmt"
  "github.com/nat-n/piper"
  "strconv"
  str "string"
)

func main() {

  // create our CLI App
  cli := piper.CLIApp{
    Name:        "pipedream",
    Description: "dreamily pipes data through your tasks",
  }

  // configure it to accept a global flag that will be visible to all tasks
  cli.Flags = append(cli.Flags, piper.Flag{
    Name:        "verbose",
    Symbol:      "v",
    Description: "Verbose mode",
  })

  // Configure some commands. Each task must have a name, and a Task configured.
  // The Description is for the sake of help generation, and a number of
  // mandatory positional arguments may also be requested.
  // The task should return a reference to a datastructure which the next
  // task in the pipeline will recieve.

  cli.Commands = append(cli.Commands, piper.Command{
    Name:        "start",
    Description: "takes two words",
    Args:        []string{"first word", "second word"},
    Task: func(data *interface{}, flags map[string]piper.Flag, args []string) *interface{} {
      // Check for global flags like so
      if _, verbose := flags["verbose"]; verbose {
        fmt.Println(" - creating some data for the pipeline with those two words")
      }
      new_data := (interface{})(args)
      return &new_data
    },
  })

  cli.Commands = append(cli.Commands, piper.Command{
    Name:        "reverse",
    Description: "reverse all the words",
    Task: func(data *interface{}, flags map[string]piper.Flag, args []string) *interface{} {
      deref := *data
      words := deref.([]string)
      if _, verbose := flags["verbose"]; verbose {
        fmt.Println(" - reversing those words")
      }
      for i, word := range words {
        words[i] = str.Reverse(word)
      }
      return data
    },
  })

  cli.Commands = append(cli.Commands, piper.Command{
    Name:        "print",
    Description: "print whatever is in the pipeline",
    Args:        []string{"times"},
    Task: func(data *interface{}, flags map[string]piper.Flag, args []string) *interface{} {
      deref := *data
      words := deref.([]string)
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
      return data
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

```bash
$ pipedream -v  start Hello world!  reverse  print 3  reverse  print 1
 - creating some data for the pipeline with those two words
 - reversing those words
 - gonna print the words now
olleH !dlrow
olleH !dlrow
olleH !dlrow
 - reversing those words
 - gonna print the words now
Hello world!

$ pipedream

* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *

pipedream - dreamily pipes data through your tasks

Usage:
   pipedream [global options]  [command [command options] [arguments...] ...]

Global options:
   -v  Verbose mode

Commands:
   start - takes two words
     args: first word, second word

   reverse - reverse all the words

   print - print whatever is in the pipeline
     args: times

* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
```

License
=======

MIT. go nuts.