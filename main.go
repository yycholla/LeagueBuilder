package main

import (
	datadragon "github.com/yycholla/LeagueBuilder/DataDragon"
	scraper "github.com/yycholla/LeagueBuilder/Scraper"
)

const (
	SPELLBLOCK      = "SpellBlock"
	BOOTS           = "Boots"
	MANAREGEN       = "ManaRegen"
	HEALTHREGEN     = "HealthRegin"
	MANA            = "Mana"
	Health          = "Health"
	Armor           = "Armor"
	SPELLDAMAGE     = "SpellDamage"
	LIFESTEAL       = "LifeSteal"
	SPELLVAMP       = "SpellVamp"
	JUNGLE          = "Jungle"
	DAMAGE          = "Damage"
	LANE            = "Lane"
	ATTACKSPEED     = "AttackSpeed"
	ONHIT           = "OnHit"
	TRINKET         = "Trinket"
	ACTIVE          = "Active"
	CONSUMABLE      = "Consumable"
	CDREDUCTION     = "CooldownReduction"
	ARMORPEN        = "ArmorPenetration"
	HASTE           = "AbilityHaste"
	STEALTH         = "Stealth"
	VISION          = "Vision"
	NONBOOTMOVEMENT = "NonbootMovement"
	TENACITY        = "Tenacity"
	MAGICPEN        = "MagicPenetration"
	CRITSTRIKE      = "CriticalStrike"
)

func main() {
	// champion := datadragon.NewChampion("Ahri")
	// datadragon.CheckChampionFields(champion)
	file := datadragon.GetChampionsFile()
	scraper.Scrape(file)
}
