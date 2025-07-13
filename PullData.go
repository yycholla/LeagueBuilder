package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/KnutZuidema/golio/static"
)

type APIClient struct {
	Client          *golio.Client
	AllChampions    []datadragon.ChampionData
	ChampionDetails ChampionDetails
}

type Config struct {
	RiotAPIKey string `json:"riot_api_key"`
}

type ChampionDetails struct {
	ExtendedInfo datadragon.ChampionDataExtended
	Spells       []datadragon.SpellData
	Stats        datadragon.ChampionDataStats
}

func LoadConfig() (*Config, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return nil, fmt.Errorf("could not open config.json")
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}
	return &config, nil
}

type RiotAPI interface {
	GetItems() ([]datadragon.Item, error)
	GetAllChampionData() ([]datadragon.ChampionData, error)
	GetChampionData(championName string) (ChampionDetails, error)
}

func NewApiClient() (*APIClient, error) {
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		config, err := LoadConfig()
		if err != nil {
			return nil, err
		}
		apiKey = config.RiotAPIKey
	}
	client := golio.NewClient(apiKey,
		golio.WithRegion(api.RegionNorthAmerica),
		// Set additional options here IE: Logger
	)

	return &APIClient{
		Client: client,
	}, nil
}

func (c *APIClient) GetItemsForMap(mapID string) ([]datadragon.Item, error) {
	// 1. Fetch the raw item data from the API.
	itemList, err := c.Client.DataDragon.GetItems()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items from API: %w", err)
	}

	// 2. Filter and de-duplicate the items.
	uniqueItemMap := make(map[string]datadragon.Item)
	for _, item := range itemList {
		// We only check if the item is available on the specified map.
		// The `InStore` flag can sometimes be unreliable for certain items.
		if item.Maps[mapID] {
			// Add the item to the map to handle duplicates by name.
			if _, exists := uniqueItemMap[item.Name]; !exists {
				uniqueItemMap[item.Name] = item
			}
		}
	}

	// 3. Convert the map back to a slice.
	filteredItems := make([]datadragon.Item, 0, len(uniqueItemMap))
	for _, item := range uniqueItemMap {
		filteredItems = append(filteredItems, item)
	}

	// 4. Sort the final slice by name.
	sort.Slice(filteredItems, func(i, j int) bool {
		return filteredItems[i].Name < filteredItems[j].Name
	})

	return filteredItems, nil
}

// --- Public Methods ---

// GetSummonersRiftItems now calls the generic getter with the correct map ID.
func (c *APIClient) GetSummonersRiftItems() ([]datadragon.Item, error) {
	return c.GetItemsForMap("11")
}

// GetARAMItems now calls the generic getter with the correct map ID.
func (c *APIClient) GetARAMItems() ([]datadragon.Item, error) {
	return c.GetItemsForMap("12")
}

// GetArenaItems now calls the generic getter with the correct map ID.
func (c *APIClient) GetArenaItems() ([]datadragon.Item, error) {
	return c.GetItemsForMap("30")
}

type AugmentFile struct {
	Augments []Augment `json:"augments"`
}

// Augment represents a single augment object from the JSON array.
type Augment struct {
	APIName      string                 `json:"apiName"`
	Calculations map[string]interface{} `json:"calculations"` // Using interface{} for complex, varied structures
	DataValues   map[string]float64     `json:"dataValues"`
	Description  string                 `json:"desc"`
	IconLarge    string                 `json:"iconLarge"`
	IconSmall    string                 `json:"iconSmall"`
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Rarity       int                    `json:"rarity"`
	Tooltip      string                 `json:"tooltip"`
}

func FetchAugments() ([]Augment, error) {
	url := "https://raw.communitydragon.org/latest/cdragon/arena/en_us.json"

	// Get request augments, err if not
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get augments from cdragon: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, status: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	// Decode returned json
	var augmentsFile AugmentFile
	if err := json.NewDecoder(resp.Body).Decode(&augmentsFile); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	// Sort output by name
	augments := augmentsFile.Augments
	sort.Slice(augments, func(i, j int) bool {
		return augments[i].Name < augments[j].Name
	})

	return augments, nil
}

func (c *APIClient) GetAllChampionData() ([]datadragon.ChampionData, error) {
	// Pull all champ data from riot api, err if bad
	championData, err := c.Client.DataDragon.GetChampions()
	if err != nil {
		return make([]datadragon.ChampionData, 0), fmt.Errorf("unable to get all champions: %w", err)
	}

	// Sort champion slice by name
	sort.Slice(championData, func(i, j int) bool {
		return championData[i].Name < championData[j].Name
	})

	return championData, err
}

func (c *APIClient) GetChampionData(championName string) (*ChampionDetails, error) {
	champion, err := c.Client.DataDragon.GetChampion(championName)
	if err != nil {
		return nil, fmt.Errorf("unable to get champion data: %w", err)
	}

	details := &ChampionDetails{
		ExtendedInfo: champion,
		Spells:       champion.Spells,
		Stats:        champion.Stats,
	}

	return details, nil
}

func (c *APIClient) GetQueueData() ([]static.Queue, error) {
	queues, err := c.Client.Static.GetQueues()
	if err != nil {
		return nil, fmt.Errorf("unable to get queues: %w", err)
	}
	return queues, nil
}

func (c *APIClient) GetMapData() ([]static.Map, error) {
	maps, err := c.Client.Static.GetMaps()
	if err != nil {
		return nil, fmt.Errorf("unable to get maps: %w", err)
	}
	return maps, nil
}

func (c *APIClient) GetGameModeData() ([]static.GameMode, error) {
	gameModes, err := c.Client.Static.GetGameModes()
	if err != nil {
		return nil, fmt.Errorf("unable to get game modes: %w", err)
	}
	return gameModes, nil
}
