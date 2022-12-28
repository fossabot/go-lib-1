package commands

import (
	"github.com/kc-workspace/go-lib/commandline/flags"
	"github.com/kc-workspace/go-lib/mapper"
)

type Command struct {
	Name     string
	Flags    *flags.Manager
	Executor Executor
}

func (c *Command) Start(p *ExecutorParameter) error {
	if c.Flags != nil {
		var option, args, err = c.Flags.Build(c.Name, p.Args)
		if err != nil {
			return err
		}

		for _, value := range option {
			p.Config = mapper.Merger(p.Config).Add(value).Merge()
		}
		p.Args = args
	}

	return c.Executor(p)
}
