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

func New(opts ...bool) *Commands {

	var async, debug, quiet bool

	if len(opts) == 1 {
		async = opts[0]
		debug = false
		quiet = false
	}

	if len(opts) == 2 {
		async = opts[0]
		debug = opts[1]
		quiet = false
	}

	if len(opts) == 3 {
		async = opts[0]
		debug = opts[1]
		quiet = opts[2]
	}

	return &Commands{
		Async: async,
		Debug: debug,
		Quiet: quiet,
	}
}

// Multi Commands

type Execer interface {
	Run([]string) error
}

type Commands struct {
	Commands []Command
	Async    bool
	Debug    bool
	Quiet    bool
}

func (c *Commands) Run(commandsCol []string) error {

	// reset collection
	c.Commands = c.Commands[:0]

	for _, commandStr := range commandsCol {

		args := strings.Split(commandStr, " ")
		if len(args) == 0 {
			log.Fatalln("Failed: empty command not allowed!")
		}
		command := &Command{
			Name:  args[0],
			Args:  args[1:],
			Debug: c.Debug,
			Quiet: c.Quiet,
		}
		c.Commands = append(c.Commands, *command)
	}

	if c.Async == true {
		c.RunAsync()
	} else {
		c.RunSequential()
	}

	return nil
}

func (c *Commands) RunAsync() error {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(len(c.Commands))

	if c.Quiet {
		fmt.Println("Run: async")
	}
	for _, stack := range c.Commands {
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

func (c *Commands) RunSequential() error {

	if c.Quiet {
		fmt.Println("Run: sequential")
	}
	for _, stack := range c.Commands {
		err := stack.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

//  Single Command

type Command struct {
	Name  string
	Args  []string
	Debug bool
	Quiet bool
}

func (c Command) Run() error {
	return c.run()
}

func (c Command) run() error {

	cmd := *exec.Command(c.Name, c.Args...)

	if c.Quiet {
		fmt.Println("Run: " + c.Name + " " + strings.Join(c.Args, " "))
	}
	if c.Debug == true {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
