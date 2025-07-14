package lolbuilder

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
	HPPerLevel           float64 `json:"hpperlevel"`
	MP                   float64 `json:"mp"`
	MPPerLevel           float64 `json:"mpperlevel"`
	MoveSpeed            float64 `json:"movespeed"`
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
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
	AttackSpeed          float64 `json:"attackspeed"`
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
