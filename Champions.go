package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ChampionFile struct {
	Type    string                    `json:"type"`
	Format  string                    `json:"format"`
	Version string                    `json:"version"`
	Data    map[string]ChampionDetail `json:"data"`
}

type Champion struct {
	ID                   string
	Key                  string
	Name                 string
	Title                string
	Lore                 string
	Blurb                string
	Attack               int
	Defense              int
	Magic                int
	Difficulty           int
	Tags                 []string
	Partype              string
	HP                   float64
	HPPerLevel           float64
	MP                   float64
	MPPerLevel           float64
	Movespeed            float64
	Armor                float64
	ArmorPerLevel        float64
	SpellBlock           float64
	SpellBlockPerLevel   float64
	AttackRange          float64
	HPRegen              float64
	HPRegenPerLevel      float64
	MPRegen              float64
	MPRegenPerLevel      float64
	Crit                 float64
	CritPerLevel         float64
	AttackDamage         float64
	AttackDamagePerLevel float64
	AttackSpeed          float64
	AttackSpeedPerLevel  float64
	Spells               []SpellAPI
	Passive              PassiveAPI
}

type ChampionDetail struct {
	ID          string        `json:"id"`
	Key         string        `json:"key"`
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Image       ImageInfo     `json:"image"`
	Skins       []SkinInfo    `json:"skins"`
	Lore        string        `json:"lore"`
	Blurb       string        `json:"blurb"`
	AllyTips    []string      `json:"allytips"`
	EnemyTips   []string      `json:"enemytips"`
	Tags        []string      `json:"tags"`
	Partype     string        `json:"partype"`
	Info        ChampionInfo  `json:"info"`
	Stats       ChampionStats `json:"stats"`
	Spells      []SpellAPI    `json:"spells"`
	Passive     PassiveAPI    `json:"passive"`
	Recommended []any         `json:"recommended"`
}

type ImageInfo struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

type SkinInfo struct {
	ID      string `json:"id"`
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Chromas bool   `json:"chromas"`
}

type ChampionInfo struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type ChampionStats struct {
	HP                   float64 `json:"hp"`
	HPPerLevel           float64 `json:"hpperlevel"`
	MP                   float64 `json:"mp"`
	MPPerLevel           float64 `json:"mpperlevel"`
	Movespeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlock           float64 `json:"spellblock"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	AttackRange          float64 `json:"attackrange"`
	HPRegen              float64 `json:"hpregen"`
	HPRegenPerLevel      float64 `json:"hpregenperlevel"`
	MPRegen              float64 `json:"mpregen"`
	MPRegenPerLevel      float64 `json:"mpregenperlevel"`
	Crit                 float64 `json:"crit"`
	CritPerLevel         float64 `json:"critperlevel"`
	AttackDamage         float64 `json:"attackdamage"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeed          float64 `json:"attackspeed"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
}

type SpellAPI struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Tooltip      string         `json:"tooltip"`
	LevelTip     SpellLevelTip  `json:"leveltip"`
	MaxRank      int            `json:"maxrank"`
	Cooldown     []float64      `json:"cooldown"`
	CooldownBurn string         `json:"cooldownBurn"`
	Cost         []int          `json:"cost"`
	CostBurn     string         `json:"costBurn"`
	DataValues   map[string]any `json:"datavalues"`
	Effect       [][]float64    `json:"effect"`
	EffectBurn   []string       `json:"effectBurn"`
	Vars         []any          `json:"vars"`
	CostType     string         `json:"costType"`
	MaxAmmo      string         `json:"maxammo"`
	Range        []int          `json:"range"`
	RangeBurn    string         `json:"rangeBurn"`
	Image        ImageInfo      `json:"image"`
	Resource     string         `json:"resource"`
}

type SpellLevelTip struct {
	Label  []string `json:"label"`
	Effect []string `json:"effect"`
}

type PassiveAPI struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       ImageInfo `json:"image"`
}

func GetChampions() ([]Champion, error) {
	champions := []Champion{}
	localVersion, err := GetLocalVersion()
	if err != nil {
		return nil, err
	}
	champDir := filepath.Join("data", "dragontail-"+localVersion, localVersion, "data", "en_US", "champion")

	files, err := os.ReadDir(champDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue // Skip non-JSON files
		}

		path := filepath.Join(champDir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var champFile struct {
			Data map[string]ChampionDetail `json:"data"`
		}
		if err := json.Unmarshal(data, &champFile); err != nil {
			return nil, err
		}

		for _, champion := range champFile.Data {
			champ := Champion{
				ID:                   champion.ID,
				Key:                  champion.Key,
				Name:                 champion.Name,
				Title:                champion.Title,
				Lore:                 champion.Lore,
				Blurb:                champion.Blurb,
				Attack:               champion.Info.Attack,
				Defense:              champion.Info.Defense,
				Magic:                champion.Info.Magic,
				Difficulty:           champion.Info.Difficulty,
				Tags:                 champion.Tags,
				Partype:              champion.Partype,
				HP:                   champion.Stats.HP,
				HPPerLevel:           champion.Stats.HPPerLevel,
				MP:                   champion.Stats.MP,
				MPPerLevel:           champion.Stats.MPPerLevel,
				Movespeed:            champion.Stats.Movespeed,
				Armor:                champion.Stats.Armor,
				ArmorPerLevel:        champion.Stats.ArmorPerLevel,
				SpellBlock:           champion.Stats.SpellBlock,
				SpellBlockPerLevel:   champion.Stats.SpellBlockPerLevel,
				AttackRange:          champion.Stats.AttackRange,
				HPRegen:              champion.Stats.HPRegen,
				HPRegenPerLevel:      champion.Stats.HPRegenPerLevel,
				MPRegen:              champion.Stats.MPRegen,
				MPRegenPerLevel:      champion.Stats.MPRegenPerLevel,
				Crit:                 champion.Stats.Crit,
				CritPerLevel:         champion.Stats.CritPerLevel,
				AttackDamage:         champion.Stats.AttackDamage,
				AttackDamagePerLevel: champion.Stats.AttackDamagePerLevel,
				AttackSpeed:          champion.Stats.AttackSpeed,
				AttackSpeedPerLevel:  champion.Stats.AttackSpeedPerLevel,
				Spells:               champion.Spells,
				Passive:              champion.Passive,
			}
			champions = append(champions, champ)
		}
	}
	return champions, nil
}

func FindChampByName(champions []Champion, name string) (*Champion, error) {
	for _, champ := range champions {
		if strings.EqualFold(champ.Name, name) {
			return &champ, nil
		}
	}
	return nil, fmt.Errorf("Champion not found") // Return nil if no champion found
}
