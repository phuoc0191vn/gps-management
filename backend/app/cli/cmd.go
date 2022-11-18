package cli

import (
	"fmt"
	"sort"
	"strconv"
)

type Handler interface {
	Handle()
}

type HandlerFunc func()

func (f HandlerFunc) Handle() {
	f()
}

type Commander struct {
	Name     string
	indexes  []int
	names    map[int]string
	handlers map[string]Handler
}

func NewCommander(name ...string) *Commander {
	if len(name) <= 0 || name == nil {
		name = make([]string, 0)
		name = append(name, "Commander")
	}

	return &Commander{
		Name:     name[0],
		names:    make(map[int]string),
		handlers: make(map[string]Handler),
	}
}

func (c *Commander) AddHandler(name string, handler Handler) {
	if c.names == nil {
		c.names = make(map[int]string)
	}

	if c.handlers == nil {
		c.handlers = make(map[string]Handler)
	}

	nextIndex := len(c.names) + 1

	c.names[nextIndex] = name
	c.handlers[name] = handler
	c.indexes = append(c.indexes, nextIndex)
}

func (c *Commander) Handle() {
	for {
		sort.Ints(c.indexes)

		fmt.Printf("\n%s menu: ", c.Name)
		fmt.Printf("exit(0)")
		for i, index := range c.indexes {
			if i == len(c.indexes)-1 {
				fmt.Printf("|%s(%d)\n", c.names[index], index)
				break
			}
			fmt.Printf("|%s(%d)", c.names[index], index)
		}

		var i string
		fmt.Scanln(&i)
		if i == "exit" || i == "0" {
			break
		}

		var name string
		var ok bool

		index, err := strconv.Atoi(i)
		if err != nil {
			name = i
		}

		name, ok = c.names[index]
		if !ok {
			continue
		}

		handler, ok := c.handlers[name]
		if !ok {
			continue
		}
		handler.Handle()
	}
}
