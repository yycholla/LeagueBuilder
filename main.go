package main

import (
	"fmt"
)

const selectedChampion = "Ahri"
const workers = 1000000
const batchSize = 10000

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// champions, err := GetChampions()
	// if err != nil {
	// 	fmt.Println("Error fetching champions:", err)
	// 	return
	// }
	// items, err := GetItems()
	// if err != nil {
	// 	fmt.Println("Error fetching items:", err)
	// 	return
	// }
	// fmt.Println(len(items))                                       // Example usage of items
	// selected, err := FindChampByName(champions, selectedChampion) // Replace "Ahri" with the desired champion name
	// if err != nil {
	// 	fmt.Println("Error searching for champion:", err)
	// 	return
	// }
	// if selected == nil {
	// 	fmt.Println("No champion selected")
	// 	return
	// }
	// // processChamp := selected

	// // testItems(*processChamp, items)
	// // ProcessAllCombinationsOptimized(items, 6, 10000, 10000)
	// // CountAllCombinations(619, 6)

	// // bestItems, err := testItems(*selected)
	// // if err != nil {
	// // 	fmt.Println("Unable to test items")
	// // 	return
	// // }
	// // count := 1
	// // for stat, item := range bestItems {
	// // 	fmt.Printf("%d | Items: %s - Value: %v \n", count, stat, item.Name)
	// // 	count++
	// // }
	// fmt.Println(GetUniqueStatNames(""))

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
