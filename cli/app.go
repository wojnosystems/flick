package cli

import "os"
import envParser "github.com/wojnosystems/go-env/v2"

func Run(cmd Commander) (err error) {
	return cmd.Switch(os.Args, &envParser.OsEnv{})
}
