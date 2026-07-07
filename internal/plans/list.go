package plans

import (
	"github.com/cherryservers/cherrygo/v3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func (c *Command) list() *cobra.Command {
	var region string
	var teamID int
	var types []string
	planGetCmd := &cobra.Command{
		Use:     `list -t <team_id> [--region <region_slug>] [--type <type>]`,
		Short:   "Retrieves a list of server plans.",
		Long:    "Retrieves a list of server plans with their corresponding hourly rates and stock volumes.",
		Example: `  # List available plans:
  cherryctl plans list`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			options := c.GetOpts()
			if options == nil {
				options = &cherrygo.GetOptions{}
			}

			if len(types) > 0 {
				options.Type = types
			}

			if region != "" {
				options.QueryParams = map[string]string{"region": region}
			}

			plans, _, err := c.Client().List(teamID, options)
			if err != nil {
				return errors.Wrap(err, "Could not list plans")
			}

			th, td := plansToTable(region, plans...)

			return c.Outputer().Output(plans, th, &td)
		},
	}

	planGetCmd.Flags().StringVarP(&region, "region", "r", "", "The Slug or ID of region.")
	planGetCmd.Flags().StringSliceVarP(&types, "type", "", []string{}, "Comma separated list of available plan types.")
	planGetCmd.Flags().IntVarP(&teamID, "team-id", "t", 0, "The team's ID. Return plans prices based on team billing details.")

	planGetCmd.MarkFlagRequired("team-id")

	return planGetCmd
}
