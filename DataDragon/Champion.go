package datadragon

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tools "github.com/yycholla/LeagueBuilder/Tools"
	"github.com/yycholla/LeagueBuilder/lolbuilder" // Import your final character structs
)

// DDragonChampion represents the structure of a single champion from Riot's JSON.
// This is used for unmarshalling the champion.json file.
type DDragonChampion struct {
	ID      string           `json:"id"`
	Key     string           `json:"key"`
	Name    string           `json:"name"`
	Title   string           `json:"title"`
	Image   lolbuilder.Image `json:"image"`
	Lore    string           `json:"lore"`
	Blurb   string           `json:"blurb"`
	Tags    []string         `json:"tags"`
	ParType string           `json:"partype"`
	Info    lolbuilder.Info  `json:"info"`
	Stats   lolbuilder.Stats `json:"stats"`
}

// DDragonResponse is the top-level structure of champion.json.
type DDragonResponse struct {
	Type    string                     `json:"type"`
	Format  string                     `json:"format"`
	Version string                     `json:"version"`
	Data    map[string]DDragonChampion `json:"data"`
}

const championFile = "Champion/champion.json"
const supplementalFile = "Champion/supplemental_abilities.json"

// GetChampionsFile ensures the local champion.json is up-to-date and returns its data.
// It now returns the map of champions, which is more useful for creating characters.
func GetChampionsFile() map[string]DDragonChampion {
	champPath := championFile

	// Ensure the Champion directory exists
	if _, err := os.Stat("Champion"); os.IsNotExist(err) {
		os.Mkdir("Champion", 0755)
	}

	// Check if the champion file exists before trying to read it
	champBytes, err := os.ReadFile(champPath)
	if os.IsNotExist(err) {
		fmt.Println("Champion file not found, downloading...")
		champBytes = nil // Ensure bytes are nil so we trigger the download logic
	} else if err != nil {
		tools.SimpleError(err)
	}

	remoteVersion, err := GetVersion()
	tools.SimpleError(err)

	var champFileResponse DDragonResponse
	// Only try to unmarshal if the file actually exists and has content
	if champBytes != nil {
		err = json.Unmarshal(champBytes, &champFileResponse)
		tools.SimpleError(err)
	}

	fmt.Println("Local Version: ", champFileResponse.Version)
	if champFileResponse.Version != remoteVersion {
		fmt.Println("Updating Champion File to version:", remoteVersion)
		url := "https://ddragon.leagueoflegends.com/cdn/" + remoteVersion + "/data/en_US/champion.json"

		// If the file exists, remove it to ensure a clean write
		if _, err := os.Stat(champPath); err == nil {
			err := os.Remove(champPath)
			tools.SimpleError(err)
		}

		resp, err := http.Get(url)
		tools.SimpleError(err)
		defer resp.Body.Close()

		file, err := os.Create(champPath)
		tools.SimpleError(err)
		defer file.Close()

		// Read the body to a variable so we can unmarshal it after writing
		body, err := io.ReadAll(resp.Body)
		tools.SimpleError(err)

		_, err = file.Write(body)
		tools.SimpleError(err)

		// Unmarshal the new data into our response struct
		err = json.Unmarshal(body, &champFileResponse)
		tools.SimpleError(err)

		fmt.Println("Updated File")
	} else {
		fmt.Println("Champion file up to date")
	}
	return champFileResponse.Data
}

// NewCharacter creates a complete, merged character object by combining data from
// Riot's Data Dragon and the supplemental data scraped from the wiki.
func NewCharacter(name string) (lolbuilder.Character, error) {
	// 1. Load the base champion data from Data Dragon
	ddragonChampions := GetChampionsFile()
	ddragonData, ok := ddragonChampions[name]
	if !ok {
		return lolbuilder.Character{}, fmt.Errorf("champion '%s' not found in Data Dragon file", name)
	}

	// 2. Load the supplemental scraped data
	scrapedBytes, err := os.ReadFile(supplementalFile)
	if err != nil {
		return lolbuilder.Character{}, fmt.Errorf("failed to read supplemental abilities file: %w", err)
	}

	var allScrapedAbilities map[string]lolbuilder.ChampionSupplementalAbilities
	err = json.Unmarshal(scrapedBytes, &allScrapedAbilities)
	if err != nil {
		return lolbuilder.Character{}, fmt.Errorf("failed to unmarshal supplemental abilities: %w", err)
	}

	scrapedData, ok := allScrapedAbilities[name]
	if !ok {
		// It's possible a champion exists in DDragon but wasn't scraped, so we just log this and continue
		log.Printf("Warning: No supplemental scraped data found for champion '%s'. Character will have basic info only.", name)
	}

	// 3. Create the final, merged character struct
	character := lolbuilder.Character{
		// Map all the data from Data Dragon
		ID:      ddragonData.ID,
		Key:     ddragonData.Key,
		Name:    ddragonData.Name,
		Title:   ddragonData.Title,
		Image:   ddragonData.Image,
		Lore:    ddragonData.Lore,
		Blurb:   ddragonData.Blurb,
		Tags:    ddragonData.Tags,
		ParType: ddragonData.ParType,
		Info:    ddragonData.Info,
		Stats:   ddragonData.Stats,

		// Embed the entire scraped abilities struct
		ChampionSupplementalAbilities: scrapedData,
	}

	return character, nil
}
