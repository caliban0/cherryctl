package plans

import (
	"fmt"
	"strconv"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/spf13/cobra"
)

func (c *Command) get() *cobra.Command {
	var region string

	cmd := cobra.Command{
		Use:   "get SLUG [-r <region_slug>]",
		Short: "Get a server plan.",
		Example: `  # Get a plan
  cherryctl plan get 'B2-1-1gb-20s-shared'
  
  # Get a plan in a specific region
  cherryctl plan get 'B2-1-1gb-20s-shared' -r "LT-Siauliai"
  `,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			opts := c.GetOpts()

			plan, _, err := c.Client().GetBySlug(slug, opts)
			if err != nil {
				return err
			}

			th, td := plansToTable(region, plan)

			c.Outputer().Output(plan, th, &td)
			return nil
		},
	}

	cmd.Flags().StringVarP(&region, "region", "r", "", "The Slug or ID of a region.")

	return &cmd
}

func planPrice(plan cherrygo.Plan, cycle string) float32 {
	for _, p := range plan.Pricing {
		if p.Unit == cycle {
			return p.Price
		}
	}
	return 0
}

// plansToTable prepares the output table of the plans in the given region.
// If `region` is "", all regions are added.
func plansToTable(region string, plans ...cherrygo.Plan) (th []string, td [][]string) {
	th = []string{"plan slug", "region slug", "stock hourly", "hourly price", "stock spot", "spot price"}
	td = make([][]string, 0, len(plans))

	for _, p := range plans {
		hourlyPrice := "-"
		if price := planPrice(p, "Hourly"); price != 0 {
			hourlyPrice = fmt.Sprintf("%f", price)
		}
		spotPrice := "-"
		if price := planPrice(p, "Spot hourly"); price != 0 {
			spotPrice = fmt.Sprintf("%f", price)
		}

		for _, r := range p.AvailableRegions {
			if region == "" || region == r.Slug || region == strconv.Itoa(r.ID) {
				tr := []string{p.Slug, r.Slug, strconv.Itoa(r.StockQty), hourlyPrice, strconv.Itoa(r.SpotQty), spotPrice}
				td = append(td, tr)
			}
		}
	}

	return
}
