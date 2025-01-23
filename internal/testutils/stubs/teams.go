package stubs

import "github.com/cherryservers/cherrygo/v3"

var Team = cherrygo.Team{
	ID:   1,
	Name: "test_team",
	Credit: cherrygo.Credit{
		Account: cherrygo.CreditDetails{
			Remaining: 1.1,
		},
		Promo: cherrygo.CreditDetails{
			Remaining: 1.1,
		},
		Resources: cherrygo.Resources{
			Pricing: cherrygo.Pricing{
				Price: 1.1,
			},
		},
	},
	Billing: cherrygo.Billing{
		Currency: "EUR",
	},
}
