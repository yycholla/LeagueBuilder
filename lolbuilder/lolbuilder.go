package lolbuilder

import (
	"fmt"
	"strings"
)

// Character is the final, unified struct combining Data Dragon and scraped data.
type Character struct {
	ID      string   `json:"id"`
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	Image   Image    `json:"image"`
	Lore    string   `json:"lore"`
	Blurb   string   `json:"blurb"`
	Tags    []string `json:"tags"`
	ParType string   `json:"partype"`
	Info    Info     `json:"info"`
	Stats   Stats    `json:"stats"`
	// Embed the detailed abilities from your scraper
	ChampionSupplementalAbilities
}

// Image holds all the image data for a champion from Data Dragon.
type Image struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
}

// Info holds the informational ratings for a champion from Data Dragon.
type Info struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

// Stats holds the base stat information for a champion from Data Dragon.
type Stats struct {
	HP                   float64 `json:"hp"`
	MP                   float64 `json:"mp"`
	MoveSpeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	SpellBlock           float64 `json:"spellblock"`
	HPRegen              float64 `json:"hpregen"`
	MPRegen              float64 `json:"mpregen"`
	AttackRange          float64 `json:"attackrange"`
	AttackDamage         float64 `json:"attackdamage"`
	AttackSpeed          float64 `json:"attackspeed"`
	Crit                 float64 `json:"crit"`
	HPPerLevel           float64 `json:"hpperlevel"`
	MPPerLevel           float64 `json:"mpperlevel"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	CritPerLevel         float64 `json:"critperlevel"`
	HPRegenPerLevel      float64 `json:"hpregenperlevel"`
	MPRegenPerLevel      float64 `json:"mpregenperlevel"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
}

// --- Your Existing Scraper Structs (unchanged) ---

// ParsedStat holds the structured data for a single ability statistic.
type ParsedStat struct {
	Type        string              `json:"type"`
	Description string              `json:"description"`
	Values      []string            `json:"values"`
	Modifiers   map[string][]string `json:"modifiers"`
	Notes       string              `json:"notes,omitempty"`
}

// SupplementalSpell defines the structure for a champion's spell.
type SupplementalSpell struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Cost        []string     `json:"cost,omitempty"`
	Cooldown    []string     `json:"cooldown,omitempty"`
	CastTime    string       `json:"cast_time,omitempty"`
	Range       []string     `json:"range,omitempty"`
	Speed       []string     `json:"speed,omitempty"`
	Width       []string     `json:"width,omitempty"`
	Stats       []ParsedStat `json:"stats,omitempty"`
}

// SupplementalPassive defines the structure for a champion's passive ability.
type SupplementalPassive struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ChampionSupplementalAbilities holds all the scraped ability data for one champion.
type ChampionSupplementalAbilities struct {
	Passive SupplementalPassive `json:"passive"`
	Q       SupplementalSpell   `json:"q"`
	W       SupplementalSpell   `json:"w"`
	E       SupplementalSpell   `json:"e"`
	R       SupplementalSpell   `json:"r"`
}

// String creates a human-readable summary of the Character.
// This method satisfies the fmt.Stringer interface.
func (c *Character) String() string {
	// Using strings.Builder is highly efficient for building strings.
	var sb strings.Builder

	// Main Identity
	sb.WriteString(fmt.Sprintf("%s - %s\n", c.Name, c.Title))
	sb.WriteString(fmt.Sprintf("\"%s\"\n", c.Blurb))
	sb.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(c.Tags, ", ")))
	sb.WriteString("\n")

	// Base Stats Summary
	sb.WriteString("--- Base Stats ---\n")
	sb.WriteString(fmt.Sprintf("  Health:           %.2f (+%.2f per level)\n", c.Stats.HP, c.Stats.HPPerLevel))
	sb.WriteString(fmt.Sprintf("  %s: %-15.2f (+%.2f per level)\n", c.ParType, c.Stats.MP, c.Stats.MPPerLevel)) // Use ParType for Mana/Energy etc.
	sb.WriteString(fmt.Sprintf("  Attack Damage:    %.2f (+%.2f per level)\n", c.Stats.AttackDamage, c.Stats.AttackDamagePerLevel))
	sb.WriteString(fmt.Sprintf("  Armor:            %.2f (+%.2f per level)\n", c.Stats.Armor, c.Stats.ArmorPerLevel))
	sb.WriteString(fmt.Sprintf("  Magic Resist:     %.2f (+%.2f per level)\n", c.Stats.SpellBlock, c.Stats.SpellBlockPerLevel))
	sb.WriteString(fmt.Sprintf("  Movement Speed:   %.0f\n", c.Stats.MoveSpeed))
	sb.WriteString("\n")

	// Abilities
	sb.WriteString("--- Abilities ---\n")
	sb.WriteString(fmt.Sprintf("  Passive: %s\n", c.Passive.Name))
	sb.WriteString(fmt.Sprintf("  Q:       %s\n", c.Q.Name))
	sb.WriteString(fmt.Sprintf("  W:       %s\n", c.W.Name))
	sb.WriteString(fmt.Sprintf("  E:       %s\n", c.E.Name))
	sb.WriteString(fmt.Sprintf("  R:       %s\n", c.R.Name))

	return sb.String()
}
