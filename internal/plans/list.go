package plans

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func (c *Client) List() *cobra.Command {
	var regionID string
	var teamID int
	var types []string
	planGetCmd := &cobra.Command{
		Use:     `list -t <team_id> [--region-id <region_slug>] [--type <type>]`,
		Aliases: []string{"get"},
		Short:   "Retrieves a list of server plans.",
		Long:    "Retrieves a list of server plans with their corresponding hourly rates and stock volumes.",
		Example: `  # List available plans:
  cherryctl plans list`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			options := c.Servicer.GetOptions()

			if len(types) > 0 {
				options.Type = types
			}

			if regionID != "" {
				options.QueryParams = map[string]string{"region": regionID}
			}

			plans, _, err := c.Service.List(teamID, options)
			if err != nil {
				return errors.Wrap(err, "Could not list plans")
			}

			data := make([][]string, 0)
			for _, p := range plans {
				priceHour := "-"
				priceSpot := "-"
				for _, pricing := range p.Pricing {
					if pricing.Unit == "Hourly" {
						priceHour = fmt.Sprintf("%f", pricing.Price)
					} else if pricing.Unit == "Spot hourly" {
						priceSpot = fmt.Sprintf("%f", pricing.Price)
					}
				}

				for _, r := range p.AvailableRegions {
					if regionID == "" || regionID == r.Slug || regionID == strconv.Itoa(r.ID) {
						data = append(data, []string{p.Slug, r.Slug, strconv.Itoa(r.StockQty), priceHour, strconv.Itoa(r.SpotQty), priceSpot})
					}
				}
			}
			header := []string{"Plan Slug", "Region Slug", "Stock Hourly", "Hourly Price", "Stock Spot", "Spot Price"}

			return c.Out.Output(plans, header, &data)
		},
	}

	planGetCmd.Flags().StringVarP(&regionID, "region", "r", "", "The Slug or ID of region.")
	planGetCmd.Flags().StringSliceVarP(&types, "type", "", []string{}, `Comma separated list of available plan types (baremetal,virtual,vps)`)
	planGetCmd.Flags().IntVarP(&teamID, "team-id", "t", 0, "The team's ID. Return plans prices based on team billing details.")

	_ = planGetCmd.MarkFlagRequired("team-id")

	return planGetCmd
}
