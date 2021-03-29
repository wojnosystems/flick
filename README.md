# Overview

I hate Vipre/Cobra and go's flags. I want to be able to create a way to load configuration and validate it as much at compile time as possible. Not only that, but configurations should be reloadable. New configurations should emerge as new objects in memory so that graceful switch overs can occur.

## Principles

* Configuration should, eventually, be interpreted in code as a structure, without having to be cast.
* Configurations are application-specific, there's no real need to genericify it within the application
* Configurations are reloadable. Don't use static anything. Parsing creates new configuration sets in memory
* Configurations are usually not needed to be loaded/interpeted with performance as the primary concern, opt for ease of use and readability over performance
* Configurations need validation from the application.
* Users of the application need to know which setting caused a validation failure and where to fix it (environment variable? file? flag? where in the file? where in the command flags?)

## Configuration life cycle

1. Definition phase (this happens at compile-time)
   1. individual values
   1. individual value validations
   1. post-load transforms
   1. grouped values
   1. grouped value validations
1. Test configs (outside of the application), make it easy to build config linters for formatting and validation.
1. Load (this happens at run-time)
   1. From file
   1. From environment
   1. From flags
1. Validate
   1. individual values
   1. groups of values
   1. groups
1. Return a configuration set

## What do you mean testing/linting?

Treat configs like you would a REST request. It's user-data. If you wouldn't put it into a database without validating it, you probably don't want to distribute an application with it. You also don't want an application to start with bad configs. We'll be using validations to enforce this. These run only after all sources have been loaded.

# How I want it to work

I'd like this scenario:

```go
package main
import (
  "fmt"
  "github.com/wojnosystems/flick/action"
  flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
  "github.com/wojnosystems/flick/cli"
  "github.com/wojnosystems/flick/parse"
  "github.com/wojnosystems/go-env/v2"
  "github.com/wojnosystems/go-optional/v2"
  "github.com/wojnosystems/okey-dokey/bad"
  "github.com/wojnosystems/okey-dokey/ok_range"
  "github.com/wojnosystems/okey-dokey/ok_string"
  "log"
  "os"
  "strings"
  "time"
)

type appConfig struct {
  Profile        optional.String   `yaml:"profile" flag:"profile" flag-short:"p" env:"PROFILE"`
  ConnectTimeout optional.Duration `yaml:"connectTimeout" flag:"connectTimeout" env:"CONNECT_TIMEOUT" usage:"Ns" help:"duration to wait for connections to complete before failing"`
  Renamed        optional.String   `yaml:"renamed"`
}
func (c *appConfig)Validate(emitter bad.Emitter) (err error) {
  ok_string.Validate(c.Profile, &ok_string.On{
    Ensure: []ok_string.Definer{
      &ok_string.IsRequired{},
      &ok_string.LengthAtLeast{
        Format: func(definition *ok_string.LengthAtLeast, value optional.String) string {
          return "was too short, buddy!"
        },
        Length: 3,
      },
    },
  }, emitter.Into("profile"))
  return
}
func main() {
  cfg := appConfig{
    ConnectTimeout: optional.DurationFrom(30*time.Second),
  }
  // Load the global configuration
  // loads the file under ~/.myapp/config.yaml, optionally the CONFIG_FILE_PATH env var, or optionally overridden with --config-file-path= flag
  // contents are stored in the global variable: "cfg"
  flagGroups, err := parse.UnmarshalWithFile( os.Getenv("HOME") + "/.myapp/config.yaml", parse.Yaml(), os.Args[1:], &cfg )
  if err != nil {
    log.Fatal("unable to parse the config file", err)
  }

  err = cli.Run(flagGroups)
  if err != nil {
    log.Panic(err)
  }
}

type connectFlags struct {
  Host optional.String `flag:"host" flag-short:"h" usage:"HOST" help:"connect to this host"`
}
func (f *connectFlags)Validate(emitter bad.MemberEmitter) (err error) {
  ok_string.Validate(f.Host, &ok_string.On{
    Ensure: []ok_string.Definer{
      &ok_string.IsRequired{},
      &ok_string.LengthBetween{
        Between: ok_range.IntBetween(3, 32),
      },
    },
  }, emitter.Into("host"))
  return
}
type connect struct {
  Flags          connectFlags
  PositionalArgs []string
}
func (c* connect)Invoke() (err error) {
  commands := strings.Builder{}
  for _, pa := range c.PositionalArgs {
    commands.WriteString(pa)
  }
  fmt.Println("I'm connecting to host", c.Flags.Host.Value(), "with profile", cfg.Profile.Value(), "with commands:", commands.String())
  return
}
```

