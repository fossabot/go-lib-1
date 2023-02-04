package plugins

import (
	"os"
	"path"

	"github.com/kc-workspace/go-lib/commandline/flags"
	"github.com/kc-workspace/go-lib/commandline/hooks"
	"github.com/kc-workspace/go-lib/dotenv"
	"github.com/kc-workspace/go-lib/fs"
	"github.com/kc-workspace/go-lib/mapper"
)

// SupportDotEnv will create --envs option for custom load .env files
//
// For parameters:
//
// If autoload is true, it will try to load $PWD/.env and log warn if not exist;
// otherwise, it will not try to load any thing by default.
func SupportDotEnv(autoload bool) Plugin {
	return func(p *PluginParameter) error {
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

		var defaultEnv = make([]string, 0)
		if autoload {
			defaultEnv = append(defaultEnv, path.Join(wd, ".env"))
		}

		p.NewFlags(flags.Array{
			Name:    "envs",
			Default: defaultEnv,
			Usage:   "environment file/directory. each file must following .env regulation",
			Action: func(data []string) mapper.Mapper {
				if len(data) > 0 {
					return mapper.New().
						Set("fs.env.type", "auto").
						Set("fs.env.mode", "multiple").
						Set("fs.env.fullpath", data)
				}

				return mapper.New()
			},
		})

		p.NewFlags(flags.Bool{
			Name:    "no-env-file",
			Default: false,
			Usage:   "disabled loading .env files completely",
			Action: func(data bool) mapper.Mapper {
				return mapper.New().
					Set("internal.flag.noenv", data)
			},
		})

		p.NewHook(hooks.AFTER_FLAG, func(m mapper.Mapper) error {
			if disabled := m.Mi("internal").Mi("flag").Bo("noenv", false); disabled {
				return nil
			}

			if m.Has("fs.env") {
				envs, err := fs.Build(m.Mi("fs").Mi("env"), m.Mi("variables"))
				if err != nil {
					p.Logger.Warn("cannot found environment file: %v, skipped", err)
					return nil
				}

				// write environment value from .env file
				err = dotenv.Overload(envs.Multiple()...)
				if err != nil {
					p.Logger.Warn("dotenv return error: %v", err)
				}
			}

			return nil
		})

		return nil
	}
}
