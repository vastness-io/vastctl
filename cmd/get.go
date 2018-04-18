package cmd

import "github.com/urfave/cli"

type Resource interface {
	Get() cli.Command
	ShortName() string
	Singular() string
	Plural() string
}

type resource struct {
	shortName string
	singular  string
	plural    string
	action    interface{}
	flags     []cli.Flag
}

func (r *resource) ShortName() string {
	return r.shortName
}

func (r *resource) Singular() string {
	return r.singular
}

func (r *resource) Plural() string {
	return r.plural
}

func (r *resource) Get() cli.Command {
	return cli.Command{
		Name:      r.Plural(),
		ShortName: r.ShortName(),
		Aliases: []string{
			r.Singular(),
		},
		Action: r.action,
		Flags:  r.flags,
	}
}

func (r *resource) Action() interface{} {
	return r.action
}

func GetResourceCommands() cli.Commands {
	return []cli.Command{
		projectResource.Get(),
	}
}
