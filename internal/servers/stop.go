package servers

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func (c *Client) Stop() *cobra.Command {
	stopServerCmd := &cobra.Command{
		Use:   `stop ID`,
		Args:  cobra.ExactArgs(1),
		Short: "Stop a server.",
		Long:  "Stops or powers off a server that is currently powered on.",
		Example: `  # Stops the specified server:
  cherryctl server stop 12345`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if serverID, err := strconv.Atoi(args[0]); err == nil {
				_, _, err := c.Service.PowerOff(serverID)
				if err != nil {
					return errors.Wrap(err, "Could not stop a Server")
				}

				c.Out.Outputln("Server", serverID, "successfully stopped.")
				return nil
			} else {
				return errors.Wrap(err, `invalid server ID`)
			}
		},
	}

	return stopServerCmd
}
