package loadbalancers

import (
	"github.com/cherryservers/cherryctl/internal/outputs"
	"github.com/cherryservers/cherrygo/v4"
	"github.com/spf13/cobra"
)

type Command struct {
	Deps
}

func (c *Command) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     `load-balancer`,
		Aliases: []string{"loadbalancer", "loadbalancers", "load-balancers", "lb", "lbs"},
		Short:   "Load balancer operations. For details on load balancers, see the product docs: https://www.cherryservers.com/knowledge/docs/networking/load-balancer",
	}

	cmd.AddCommand(
		c.Get(),
	)

	return cmd
}

type Deps interface {
	Client() cherrygo.LoadBalancerService
	GetOpts() *cherrygo.GetOptions
	Outputer() outputs.Outputer
}

func NewCommand(dep Deps) *Command {
	return &Command{
		Deps: dep,
	}
}
