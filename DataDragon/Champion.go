package datadragon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tools "github.com/yycholla/LeagueBuilder/Tools"
)

type Champion struct {
	Version string   `json:"version"`
	ID      string   `json:"id"`
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	Info    Info     `json:"info"`
	Tags    []string `json:"tags"`
	Spells  []Spell
	Passive Passive
	Stats   BaseStats `json:"stats"`
}

type Info struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type Spell struct {
	ID          string
	Name        string
	Description string
	Tooltip     string
	LevelTip    struct {
		Label  []string
		Effect []string
	}
	MaxRank      int
	Cooldown     []float64
	CooldownBurn string
	Cost         []float64
	CostBurn     string
	effect       [][]float64
}

type Passive struct {
}

type BaseStats struct {
	Hp           float64 `json:"hp"`
	Mp           float64 `json:"mp"`
	MoveSpeed    float64 `json:"movespeed"`
	Armor        float64 `json:"armor"`
	SpellBlock   float64 `json:"spellblock"`
	AttackRange  float64 `json:"attackrange"`
	HpRegen      float64 `json:"hpregen"`
	MpRegen      float64 `json:"mpregen"`
	Crit         float64 `json:"crit"`
	AttackDamage float64 `json:"attackdamage"`
	AttackSpeed  float64 `json:"attackspeed"`
	// Growth stats are also in this same object
	HpPerLevel           float64 `json:"hpperlevel"`
	MpPerLevel           float64 `json:"mpperlevel"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	HpRegenPerLevel      float64 `json:"hpregenperlevel"`
	MpRegenPerLevel      float64 `json:"mpregenperlevel"`
	CritPerLevel         float64 `json:"critperlevel"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
}

const championFile = "Champion/champion.json"

func GetChampionsFile() ChampionFile {
	champPath := "Champion/champion.json"
	champBytes, err := os.ReadFile(champPath)
	tools.SimpleError(err)

	remoteVersion, err := GetVersion()
	tools.SimpleError(err)

	var champFile ChampionFile
	err = json.Unmarshal(champBytes, &champFile)
	tools.SimpleError(err)
	fmt.Println("Local Version: ", champFile.Version) // blank version on output
	if champFile.Version != remoteVersion {
		fmt.Println("Updating Champion File")
		url := "https://ddragon.leagueoflegends.com/cdn/" + remoteVersion + "/data/en_US/champion.json"
		// delete file
		err := os.Remove("Champion/champion.json")
		tools.SimpleError(err)

		// make request
		resp, err := http.Get(url)
		tools.SimpleError(err)
		defer resp.Body.Close()

		// create file
		file, err := os.Create(champPath)
		tools.SimpleError(err)
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		tools.SimpleError(err)
		fmt.Println("Updated File")
	} else {
		fmt.Println("Champion file up to date")
	}
	return champFile
}

// create champion with data present
func NewChampion(name string) Champion {
	file, err := os.ReadFile(championFile)
	tools.SimpleError(err)

	var championFile ChampionFile
	err = json.Unmarshal(file, &championFile)
	tools.SimpleError(err)

	champion := championFile.Data[name]
	return champion
}

// func GetChampionNames()

// test fields for empty values
func CheckChampionFields(champion Champion) {
	champ, err := json.Marshal(champion)
	tools.SimpleError(err)

	var champMap map[string]any
	err = json.Unmarshal(champ, &champMap)
	tools.SimpleError(err)

	for key, value := range champMap {
		fmt.Println(key, value)
	}
}
