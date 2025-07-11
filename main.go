package main

import (
	"fmt"
	"runtime"
)

const selectedChampion = "Ahri"
const workers = 1000000
const batchSize = 10000

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	champions, err := GetChampions()
	if err != nil {
		fmt.Println("Error fetching champions:", err)
		return
	}
	items, err := GetItems()
	if err != nil {
		fmt.Println("Error fetching items:", err)
		return
	}
	fmt.Println(items[1].Name)                                     // Example usage of items
	selected := SearchChampionsByName(champions, selectedChampion) // Replace "Ahri" with the desired champion name
	if err != nil {
		fmt.Println("Error searching for champion:", err)
		return
	}
	if selected != nil {
		fmt.Println("You selected:", selected.Name, "")
		fmt.Println()
		ProcessAllCombinationsBatched(items, 6, workers, batchSize)

	}

}
