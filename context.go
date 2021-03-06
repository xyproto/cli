package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type (
	Args []string

	// Context is a type that is passed through to each Handler action in a cli application.
	// Context can be used to retrieve context-specific Args and parsed command-line options.
	Context struct {
		App       *App
		Command   Command
		flagSet   *flag.FlagSet
		globalSet *flag.FlagSet
		setFlags  map[string]bool
	}
)

// NewContext creates a new context. For use in when invoking an App or Command action.
func NewContext(app *App, set *flag.FlagSet, globalSet *flag.FlagSet) *Context {
	return &Context{App: app, flagSet: set, globalSet: globalSet}
}

// Int looks up the value of a local int flag, returns 0 if no int flag exists.
func (c *Context) Int(name string) int {
	return lookupInt(name, c.flagSet)
}

// Float64 looks up the value of a local float64 flag, returns 0 if no float64 flag exists.
func (c *Context) Float64(name string) float64 {
	return lookupFloat64(name, c.flagSet)
}

// Bool looks up the value of a local bool flag, returns false if no bool flag exists.
func (c *Context) Bool(name string) bool {
	return lookupBool(name, c.flagSet)
}

// BoolT looks up the value of a local boolT flag, returns false if no bool flag exists.
func (c *Context) BoolT(name string) bool {
	return lookupBoolT(name, c.flagSet)
}

// String looks up the value of a local string flag, returns "" if no string flag exists.
func (c *Context) String(name string) string {
	return lookupString(name, c.flagSet)
}

// StringSlice looks up the value of a local string slice flag, returns nil if no string slice flag exists.
func (c *Context) StringSlice(name string) []string {
	return lookupStringSlice(name, c.flagSet)
}

// IntSlice looks up the value of a local int slice flag, returns nil if no int slice flag exists.
func (c *Context) IntSlice(name string) []int {
	return lookupIntSlice(name, c.flagSet)
}

// GlobalInt looks up the value of a global int flag, returns 0 if no int flag exists
func (c *Context) GlobalInt(name string) int {
	return lookupInt(name, c.globalSet)
}

// GlobalBool looks up the value of a global bool flag, returns false if no bool flag exists.
func (c *Context) GlobalBool(name string) bool {
	return lookupBool(name, c.globalSet)
}

// GlobalString looks up the value of a global string flag, returns "" if no string flag exists.
func (c *Context) GlobalString(name string) string {
	return lookupString(name, c.globalSet)
}

// GlobalStringSlice looks up the value of a global string slice flag, returns nil if no string slice flag exists.
func (c *Context) GlobalStringSlice(name string) []string {
	return lookupStringSlice(name, c.globalSet)
}

// GlobalIntSlice looks up the value of a global int slice flag, returns nil if no int slice flag exists.
func (c *Context) GlobalIntSlice(name string) []int {
	return lookupIntSlice(name, c.globalSet)
}

// IsSet determines if the flag was actually set exists.
func (c *Context) IsSet(name string) bool {
	if c.setFlags == nil {
		c.setFlags = make(map[string]bool)
		c.flagSet.Visit(func(f *flag.Flag) {
			c.setFlags[f.Name] = true
		})
	}
	return c.setFlags[name] == true
}

// Args returns the command line arguments associated with the context.
func (c *Context) Args() Args {
	args := Args(c.flagSet.Args())
	return args
}

// Get returns the nth argument, or else a blank string.
func (a Args) Get(n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

// First returns the first argument, or else a blank string.
func (a Args) First() string {
	return a.Get(0)
}

// Tail return the rest of the arguments (not the first one) or else an empty string slice.
func (a Args) Tail() []string {
	if len(a) >= 2 {
		return []string(a)[1:]
	}
	return []string{}
}

// Present checks if there are any arguments present.
func (a Args) Present() bool {
	return len(a) != 0
}

// lookupInt retrieves the Int value of a named flag.
func lookupInt(name string, set *flag.FlagSet) int {
	f := set.Lookup(name)
	// bail out if name is not found in set
	if f == nil {
		return 0
	}
	// get the Int value
	val, err := strconv.Atoi(f.Value.String())
	if err != nil {
		return 0
	}
	return val
}

// lookupFloat64 retrieves the Float64 value of a named flag.
func lookupFloat64(name string, set *flag.FlagSet) float64 {
	f := set.Lookup(name)
	// bail out if name is not found in set
	if f == nil {
		return 0
	}
	// get the Float64 value
	val, err := strconv.ParseFloat(f.Value.String(), 64)
	if err != nil {
		return 0
	}
	return val
}

// lookupString retrieves the String value of a named flag.
func lookupString(name string, set *flag.FlagSet) string {
	f := set.Lookup(name)
	// bail out if name is not found in set
	if f == nil {
		return ""
	}
	// get and return the string value
	return f.Value.String()
}

// lookupStringSlice retrieves the StringSlice value of a named flag.
func lookupStringSlice(name string, set *flag.FlagSet) []string {
	f := set.Lookup(name)
	// bail out if name is not found in set
	if f == nil {
		return nil
	}
	// get and return the string slice value
	return (f.Value.(*StringSlice)).Value()
}

// lookupIntSlice retrieves the IntSlice value of a named flag.
func lookupIntSlice(name string, set *flag.FlagSet) []int {
	f := set.Lookup(name)
	// bail out if name is not found in set
	if f == nil {
		return nil
	}
	// get and return the int slice value
	return (f.Value.(*IntSlice)).Value()
}

// lookupBool retrieves the Bool value of a named flag.
func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	// bail out if name is not found in set
	if f == nil {
		return false
	}
	// get the Bool value
	val, err := strconv.ParseBool(f.Value.String())
	if err != nil {
		return false
	}
	return val
}

// lookupBoolT retrieves the BoolT value of a named flag.
func lookupBoolT(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f == nil {
		return false
	}
	// get the BoolT value
	val, err := strconv.ParseBool(f.Value.String())
	if err != nil {
		return true
	}
	return val
}

func copyFlag(name string, ff *flag.Flag, set *flag.FlagSet) {
	switch ff.Value.(type) {
	case *StringSlice:
	default:
		set.Set(name, ff.Value.String())
	}
}

// mapS applies the given function on each element in the string slice.
func mapS(sl []string, f func(string) string) []string {
	newlist := make([]string, len(sl))
	for i, element := range sl {
		newlist[i] = f(element)
	}
	return newlist
}

func normalizeFlags(flags []Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		// split flags by comma and strip the whitespace from each element
		parts := mapS(strings.Split(f.getName(), ","), strings.TrimSpace)
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			if !visited[name] {
				copyFlag(name, ff, set)
			}
		}
	}
	return nil
}
