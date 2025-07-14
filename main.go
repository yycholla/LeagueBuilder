package main

import (
	"github.com/kr/pretty"
	datadragon "github.com/yycholla/LeagueBuilder/DataDragon"
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
	ahri, err := datadragon.NewCharacter("Ahri")
	if err != nil {
		panic(err)
	}
	pretty.Print(ahri)

}
