package lolbuilder

import (
	"fmt"
)

// Update champion stats {Statistic}}=base+bonus+g\times (n-1)\times (0.7025+0.0175\times (n-1))
func CalculateStatsOnLevel(champion Character, level int) (Character, error) {
	if level < 1 {
		return Character{}, fmt.Errorf("level must be greater than 0")
	}

	if level == 1 {
		return champion, nil
	}

	leveledChampion := champion
	levelFloat := float64(level)
	growthMultiplier := levelFloat * (0.7025 + 0.0175*levelFloat)

	stats := leveledChampion.Stats
	baseStats := champion.Stats

	stats.HP = baseStats.HP + baseStats.HPPerLevel*growthMultiplier
	stats.MP = baseStats.MP + baseStats.MPPerLevel*growthMultiplier
	stats.MoveSpeed = baseStats.MoveSpeed
	stats.Armor = baseStats.Armor + baseStats.ArmorPerLevel*growthMultiplier
	stats.SpellBlock = baseStats.SpellBlock + baseStats.SpellBlockPerLevel*growthMultiplier
	stats.HPRegen = baseStats.HPRegen + baseStats.HPRegenPerLevel*growthMultiplier
	stats.MPRegen = baseStats.MPRegen + baseStats.MPRegenPerLevel*growthMultiplier
	stats.AttackRange = baseStats.AttackRange
	stats.AttackDamage = baseStats.AttackDamage + baseStats.AttackDamagePerLevel*growthMultiplier
	stats.Crit = baseStats.Crit + baseStats.CritPerLevel*growthMultiplier

	stats.AttackSpeed = baseStats.AttackSpeed * (1 + (baseStats.AttackSpeedPerLevel/100)*levelFloat)

	return leveledChampion, nil
}
