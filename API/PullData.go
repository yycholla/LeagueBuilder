package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/KnutZuidema/golio/static"
)

type APIClient struct {
	Client    *golio.Client
	allItems  []datadragon.Item
	itemMutex sync.RWMutex
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

func (c *APIClient) getCachedItems() error {
	// First, try a fast read-only check to see if the cache is already populated.
	c.itemMutex.RLock()
	if c.allItems != nil {
		c.itemMutex.RUnlock()
		return nil // Cache is hot, nothing to do.
	}
	c.itemMutex.RUnlock()

	// If we got here, the cache was empty. Now get a full write lock to populate it.
	c.itemMutex.Lock()
	defer c.itemMutex.Unlock() // Ensure the lock is always released

	// It's possible another goroutine populated the cache while we waited for the lock.
	// We do a "double-check" to prevent redundant API calls.
	if c.allItems != nil {
		return nil
	}

	// The cache is definitely empty, so let's make the API call.
	itemList, err := c.Client.DataDragon.GetItems()
	if err != nil {
		return fmt.Errorf("failed to fetch items from API: %w", err)
	}

	// Sort the slice for consistent ordering.
	sort.Slice(itemList, func(i, j int) bool {
		return itemList[i].Name < itemList[j].Name
	})

	// Store the result in the cache.
	c.allItems = itemList
	return nil
}

// filterAndUniqueItems is a standalone helper function to process a list of items.
// It filters by map ID and ensures the returned list has unique names.
func filterAndUniqueItemsChampion(allItems []datadragon.Item, mapID string, championName string) []datadragon.Item {
	uniqueItemMap := make(map[string]datadragon.Item)

	for _, item := range allItems {
		if item.Maps[mapID] && (item.RequiredChampion == "" || item.RequiredChampion == championName) {
			if _, exists := uniqueItemMap[item.Name]; !exists {
				uniqueItemMap[item.Name] = item
			}
		}
	}

	// Convert the map back to a slice
	filteredItems := make([]datadragon.Item, 0, len(uniqueItemMap))
	for _, item := range uniqueItemMap {
		filteredItems = append(filteredItems, item)
	}

	// Sort the final slice by name
	sort.Slice(filteredItems, func(i, j int) bool {
		return filteredItems[i].Name < filteredItems[j].Name
	})

	return filteredItems
}

func filterAndUniqueItems(allItems []datadragon.Item, mapID string) []datadragon.Item {
	return filterAndUniqueItemsChampion(allItems, mapID, "")
}

// GetSummonersRiftItems ensures the master list is cached, then returns a filtered list.
func (c *APIClient) GetSummonersRiftItems() ([]datadragon.Item, error) {
	if err := c.getCachedItems(); err != nil {
		return nil, err
	}
	return filterAndUniqueItems(c.allItems, "11"), nil
}

// GetARAMItems ensures the master list is cached, then returns a filtered list.
func (c *APIClient) GetARAMItems() ([]datadragon.Item, error) {
	if err := c.getCachedItems(); err != nil {
		return nil, err
	}
	return filterAndUniqueItems(c.allItems, "12"), nil
}

// GetArenaItems ensures the master list is cached, then returns a filtered list.
func (c *APIClient) GetArenaItems() ([]datadragon.Item, error) {
	if err := c.getCachedItems(); err != nil {
		return nil, err
	}
	return filterAndUniqueItems(c.allItems, "30"), nil
}

func (c *APIClient) GetBasicItem() ([]datadragon.Item, error) {
	if err := c.getCachedItems(); err != nil {
		return nil, err
	}
	for _, item := range c.allItems {
		fmt.Println(item.Name)
	}
	return c.allItems, nil

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

func (c *APIClient) GetVersion() string {
	return c.Client.DataDragon.Version
}

func (c *APIClient) WriteJsonToFile() error {
	exampleChampion := "Ahri"
	os.MkdirAll("data", os.FileMode.Perm(0755))

	// try to get items from api
	itemData, itemErr := c.Client.DataDragon.GetItems()
	allChampData, allChampErr := c.Client.DataDragon.GetChampions()
	exampleChampionData, exChampErr := c.Client.DataDragon.GetChampion(exampleChampion)
	exampleChampionURL := "https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default/v1/champions/" + exampleChampionData.ID + ".json"
	exChampionData, exErr := http.Get(exampleChampionURL)
	if itemErr != nil || allChampErr != nil || exChampErr != nil || exErr != nil {
		return fmt.Errorf("unable to fetch data %w%w%w", itemErr, allChampErr, exChampErr)
	}
	fmt.Println("WRITEAPITOFILE: Got data from API")

	// try to marshal data to json
	itemJson, itemErr := json.MarshalIndent(itemData, "", "   ")
	allChampJson, allChampErr := json.MarshalIndent(allChampData, "", "   ")
	exampleChampionJson, exChampErr := json.MarshalIndent(exChampionData, "", "   ")
	if itemErr != nil || allChampErr != nil || exChampErr != nil {
		return fmt.Errorf("unable to marshal data %w%w%w", itemErr, allChampErr, exChampErr)
	}
	fmt.Println("WRITEAPITOFILE: Wrote data to json")

	// try to write files
	itemWriteErr := os.WriteFile("data/items.json", itemJson, 0644)
	allChamWriteErr := os.WriteFile("data/champions.json", allChampJson, 0644)
	exChampWriteErr := os.WriteFile("data/examplechampion.json", exampleChampionJson, 0644)
	if itemWriteErr != nil || allChamWriteErr != nil || exChampWriteErr != nil {
		return fmt.Errorf("unable to write data to file: %w%w%w", itemWriteErr, allChamWriteErr, exChampWriteErr)
	}
	fmt.Println("WRITEAPITOFILE: Wrote json to file")
	return nil

}
