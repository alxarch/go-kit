package kit

import (
	"flag"
	"os"
	"strings"
)

func envName(n string, prefix string) string {
	n = strings.Replace(n, "-", "_", -1)
	n = strings.Replace(n, ".", "_", -1)
	n = strings.ToUpper(n)
	if prefix != "" {
		prefix = envName(prefix, "")
		n = prefix + "_" + n
	}
	return n
}

// ParseEnv parses flags for a FlagSet reading env vars
func ParseEnv(fs *flag.FlagSet, args []string) error {
	return ParseEnvPrefix(fs, args, "")
}

// ParseEnvPrefix parses flags for a FlagSet reading env vars prefixed with prefix
func ParseEnvPrefix(fs *flag.FlagSet, args []string, prefix string) error {
	if err := fs.Parse(args); err != nil {
		return err
	}
	unset := make(map[string]flag.Value)
	fs.VisitAll(func(f *flag.Flag) {
		unset[f.Name] = f.Value
	})
	fs.Visit(func(f *flag.Flag) {
		delete(unset, f.Name)
	})
	for key, v := range unset {
		key = envName(key, prefix)
		if val, ok := os.LookupEnv(key); ok {
			v.Set(val)
		}
	}
	return nil
}
