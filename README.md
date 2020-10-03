# Overview

I hate Vipre/Cobra and go's flags. I want to be able to create a way to load configuration and validate it as much at compile time as possible.

## Principles

* Configuration should, eventually, be interpreted in code as a structure, without having to be cast.
* Configurations are application-specific, there's no real need to genericify it within the application
* Configurations are reloadable. Don't use static anything. Parsing creates new configuration sets in memory
* Configurations are usually not needed to be loaded/interpeted with performance as the primary concern, opt for ease of use and readability over performance
* Configurations need validation from the application.
* Users of the application need to know which setting caused a validation failure and where to fix it (environment variable? file? flag? where in the file? where in the command flags?)

## Configuration life cycle

1. Definition phase
   1. individual values
   1. individual value validations
   1. post-load transforms
   1. grouped values
   1. grouped value validations
1. Create configs!
1. Test configs (outside of the application), make it easy to build config linters for formatting and validation.
1. Load
   1. From file
   1. From environment
   1. From flags
1. Validate
   1. individual values
   1. groups of values
   1. groups
1. Return a configuration set

## What do you mean testing/linting?

Treat configs like you would a REST request. It's user-data. If you wouldn't put it into a database without validation it, you probably don't want to distribute an application with it.

# How I want it to work

I'd like this scenario:

```go
package main
import (

"log"
"os"
)

type configConfig struct {
  config.Default
  Path value.String
}

func DefaultConfigConfig() configConfig {
  return configConfig{
    Path: os.Getenv("HOME") + "/.myapp/conf",
  }
}

type appConfig struct {
  config.Default
  Name value.String
  Age value.Int
  Databases []dbConfig
}

func NewAppConfig() appConfig {
  return appConfig{
    Name: value.String{
      Default: "puppy",
      Validations: []validation.Definition{
        validation.
      },
    },
  }
}

type dbConfig struct {
  config.Default
  Host value.String
  Username value.String
  Password value.String
}

func main() {
   configLocation := DefaultConfigConfig() 
   _ = config.LoadAndValidate(&configLocation, validation.Skip, source.SourceEnv() )
   appCfg := appConfig{}
   validationErrors := validation.Collection{}
   err := config.Load(
      &appCfg,
      &validationErrors,
      source.YamlFile(configLocation.Path),
      source.Env(),
      source.SourceFlags )
  if err != nil {
    log.Fatal("unable to parse the config file at", configLocation.Path)
  }
  if validationErrors.HasAny() {
     validation.Println(&validationErrors)
     log.Fatal("Fix the above issues and try again!")
  }
}
```
