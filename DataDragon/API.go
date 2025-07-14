package datadragon

import (
	"encoding/json"
	"io"
	"net/http"
)

// GetVersion fetches the latest version string from the Riot API.
func GetVersion() (string, error) {
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versions []string
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return "", err
	}
	return versions[0], nil
}
