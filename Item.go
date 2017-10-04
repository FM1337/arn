package arn

const (
	// ItemRarityCommon ...
	ItemRarityCommon = "common"

	// ItemRaritySuperior ...
	ItemRaritySuperior = "superior"

	// ItemRarityRare ...
	ItemRarityRare = "rare"

	// ItemRarityUnique ...
	ItemRarityUnique = "unique"

	// ItemRarityLegendary ...
	ItemRarityLegendary = "legendary"
)

// Item ...
type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Icon        string `json:"icon"`
	Rarity      string `json:"rarity"`
	Order       int    `json:"order"`
	Consumable  bool   `json:"consumable"`
}

// AllItems returns a slice of all items.
func AllItems() []*Item {
	return []*Item{
		&Item{
			ID:    "pro-account-3",
			Name:  "PRO Account (1 season)",
			Price: 900,
			Description: `PRO account for 1 anime season (3 months).

Includes:

* Special highlight on the forums
* Customizable cover image for your profile
* Custom title for profile and forums
* Your suggestions will have a high priority
* Access to the VIP channel on Discord`,
			Icon:       "star",
			Rarity:     ItemRaritySuperior,
			Order:      1,
			Consumable: true,
		},
		&Item{
			ID:    "pro-account-6",
			Name:  "PRO Account (2 seasons)",
			Price: 1600,
			Description: `PRO account for 2 anime seasons (6 months).

Includes:

* Special highlight on the forums
* Customizable cover image for your profile
* Custom title for profile and forums
* Your suggestions will have a high priority
* Access to the VIP channel on Discord`,
			Icon:       "star",
			Rarity:     ItemRarityRare,
			Order:      2,
			Consumable: true,
		},
		&Item{
			ID:    "pro-account-12",
			Name:  "PRO Account (4 seasons)",
			Price: 3000,
			Description: `PRO account for 4 anime seasons (12 months).

Includes:

* Special highlight on the forums
* Customizable cover image for your profile
* Custom title for profile and forums
* Your suggestions will have a high priority
* Access to the VIP channel on Discord`,
			Icon:       "star",
			Rarity:     ItemRarityUnique,
			Order:      3,
			Consumable: true,
		},
		&Item{
			ID:    "pro-account-24",
			Name:  "PRO Account (8 seasons)",
			Price: 5900,
			Description: `PRO account for 8 anime seasons (24 months).

Includes:

* Special highlight on the forums
* Customizable cover image for your profile
* Custom title for profile and forums
* Your suggestions will have a high priority
* Access to the VIP channel on Discord`,
			Icon:       "star",
			Rarity:     ItemRarityLegendary,
			Order:      4,
			Consumable: true,
		},
	}
}

//- ShopItem("PRO Account", "6 months", "1600", "star", strings.Replace(strings.Replace(proAccountMarkdown, "3 months", "6 months", 1), "1 anime season", "2 anime seasons", 1))
//- ShopItem("PRO Account", "1 year", "3000", "star", strings.Replace(strings.Replace(proAccountMarkdown, "3 months", "12 months", 1), "1 anime season", "4 anime seasons", 1))
//- ShopItem("PRO Account", "2 years", "5900", "star", strings.Replace(strings.Replace(proAccountMarkdown, "3 months", "24 months", 1), "1 anime season", "8 anime seasons", 1))
//- ShopItem("Anime Support Ticket", "", "100", "ticket", "Support the makers of your favourite anime by using an anime support ticket. Anime Notifier uses 8% of the money to handle the transaction fees while the remaining 92% go directly to the studios involved in the creation of your favourite anime.")
//- ShopItem("Artwork Support Ticket", "", "100", "ticket", "Support the makers of your favourite artwork by using an artwork support ticket. Anime Notifier uses 8% of the money to handle the transaction fees while the remaining 92% go directly to the creator.")
//- ShopItem("Soundtrack Support Ticket", "", "100", "ticket", "Support the makers of your favourite soundtrack by using a soundtrack support ticket. Anime Notifier uses 8% of the money to handle the transaction fees while the remaining 92% go directly to the creator.")
//- ShopItem("AMV Support Ticket", "", "100", "ticket", "Support the makers of your favourite AMV by using an AMV support ticket. Anime Notifier uses 8% of the money to handle the transaction fees while the remaining 92% go directly to the creator.")
