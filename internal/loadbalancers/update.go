package loadbalancers

import (
	"fmt"
	"strconv"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/spf13/cobra"
)

func (c *Command) Update() *cobra.Command {
	var (
		name,
		plan,
		stickyCookie string
		stickyEnabled,
		httpsRedirectEnabled,
		proxyEnabled bool
	)
	cmd := cobra.Command{
		Use:     `update ID`,
		Args:    cobra.ExactArgs(1),
		Short:   "Update a load balancer",
		Example: `cherryctl load-balancer update 123 --name my_lb --plan load_balancer_2`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse load balancer id: %w", err)
			}

			req := cherrygo.UpdateLoadBalancer{
				Name:         name,
				Plan:         plan,
				StickyCookie: stickyCookie,
			}

			if f := cmd.Flag("sticky"); f.Changed {
				req.StickyEnabled = &stickyEnabled
			}
			if f := cmd.Flag("redirect"); f.Changed {
				req.HTTPSRedirectEnabled = &httpsRedirectEnabled
			}
			if f := cmd.Flag("proxy"); f.Changed {
				req.ProxyEnabled = &proxyEnabled
			}

			lb, _, err := c.Client().Update(ctx, id, req)
			if err != nil {
				return fmt.Errorf("failed to update load balancer %d: %w", id, err)
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

	cmd.Flags().StringVar(&name, "name", "", "Load balancer name.")
	cmd.Flags().StringVar(&plan, "plan", "", "Plan slug.")
	cmd.Flags().BoolVar(&httpsRedirectEnabled, "redirect", false, "Redirect all HTTP traffic through HTTPS with a 307 redirect.")
	cmd.Flags().BoolVar(&proxyEnabled, "proxy", false, "Preserve client IP as the request passes through the load balancer.")

	cmd.Flags().StringVar(&stickyCookie, "sticky-cookie", "", "Cookie used by sticky sessions.")
	cmd.Flags().BoolVar(&stickyEnabled, "sticky", false, "Enable sticky sessions to route subsequent requests from the same client to the same server.")

	return &cmd
}
