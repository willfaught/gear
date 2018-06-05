package cli

import "fmt"

type Cmd struct {
	Name  string
	Run   func([]string)
	Sub   []*Cmd
	Usage string
}

func (c *Cmd) Parse(args []string) error {
	if len(c.Sub) == 0 {
		if err := c.Args.Parse(args); err != nil {
			return err
		}
		c.Run()
	}
	if len(c.Sub) == 0 {
		c.Run()
		return nil
	}
	for _, s := range c.Sub {
		if args[0] == s.Args.Flags.Name() {
			return s.Parse(args[1:])
		}
	}
	c.usage()
	return fmt.Errorf("cli: invalid command: %s", args[0])
}

func (c *Cmd) PrintCommands() {

}

func (c *Cmd) defaultUsage() {

}

func (c *Cmd) usage() {

}
