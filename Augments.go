package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AugmentsFile struct {
	Augments []Augment `json:"augments"`
}

type Augment struct {
	APIName      string             `json:"apiName"`
	Calculations map[string]any     `json:"calculations"`
	DataValues   map[string]float64 `json:"dataValues"`
	Desc         string             `json:"desc"`
	IconLarge    string             `json:"iconLarge"`
	IconSmall    string             `json:"iconSmall"`
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	Rarity       int                `json:"rarity"`
	Tooltip      string             `json:"tooltip"`
}

func trimPatch(version string) string {
	if idx := strings.LastIndex(version, "."); idx != -1 {
		return version[:idx]
	}
	return version
}

func FetchAugments() error {
	localVersion, err := GetLocalVersion()
	if err != nil {
		return fmt.Errorf("error getting local version: %w", err)
	}
	baseVersion := trimPatch(localVersion)
	url := "https://raw.communitydragon.org/" + baseVersion + "/cdragon/arena/en_us.json"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, status: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	var augmentsFile AugmentsFile
	if err := json.NewDecoder(resp.Body).Decode(&augmentsFile); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	// Save the augments data to a local file
	augmentsDir := filepath.Join("data", "dragontail-"+localVersion, localVersion, "data", "en_US", "augments")
	if err := os.MkdirAll(augmentsDir, 0755); err != nil {
		return fmt.Errorf("error creating augments directory: %w", err)
	}
	augmentsFilePath := filepath.Join(augmentsDir, "augments.json")
	file, err := os.Create(augmentsFilePath)
	if err != nil {
		return fmt.Errorf("error creating augments file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(augmentsFile); err != nil {
		return fmt.Errorf("error writing augments data to file: %w", err)
	}
	return nil
}

func GetAugments() ([]Augment, error) {
	localVersion, err := GetLocalVersion()
	if err != nil {
		return nil, err
	}
	augmentsDir := filepath.Join("data", "dragontail-"+localVersion, localVersion, "data", "en_US", "augments")
	augmentsFilePath := filepath.Join(augmentsDir, "augments.json")

	data, err := os.ReadFile(augmentsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading augments file: %w", err)
	}

	var augmentsFile AugmentsFile
	if err := json.Unmarshal(data, &augmentsFile); err != nil {
		return nil, fmt.Errorf("error unmarshalling augments data: %w", err)
	}

	return augmentsFile.Augments, nil
}
