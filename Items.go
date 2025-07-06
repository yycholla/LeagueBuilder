package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ItemFile struct {
	Type    string                `json:"type"`
	Version string                `json:"version"`
	Basic   ItemDetail            `json:"basic"`
	Data    map[string]ItemDetail `json:"data"`
}

type Item struct {
	ID               string
	Name             string
	Description      string
	Colloq           string
	Plaintext        string
	From             []string
	Into             []string
	Image            ItemImage
	Gold             ItemGold
	Tags             []string
	Maps             map[string]bool
	Stats            map[string]float64
	Effect           map[string]string
	Depth            int
	Stacks           int
	Consumed         bool
	ConsumeOnFull    bool
	InStore          bool
	HideFromAll      bool
	RequiredChampion string
	RequiredAlly     string
	SpecialRecipe    int
}

type ItemDetail struct {
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Colloq           string             `json:"colloq"`
	Plaintext        string             `json:"plaintext"`
	From             []string           `json:"from"`
	Into             []string           `json:"into"`
	Image            ItemImage          `json:"image"`
	Gold             ItemGold           `json:"gold"`
	Tags             []string           `json:"tags"`
	Maps             map[string]bool    `json:"maps"`
	Stats            map[string]float64 `json:"stats"`
	Effect           map[string]string  `json:"effect"`
	Depth            int                `json:"depth"`
	Stacks           int                `json:"stacks"`
	Consumed         bool               `json:"consumed"`
	ConsumeOnFull    bool               `json:"consumeOnFull"`
	InStore          bool               `json:"inStore"`
	HideFromAll      bool               `json:"hideFromAll"`
	RequiredChampion string             `json:"requiredChampion"`
	RequiredAlly     string             `json:"requiredAlly"`
	SpecialRecipe    int                `json:"specialRecipe"`
}

type ItemImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

type ItemGold struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}

// Example loader function, similar to GetChampions
func GetItems() ([]Item, error) {
	localVersion, err := GetLocalVersion()
	if err != nil {
		return nil, err
	}
	itemPath := filepath.Join("data", "dragontail-"+localVersion, localVersion, "data", "en_US", "item.json")
	data, err := os.ReadFile(itemPath)
	if err != nil {
		return nil, err
	}
	var itemFile ItemFile
	if err := json.Unmarshal(data, &itemFile); err != nil {
		return nil, err
	}
	items := []Item{}
	for id, detail := range itemFile.Data {
		item := Item{
			ID:               id,
			Name:             detail.Name,
			Description:      detail.Description,
			Colloq:           detail.Colloq,
			Plaintext:        detail.Plaintext,
			From:             detail.From,
			Into:             detail.Into,
			Image:            detail.Image,
			Gold:             detail.Gold,
			Tags:             detail.Tags,
			Maps:             detail.Maps,
			Stats:            detail.Stats,
			Effect:           detail.Effect,
			Depth:            detail.Depth,
			Stacks:           detail.Stacks,
			Consumed:         detail.Consumed,
			ConsumeOnFull:    detail.ConsumeOnFull,
			InStore:          detail.InStore,
			HideFromAll:      detail.HideFromAll,
			RequiredChampion: detail.RequiredChampion,
			RequiredAlly:     detail.RequiredAlly,
			SpecialRecipe:    detail.SpecialRecipe,
		}
		items = append(items, item)
	}
	return items, nil
}
