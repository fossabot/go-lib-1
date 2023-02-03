package plugins

import (
	"os"

	"github.com/kc-workspace/go-lib/commandline/flags"
	"github.com/kc-workspace/go-lib/commandline/hooks"
	"github.com/kc-workspace/go-lib/configs"
	"github.com/kc-workspace/go-lib/mapper"
)

// SupportConfig will load configuration from configs directory
// It also add --pwd for current directory, --configs for custom config directory
func SupportConfig(p *PluginParameter) error {
	var wd, err = os.Getwd()
	if err != nil {
		return err
	}

	p.NewFlags(flags.String{
		Name:    "pwd",
		Default: wd,
		Usage:   "current directory",
		Action: func(data string) mapper.Mapper {
			return mapper.New().Set("variables.current", data)
		},
	})

	p.NewFlags(flags.Array{
		Name:    "configs",
		Default: []string{},
		Usage:   "configuration file/directory. directory must contains only json files and file must be json",
		Action: func(data []string) mapper.Mapper {
			var result = mapper.New()
			if len(data) > 0 {
				result.
					Set("fs.config.type", "auto").
					Set("fs.config.mode", "multiple").
					Set("fs.config.fullpath", data)
			}
			return result
		},
	})

	p.NewHook(hooks.BEFORE_COMMAND, func(config mapper.Mapper) error {
		var name = p.Metadata.Name
		if name == "" {
			name = DEFAULT_ENV_PREFIX
		}

		var addition, err = configs.New(name, config).Build(os.Environ())
		if err != nil {
			return err
		}

		addition.ForEach(func(key string, value interface{}) {
			config.Set(key, value)
		})
		return nil
	})

	return nil
}
