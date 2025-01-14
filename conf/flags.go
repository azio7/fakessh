package conf

import (
	"flag"
	"fmt"
	"os"
)

// var _args *ArgsStruct

// FlagArgsStruct : a struct of args
type FlagArgsStruct struct {
	Help       bool
	AppVersion bool

	// Log
	LogFile   string
	LogLevel  string
	LogFormat string

	// Key
	KeyFiles   []string
	GenKeyFile bool
	KeyType    string

	// Serve
	ServPort   string
	SSHVersion string

	// Wait time
	Delay     int
	Deviation int

	// Log password
	IsLogPasswd bool

	// Anti honeypot scan
	AntiScan bool

	// Max try times
	MaxTry int

	// ConfigPath
	ConfigPath string
}

// GetArg : get args
func GetArg() (args *FlagArgsStruct, set StringSet, helper func()) {
	/*if _args != nil {
		return *_args
	}*/
	args = &FlagArgsStruct{}
	f := flag.NewFlagSet("FakeSSH", flag.ExitOnError)

	f.BoolVar(&args.Help, "h", false, "show this page")
	f.BoolVar(&args.Help, "help", false, "show this page")
	f.BoolVar(&args.AppVersion, "V", false, "show version of this binary")

	f.StringVar(&args.LogFile, "log", "", "log `file`")
	f.StringVar(&args.LogLevel, "level", DefaultLogLevel, "log level: `[debug|info|warning]`")
	f.StringVar(&args.LogFormat, "format", DefaultLogFormat, "log format: `[plain|json]`")

	var files = FlagValues{}
	f.Var(&files, "key", "key file `path`, can set more than one")
	f.BoolVar(&args.GenKeyFile, "gen", false, "generate a private key to key file path")
	f.StringVar(&args.KeyType, "type", "", "type for generate private key (default \"ed25519\")")

	f.StringVar(&args.ServPort, "bind", DefaultBind, "binding `addr`")
	f.StringVar(&args.SSHVersion, "version", DefaultSSHVersion, "ssh server version")

	f.IntVar(&args.Delay, "delay", DefaultDelay, "wait time for each login (ms)")
	f.IntVar(&args.Deviation, "devia", DefaultDeviation, "deviation for wait time (ms)")

	f.BoolVar(&args.IsLogPasswd, "passwd", false, "log password to file")

	var NoAntiScan, AntiScan bool
	f.BoolVar(&NoAntiScan, "A", false, "disable anti honeypot scan")
	f.BoolVar(&AntiScan, "a", false, "enable anti honeypot scan (default)")

	f.IntVar(&args.MaxTry, "try", DefaultMaxTry, "max try times")

	f.StringVar(&args.ConfigPath, "c", "", "config `path`")
	f.StringVar(&args.ConfigPath, "config", "", "config `path`")

	f.Parse(os.Args[1:])
	//_args = &args

	// if NoAntiScan is set and AntiScan not set, disable it
	args.AntiScan = true
	if !AntiScan && NoAntiScan {
		args.AntiScan = false
	}

	args.KeyFiles = files

	// detect used flags
	usedFlagsSet := StringSet{}
	f.Visit(func(f *flag.Flag) {
		usedFlagsSet.Add(f.Name)
	})

	return args, usedFlagsSet, f.Usage
}

// FlagValues : for multi values
type FlagValues []string

// String : implement for `flag.Value`
func (p *FlagValues) String() string {
	return fmt.Sprint(*p)
}

// Set : implement for `flag.Value`
func (p *FlagValues) Set(v string) error {
	*p = append(*p, v)
	return nil
}

func StringArrayVar(ps *[]string, name, usage string) {
	// TODO
}
