package datadragon

import (
	"encoding/json"
	"fmt"
	"net/http"

	tools "github.com/yycholla/LeagueBuilder/Tools"
)

type ChampionFile struct {
	Type    string              `json:"type"`
	Format  string              `json:"format"`
	Version string              `json:"version"`
	Data    map[string]Champion `json:"data"`
}

func GetVersion() (string, error) {
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	tools.SimpleError(err)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("unable to get response: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("non 200 response: %v", err)
	}
	var versions []string
	err = json.NewDecoder(resp.Body).Decode(&versions)
	tools.SimpleError(err)

	fmt.Println("Remote version: ", versions[0])

	return versions[0], nil
}
