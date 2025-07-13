package main

import (
	"fmt"
)

func main() {
	c, err := NewApiClient()
	if err != nil {
		fmt.Println("Error creating client: ", err)
		return
	}
	if err != nil {
		fmt.Println("Error grabbing arena items: ", err)
		return
	}

	// champions, err := c.GetAllChampionData()
	if err != nil {
		fmt.Println("Error getting champions: ", err)
	}

	arenaItems, err := c.GetArenaItems()
	if err != nil {
		return
	}
	for _, item := range arenaItems {
		fmt.Println(item.Name)
	}
	fmt.Println("found ", len(arenaItems), "arena items")
}
