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

	return &Commands{
		Async: async,
		Debug: debug,
	}
}

/////////////////////////
//// Multi Commands /////
/////////////////////////

type Execers interface {
	Run() error
}

type Commands struct {
	Commands []Command
	Async    bool
	Debug    bool
}

func (c *Commands) Run(commandsCol []string) error {

	for _, commandStr := range commandsCol {

		args := strings.Split(commandStr, " ")
		if len(args) == 0 {
			log.Fatalln("Failed: empty command not allowed!")
		}

		command := &Command{
			args[0],
			args[1:],
			c.Debug,
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

	fmt.Println("Run: async")
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

	fmt.Println("Run: sequential")
	for _, stack := range c.Commands {
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
