package loadbalancers

import (
	"fmt"
	"strconv"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/spf13/cobra"
)

func (c *Command) Get() *cobra.Command {
	return &cobra.Command{
		Use:     `get ID`,
		Args:    cobra.ExactArgs(1),
		Short:   "Retrieves load balancer details.",
		Example: `cherryctl load-balancer get 123`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: investigate silence usage
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse load balancer id: %w", err)
			}

			lb, _, err := c.Client().Get(ctx, id, c.GetOpts())
			if err != nil {
				return fmt.Errorf("failed to get load balancer %d: %w", id, err)
			}

			pubIP, _ := cherrygo.IPOfType(lb.IPAdresses, cherrygo.PrimaryIP)
			privIP, _ := cherrygo.IPOfType(lb.IPAdresses, cherrygo.PrivateIP)

			header := []string{
				"ID",
				"Plan",
				"Name",
				"Status",
				"Public IP",
				"Private IP",
				"Region"}
			data := make([][]string, 1)
			data[0] = []string{
				strconv.Itoa(lb.ID),
				lb.Plan.Slug,
				lb.Name,
				lb.Status,
				pubIP.Address,
				privIP.Address,
				lb.Region.Slug,
			}
			return c.Outputer().Output(lb, header, &data)
		},
	}
}
