package cli

import (
	"flag"
	"time"
)

type Flag struct {
	Default string
	Name    string
	Usage   string
	Value   Value
}

type Value interface {
	Get() interface{}
	Set(string) error
	String() string
}

type flagValue struct {
	value flag.Value
}

func (fv flagValue) Get() interface{} {
	if g, ok := fv.value.(flag.Getter); ok {
		return g.Get()
	}
	return nil
}

func (fv flagValue) Set(s string) error {
	return fv.value.Set(s)
}

func (fv flagValue) String() string {
	return fv.value.String()
}

type FlagSet struct {
	flagSet flag.FlagSet
}

func (fs *FlagSet) Bool(name string, value bool, usage string) *bool {
	return fs.flagSet.Bool(name, value, usage)
}

func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	fs.flagSet.BoolVar(p, name, value, usage)
}

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	fs.flagSet.Duration(name, value, usage)
}

func (fs *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	fs.flagSet.DurationVar(p, name, value, usage)
}

func (fs *FlagSet) Flags() int {
	return fs.flagSet.NFlag()
}

func (fs *FlagSet) Float64(name string, value float64, usage string) *float64 {
	fs.flagSet.Float64(name, value, usage)
}

func (fs *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	fs.flagSet.Float64Var(p, name, value, usage)
}

func (fs *FlagSet) Int(name string, value int, usage string) *int {
	fs.flagSet.Int(name, value, usage)
}

func (fs *FlagSet) Int64(name string, value int64, usage string) *int64 {
	fs.flagSet.Int64(name, value, usage)
}

func (fs *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	fs.flagSet.Int64Var()
}

func (fs *FlagSet) IntVar(p *int, name string, value int, usage string) {
	fs.flagSet.IntVar()
}

func (fs *FlagSet) Lookup(name string) *Flag {
	var l = fs.flagSet.Lookup(name)
	if l == nil {
		return nil
	}
	return &Flag{
		Default: l.DefValue,
		Name:    l.Name,
		Usage:   l.Usage,
		Value:   flagValue{l.Value},
	}
}

func (fs *FlagSet) String(name, value, usage string) *string {
	return fs.flagSet.String(name, value, usage)
}

func (fs *FlagSet) StringVar(p *string, name, value, usage string) {
	fs.flagSet.StringVar(p, name, value, usage)
}

func (fs *FlagSet) Uint(name string, value uint, usage string) *uint {
	fs.flagSet.Uint(name, value, usage)
}

func (fs *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	fs.flagSet.Uint64(name, value, usage)
}

func (fs *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	fs.flagSet.Uint64Var(p, name, value, usage)
}

func (fs *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	fs.flagSet.UintVar(p, name, value, usage)
}

func (fs *FlagSet) Var(value FlagValue, name string, usage string) {
	fs.flagSet.Var(value, name, usage)
}

func (fs *FlagSet) Visit(f func(*Flag)) {
	fs.flagSet.Visit(func(g *flag.Flag) {
		f(&Flag{
			Default: g.DefValue,
			Name:    g.Name,
			Usage:   g.Usage,
			Value:   flagValue{g.Value},
		})
	})
}

func (fs *FlagSet) VisitAll(f func(*Flag)) {
	fs.flagSet.VisitAll(func(g *flag.Flag) {
		f(&Flag{
			Default: g.DefValue,
			Name:    g.Name,
			Usage:   g.Usage,
			Value:   flagValue{g.Value},
		})
	})
}
