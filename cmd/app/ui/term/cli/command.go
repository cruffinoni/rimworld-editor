package cli

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Command struct {
	Name        string
	Usage       string
	Description string
	params      []*Params
	alias       []string
	childCmds   []*Command
	Handler     any
	logger      logging.Logger
}

type CmdSpan func(*Command)

func (c *Command) init() {
	c.alias = strings.Split(strings.TrimSpace(c.Name), "|")
	for i := range c.params {
		n := strings.Split(strings.TrimSpace(c.params[i].Name), "|")
		if len(n) > 1 {
			c.params[i].alias = n[1:]
		}
		c.params[i].Name = n[0]
	}
}

func (c *Command) NewCommand(cmd *Command, f CmdSpan, params ...Params) {
	c.params = make([]*Params, 0)
	duplicate := make(map[string]bool)
	for _, p := range params {
		if _, ok := duplicate[strings.ToLower(p.Name)]; ok {
			log.Fatalf("duplicate param name %s", p.Name)
		}
		duplicate[strings.ToLower(p.Name)] = true
		c.params = append(c.params, &p)
	}
	c.childCmds = make([]*Command, 0)
	c.childCmds = append(c.childCmds, cmd)
	cmd.logger = c.logger
	cmd.init()
	if f != nil {
		f(cmd)
	}
}

func (c *Command) WithParameter(p *Params) *Command {
	c.params = append(c.params, p)
	return c
}

func (c *Command) getCommandNames(withCommandName bool) []string {
	names := make([]string, 0, len(c.alias)+1)
	for _, cmd := range c.alias {
		names = append(names, cmd)
	}
	if withCommandName {
		names = append(names, c.Name)
	}
	return names
}

func (c *Command) getChildCommandNames() []string {
	names := make([]string, 0, len(c.childCmds))
	for _, cmd := range c.childCmds {
		names = append(names, cmd.Name)
	}
	return names
}

var (
	ErrCommandNotFound   = errors.New("command not found")
	ErrEmptyInput        = errors.New("empty input")
	ErrNoHandler         = errors.New("no handler found")
	ErrMustBeAFunction   = errors.New("handler must be a function")
	ErrParamDoesntExists = errors.New("a parameter doesn't exists")
)

func (c *Command) PrintUsage() {
	c.logger.WithFields(logging.Fields{
		"command": c.Name,
		"usage":   c.Usage,
	}).Info("Command usage")
	//printer.Debugf("%s [%q]", c.Name, strings.Join(c.getChildCommandNames(), "|"))
}

func (c *Command) PrintHelp() {
	c.logger.WithFields(logging.Fields{
		"command": c.Name,
		"help":    c.Description,
		"children": strings.Join(c.getChildCommandNames(), "|"),
	}).Info("Command help")
}

func (c *Command) findParam(name string) *Params {
	for _, p := range c.params {
		if p.Name == name {
			return p
		}
	}
	return nil
}

/*
pawn injury heal PAWN_3 --wounds
faction delete FACTION_18 -f
*/

func (c *Command) applyParameter(param *Params, args string, t reflect.Type) (*reflect.Value, error) {
	var namedParam bool
	for args[0] == '-' {
		namedParam = true
		args = args[1:]
	}
	if !param.isNamed(args) || !namedParam {
		return nil, nil
	}
	if len(args) == 0 && param.OptionalValue == nil {
		return nil, errors.New("no param left")
	}
	v := reflect.New(t)
	switch param.Type {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(args, 10, 64); err == nil {
			v.SetInt(i)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(args, 64); err == nil {
			v.SetFloat(f)
		}
	case reflect.String:
		v.SetString(args)
	}
	return &v, nil
}

func (c *Command) callHandler(args []string) error {
	if c.Handler == nil {
		return ErrNoHandler
	}
	//h := reflect.TypeOf(c.Handler)
	//if h.Kind() != reflect.Func {
	//	return ErrMustBeAFunction
	//}
	//hv := reflect.ValueOf(c.Handler)
	//in := h.NumIn()
	//val := make([]reflect.Val, 0)
	//var cpyArgs []string
	//copy(cpyArgs, args)
	//if in > 0 {
	//	if len(args) != in {
	//		printer.Debugf("Size difference: %v & %v", len(args), in)
	//	}
	//	for i := 0; i < in; i++ {
	//		arg := h.In(i)
	//		//switch arg.Kind() {
	//		//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//		//	val = append(val, reflect.ValueOf())
	//		//}
	//	}
	//}
	return nil
}

func (c *Command) RunWithArgs(args []string) error {
	if len(args) == 0 {
		return ErrEmptyInput
	}
	if strings.ToLower(args[0]) == "help" {
		c.PrintHelp()
		return nil
	}
	for _, n := range c.childCmds {
		names := n.getCommandNames(true)
		for _, name := range names {
			if name == args[0] {
				c.logger.WithField("command", name).Debug("Command found")
				if len(n.childCmds) > 0 {
					if len(args) == 1 {
						n.PrintUsage()
						return nil
					}
					return n.RunWithArgs(args[1:])
				}
				//n.Handler(args[1:])
				return nil
			}
		}
	}
	return ErrCommandNotFound
}
