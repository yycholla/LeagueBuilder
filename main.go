package main

import (
	"fmt"
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
	CRITCHANCE      = ""
)

func main() {
	c, err := NewApiClient()
	if err != nil {
		fmt.Println("Error creating client: ", err)
		return
	}
	if err != nil {
		fmt.Println("Error grabbing arena items: ", err)
		return
	}

	// champions, err := c.GetAllChampionData()
	if err != nil {
		fmt.Println("Error getting champions: ", err)
	}

	arenaItems, err := c.GetArenaItems()
	if err != nil {
		return
	}
	for itemName, item := range FindBestItemPerStat(arenaItems) {
		fmt.Println(itemName, ": ", item.Name)
	}
	for _, item := range FilterItemsByTag(arenaItems, SPELLBLOCK) {
		fmt.Println(item.Name)
	}
	fmt.Println(c.GetVersion())
}
