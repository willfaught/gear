package cli

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Arg struct {
	Default string
	Name    string
	Usage   string
	Value   ArgValue
}

type ArgList struct {
	Flags flag.FlagSet
	Usage func()
	args  map[string]*Arg
	order []*Arg
}

func (l *ArgList) handleError(err error) error {
	switch l.Flags.ErrorHandling() {
	case flag.ExitOnError:
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
	return err
}

func (l *ArgList) Lookup(name string) *Arg {
	return l.args[name]
}

func (l *ArgList) NArgs() int {
	return len(l.args)
}

func (l *ArgList) Parse(args []string) error {
	l.Flags.Usage = l.usage
	if err := l.Flags.Parse(args); err != nil {
		return l.handleError(err)
	}
	args = l.Flags.Args()
	for _, a := range l.order {
		if len(args) == 0 {
			l.usage()
			return l.handleError(fmt.Errorf("cli: missing value for argument %s", a.Name))
		}
		if err := a.Value.Set(args[0]); err != nil {
			l.usage()
			return l.handleError(fmt.Errorf("cli: cannot set value for argument %s: %v", a.Name, err))
		}
		args = args[1:]
	}
	if len(args) > 0 {
		l.usage()
		return l.handleError(fmt.Errorf("cli: unused arguments: %v", strings.Join(args, " ")))
	}
	return nil
}

func (l *ArgList) Parsed() bool {
	return l.Flags.Parsed()
}

func unquoteUsage(a *Arg) (name string, usage string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = a.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}
	// No explicit name, so use type if we can find one.
	name = "value"
	switch a.Value.(type) {
	case *float64ArgValue:
		name = "float"
	case *intArgValue:
		name = "int"
	case *stringArgValue:
		name = "string"
	}
	return
}

func (l *ArgList) PrintDefaults() {
	var output = l.Flags.Output()
	for _, a := range l.order {
		var s = fmt.Sprintf("  %s", a.Name)
		var name, usage = unquoteUsage(a)
		if len(name) > 0 {
			s += " " + name
		}
		if len(s) <= 4 {
			s += "\t"
		} else {
			s += "\n    \t"
		}
		s += strings.Replace(usage, "\n", "\n    \t", -1)
		// if !isZeroValue(a, a.Default) {
		// 	if _, ok := flag.Value.(*stringValue); ok {
		// 		// put quotes on the value
		// 		s += fmt.Sprintf(" (default %q)", a.Default)
		// 	} else {
		// 		s += fmt.Sprintf(" (default %v)", a.Default)
		// 	}
		// }
		fmt.Fprint(output, s, "\n")
	}
}

func (l *ArgList) Set(name, value string) error {
	var a, ok = l.args[name]
	if !ok {
		return l.handleError(fmt.Errorf("cli: argument name %s is invalid", name))
	}
	return a.Value.Set(value)
}

func (l *ArgList) String(name, value, usage string) *string {
	var s string
	l.StringVar(&s, name, value, usage)
	return &s
}

func (l *ArgList) StringVar(p *string, name, value, usage string) {
	*p = value
	l.Var((*stringArgValue)(p), name, usage)
}

func (l *ArgList) Var(value ArgValue, name string, usage string) {
	if _, ok := l.args[name]; ok {
		panic(fmt.Sprintf("cli: duplicate argument name %s", name))
	}
	var arg = &Arg{Default: value.String(), Name: name, Usage: usage, Value: value}
	if l.args == nil {
		l.args = map[string]*Arg{}
	}
	l.args[name] = arg
	l.order = append(l.order, arg)
}

func (l *ArgList) defaultUsage() {
	var b = bytes.NewBufferString("Usage: ")
	if n := l.Flags.Name(); n != "" {
		b.WriteString(n + " ")
	}
	if l.Flags.NFlag() > 0 {
		b.WriteString("[<options>] ")
	}
	for i, a := range l.order {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(a.Name)
	}
	var output = l.Flags.Output()
	fmt.Fprintln(output, b.String())
	if l.Flags.NFlag() > 0 {
		fmt.Fprint(output, "\nOptions:\n")
		l.Flags.PrintDefaults()
	}
	if len(l.args) > 0 {
		fmt.Fprint(output, "\nArguments:\n")
		l.PrintDefaults()
	}
}

func (l *ArgList) usage() {
	if l.Usage == nil {
		l.defaultUsage()
	} else {
		l.Usage()
	}
}

type ArgValue interface {
	Get() interface{}
	String() string
	Set(string) error
}

type float64ArgValue float64

func (v float64ArgValue) Get() interface{} {
	return float64(v)
}

func (v *float64ArgValue) Set(s string) error {
	var f, err = strconv.ParseFloat(s, 64)
	if err == nil {
		*v = float64ArgValue(f)
	}
	return err
}

func (v float64ArgValue) String() string {
	return fmt.Sprint(float64(v))
}

type intArgValue int

func (v intArgValue) Get() interface{} {
	return int(v)
}

func (v *intArgValue) Set(s string) error {
	var i, err = strconv.Atoi(s)
	if err == nil {
		*v = intArgValue(i)
	}
	return err
}

func (v intArgValue) String() string {
	return strconv.Itoa(int(v))
}

type stringArgValue string

func (v stringArgValue) Get() interface{} {
	return string(v)
}

func (v *stringArgValue) Set(s string) error {
	*v = stringArgValue(s)
	return nil
}

func (v stringArgValue) String() string {
	return string(v)
}
