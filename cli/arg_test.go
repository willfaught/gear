package cli

import (
	"fmt"
	"testing"
)

func TestPrintDefaults(t *testing.T) {
	var l ArgList
	l.Flags.Init("terraform", 0)
	// l.String("x", "aa", "test arg")
	// l.String("test", "bb", "test arg")
	l.String("about", "cc", "test arg")
	l.Flags.String("g", "gg", "test flag")
	l.Flags.String("h", "hh", "test flag")
	l.Flags.String("i", "ii", "test flag")
	fmt.Println(l.Parse([]string{"-g", "ggg", "-h", "hhh", "-i", "iii", "aaa", "bbb"}))
}
