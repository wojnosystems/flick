package cli

import (
	envParser "github.com/wojnosystems/go-env/v2"
	"os"
)

func Run(cmd Commander) (err error) {
	return cmd.Switch(os.Args, &envParser.OsEnv{})
}
