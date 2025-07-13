package main

import (
	"reflect"
	"strings"

	"github.com/KnutZuidema/golio/datadragon"
)

func GetAllItemStats(filteredItems []datadragon.Item) map[string]map[string]float64 {
	allStats := make(map[string]map[string]float64)
	for _, item := range filteredItems {
		statsStruct := reflect.ValueOf(item.Stats)
		statsType := statsStruct.Type()
		currentItemStats := make(map[string]float64)

		for i := 0; i < statsStruct.NumField(); i++ {
			field := statsStruct.Field(i)
			fieldName := statsType.Field(i).Name

			if field.Kind() == reflect.Float64 {
				statValue := field.Float()
				if statValue != 0 {
					currentItemStats[fieldName] = statValue
				}
			}
		}
		if len(currentItemStats) > 0 {
			allStats[item.Name] = currentItemStats
		}
	}
	return allStats
}

func FindBestItemPerStat(items []datadragon.Item) map[string]datadragon.Item {
	bestItemsForStat := make(map[string]datadragon.Item)
	highestStatValues := make(map[string]float64)

	// Get all possible stat names from the ItemStats struct definition.
	statsType := reflect.TypeOf(datadragon.ItemStats{})
	var statNames []string
	for i := 0; i < statsType.NumField(); i++ {
		statNames = append(statNames, statsType.Field(i).Name)
	}

	// Initialize the highest value for each stat to a very small number.
	for _, name := range statNames {
		highestStatValues[name] = 0
	}

	// Iterate through every item in the provided list.
	for _, item := range items {
		// Use reflection to inspect the item's Stats struct.
		statsStruct := reflect.ValueOf(item.Stats)

		// Check each possible stat for the current item.
		for _, statName := range statNames {
			field := statsStruct.FieldByName(statName)
			if field.IsValid() && field.Kind() == reflect.Float64 {
				currentValue := field.Float()

				// If this item's stat is the best we've seen so far, record it.
				if currentValue > highestStatValues[statName] {
					highestStatValues[statName] = currentValue
					bestItemsForStat[statName] = item
				}
			}
		}
	}

	return bestItemsForStat
}

func FilterItemsByTag(items []datadragon.Item, targetTag string) []datadragon.Item {
	var itemList []datadragon.Item
	for _, item := range items {
		for _, tag := range item.Tags {
			if strings.Contains(tag, targetTag) {
				itemList = append(itemList, item)
				continue
			}
		}
	}
	return itemList
}

// func LevelStats(base, bonus datadragon.ChampionDataStats, level float64) (datadragon.ChampionDataStats, error) {
// 	// Handle base stats for character
// 	jsonBytes, err := json.Marshal(base)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to marshal stat struct: %v", err)
// 	}

// 	var statMap map[string]float64
// 	err = json.Unmarshal(jsonBytes, &statMap)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to unmarshal into map: %v", err)
// 	}

// 	// temporary bonus implementation
// 	jsonBytes, err = json.Marshal(bonus)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to marshal stat struct: %v", err)
// 	}

// 	var bonusMap map[string]float64
// 	err = json.Unmarshal(jsonBytes, &bonusMap)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to unmarshal into map: %v", err)
// 	}

// 	// temporary growth implementation
// 	jsonBytes, err = json.Marshal(growth)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to marshal stat struct: %v", err)
// 	}

// 	var growthMap map[string]float64
// 	err = json.Unmarshal(jsonBytes, &growthMap)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to unmarshal into map: %v", err)
// 	}

// 	// calculate stats
// 	for stat, value := range statMap {
// 		if _, ok := statMap[stat+"PerLevel"]; ok {
// 			statMap[stat] = value + bonusMap[stat] + statMap[stat+"PerLevel"]*(level-1)*(0.7025+0.0175*(level-1))
// 		}
// 	}

// 	// convert final back to struct
// 	finalStatsBytes, err := json.Marshal(statMap)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to marshal stat map: %v", err)
// 	}

// 	var finalStats datadragon.ChampionDataStats
// 	err = json.Unmarshal(finalStatsBytes, &finalStats)
// 	if err != nil {
// 		return datadragon.ChampionDataStats{}, fmt.Errorf("unable to unmarshal stat bytes: %v", err)
// 	}

// 	return finalStats, nil
// }
