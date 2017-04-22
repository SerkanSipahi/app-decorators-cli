package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func New(commandsCol []string, opts ...bool) *Commands {

	var async bool
	var debug bool

	if len(opts) == 1 {
		async = opts[0]
		debug = false
	}

	if len(opts) == 2 {
		async = opts[0]
		debug = opts[1]
	}

	commands := &Commands{
		Async: async,
	}
	for _, commandStr := range commandsCol {

		args := strings.Split(commandStr, " ")
		if len(args) == 0 {
			log.Fatalln("Failed: empty command not allowed!")
		}

		command := &Command{
			args[0],
			args[1:],
			debug,
		}

		commands.Stack = append(commands.Stack, *command)
	}
	return commands
}

/////////////////////////
//// Multi Commands /////
/////////////////////////

type ICommands interface {
	Run() error
	RunAsync() error
	RunSequential() error
}

type Commands struct {
	Stack []Command
	Async bool
	Debug bool
}

func (c Commands) Run() error {

	var commands []Command = c.Stack
	if c.Async == true {
		c.RunAsync(commands)
	} else {
		c.RunSequential(commands)
	}

	return nil
}

func (c Commands) RunAsync(commands []Command) error {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(len(c.Stack))

	fmt.Println("Run: async")
	for _, stack := range commands {
		go func(stack Command) {

			defer wg.Done()
			err := stack.Run()

			if err != nil {
				log.Fatal(err)
			}
		}(stack)
	}

	wg.Wait()

	return nil
}

func (c Commands) RunSequential(commands []Command) error {

	fmt.Println("Run: sequential")
	for _, stack := range commands {
		err := stack.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

/////////////////////////
//// Single Command /////
/////////////////////////

type ICommand interface {
	Run() error
}

type Command struct {
	Name  string
	Args  []string
	Debug bool
}

func (c Command) Run() error {
	return c.run()
}

func (c Command) run() error {

	cmd := *exec.Command(c.Name, c.Args...)
	fmt.Println("Run: " + c.Name + " " + strings.Join(c.Args, " "))
	if c.Debug == true {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
