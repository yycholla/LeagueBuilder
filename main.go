package main

import (
	"fmt"
)

func main() {
	FetchUpdate()
	fmt.Println("Update check completed.")

	champs, err := GetChampions()
	if err != nil {
		fmt.Println("Error fetching champions:", err)
		return
	}
	for _, champ := range champs {
		fmt.Printf("Champion: %s, Title: %s\n", champ.Name, champ.Title)
	}

	err = FetchAugments()
	if err != nil {
		fmt.Println("Error fetching augments:", err)
		return
	}
	fmt.Println("Augments fetched successfully.")
	augments, err := GetAugments()
	if err != nil {
		fmt.Println("Error fetching augments:", err)
		return
	}
	for _, augment := range augments {
		fmt.Printf("Augment: %s, Description: %s\n", augment.Name, augment.Desc)
	}

	if err != nil {
		fmt.Println("Error fetching items:", err)
		return
	}
	items, err := GetItems()
	if err != nil {
		fmt.Println("Error fetching items:", err)
		return
	}
	for _, item := range items {
		fmt.Printf("Item: %s, Description: %s\n", item.Name, item.Description)
	}
}