If you call myapp with:

```shell script
$ myapp version
> profile is required
$ myapp --profile=ch version
> profile was too short, buddy!
$ myapp --profile=chris version
> v1.0.0
$ echo "rename: test" > /home/wojno/.myapp/config.yaml
$ CONNECT_TIMEOUT=45s myapp --flagTrace --profile=chris
ConfigFile:
  ConfigFilePath: /home/wojno/.myapp/config.yaml (default)
appConfig:
  Profile: chris (flag:--profile)
  ConnectTimeout: 45s (env:CONNECT_TIMEOUT)
  Renamed: test (file:config.yaml)
error no command specified
Usage: myapp [--configFilePath=PATH] [--profile=STRING] [--connectTimeout=Ns] COMMAND
Available Flags:
  --configFilePath, -c: file path to the configuration file that will set the global defaults
  --profile, -p
  --connectTimeout: duration to wait for connections to complete before failing
Available Commands:
  connect
  help
  version

$ myapp connect help
error host is required
Usage: myapp connect [--host=HOST]
Available Flags:
  --host, -h: connect to this host
Available Commands:
  help
```

## optionapi

```yaml
optionapi:
   version: 1
commands:
   server:
      options:
         - $ref: "#/components/options/ConnectTimeout"
      commands:
         start:
            usage: "how to use"
            description: "Long text explaining things"
            options:
               - $ref: "#/components/options/HasBanana"
         stop:
            
         restart:
            
components:
   options:
      ConnectTimeout:
         type: duration
         description: "how long to wait when connecting to the server"
         usage: "Ns"
         env:
            name: "CONNECT_TIMEOUT"
         flag:
            name: "connectTimeout"
            aliases: ["c"]
         default: 30s
      HasBanana:
         type: bool
         description: "specify to use bananas"
         env:
            name: "BANANA"
         flag:
            name: "banana"
            aliases: ["b"]
         default: false
```

Produces the following GoLang file:

```go
package flickstub

import (
   "context"
   "errors"
   "github.com/wojnosystems/go-optional/v2"
)

type Interface interface {
   HookBefore(ctx context.Context, opts *AllCommandsOptions) error
   ServerHookBefore(ctx context.Context, opts *ServerOptions) error
   ServerStart(ctx context.Context, opts *ServerStartOptions) error
   ServerStop(ctx context.Context, opts *ServerStopOptions) error
   ServerRestart(ctx context.Context, opts *ServerRestartOptions) error
   ServerHookAfter(ctx context.Context, opts *ServerOptions, err error) error
   HookAfter(ctx context.Context, opts *AllCommandsOptions, err error) error
}

type AllCommandsOptions struct {
	
}

type ServerOptions struct {
	AllCommands AllCommandsOptions
   ConnectTimeout optional.Duration
}

type ServerStartOptions struct {
   Server    ServerOptions
   HasBanana optional.Bool
}

type ServerStopOptions struct {
   Server ServerOptions
}

type ServerRestartOptions struct {
   Server ServerOptions
}

var ErrUnimplemented = errors.New("command has not been implemented")

type Unimplemented struct {
}

func (u *Unimplemented) HookBefore(ctx context.Context, opts *AllCommandsOptions) error {
   return nil
}
func (u *Unimplemented) ServerHookBefore(ctx context.Context, opts *ServerOptions) error {
   return nil
}
func (u *Unimplemented) ServerStart(ctx context.Context, opts *ServerStartOptions) error {
   return ErrUnimplemented
}
func (u *Unimplemented) ServerStop(ctx context.Context, opts *ServerStopOptions) error {
   return ErrUnimplemented
}
func (u *Unimplemented) ServerRestart(ctx context.Context, opts *ServerRestartOptions) error {
   return ErrUnimplemented
}
func (u *Unimplemented) ServerHookAfter(ctx context.Context, opts *ServerOptions, err error) error {
   return nil
}
func (u *Unimplemented) HookAfter(ctx context.Context, opts *AllCommandsOptions, err error) error {
   return nil
}

```


# Testing

You'll want to run some checks on your command line interface definition. You can do this easily within a main_test.go file:

```go
package main

func TestCli(t *testing.T) {
	flick.TestRun(t, &app)
}
```

Primarily, this will ensure that all of your functions for your actions match up with your data types. It will also run other tests as deemed useful later :D.
