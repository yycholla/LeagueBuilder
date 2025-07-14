package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	datadragon "github.com/yycholla/LeagueBuilder/DataDragon"
	lolbuilder "github.com/yycholla/LeagueBuilder/lolbuilder"
)

// cleanAndSplitValues takes a raw string, splits it by " / ", and trims whitespace.
func cleanAndSplitValues(input string) []string {
	parts := strings.Split(input, " / ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

// parseStatValue is the core logic to transform a stat string into the structured ParsedStat model.
func parseStatValue(statType string, ddSelection *goquery.Selection) (lolbuilder.ParsedStat, error) {
	// Clone the selection so we can modify it without affecting other operations.
	ddClone := ddSelection.Clone()

	// **FIX**: Remove ALL known tooltip elements before processing the text.
	// This now includes .pp-tooltip for scaling stats and .glossary for keywords.
	ddClone.Find("span.ll-item.navbox, .pp-tooltip, .glossary").Remove()

	// Get the fully cleaned text.
	cleanValueStr := ddClone.Text()

	// **FIX**: Use the clean string for BOTH the description and for parsing.
	parsed := lolbuilder.ParsedStat{
		Type:        statType,
		Description: strings.TrimSpace(cleanValueStr), // Use the clean string here
		Modifiers:   make(map[string][]string),
	}

	modifierRegex := regexp.MustCompile(`\(\+\s*([^)]+)\)`)
	notesAndBases := modifierRegex.ReplaceAllString(cleanValueStr, "")

	matches := modifierRegex.FindAllStringSubmatch(cleanValueStr, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		modText := strings.TrimSpace(match[1])

		keys := []string{"bonus AD", "bonus health", "of target's maximum health", "AD", "AP"}
		foundKey := "unknown"

		for _, k := range keys {
			if strings.Contains(modText, k) {
				foundKey = k
				break
			}
		}

		valuesPart := strings.ReplaceAll(modText, foundKey, "")
		valuesPart = strings.ReplaceAll(valuesPart, "%", "")
		values := strings.Split(valuesPart, "/")
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}
		parsed.Modifiers[foundKey] = values
	}

	baseValuesRegex := regexp.MustCompile(`^[\d\s/.-]+`)
	baseValuesMatch := baseValuesRegex.FindString(notesAndBases)

	if baseValuesMatch != "" {
		cleanedBases := strings.Trim(baseValuesMatch, " /")
		// Further cleanup to handle non-breaking spaces sometimes used as separators
		cleanedBases = strings.ReplaceAll(cleanedBases, " − ", "/")
		parsed.Values = strings.Split(cleanedBases, "/")
		for i, v := range parsed.Values {
			parsed.Values[i] = strings.TrimSpace(v)
		}
		parsed.Notes = strings.TrimSpace(strings.Replace(notesAndBases, baseValuesMatch, "", 1))
	} else {
		parsed.Notes = strings.TrimSpace(notesAndBases)
	}

	// Final regex cleanup on notes for stray characters or phrases
	spaceCleanupRegex := regexp.MustCompile(`\s{2,}`)
	parsed.Notes = spaceCleanupRegex.ReplaceAllString(parsed.Notes, " ")

	if len(parsed.Values) > 0 && strings.HasSuffix(cleanValueStr, "seconds") {
		parsed.Notes = "seconds"
	}
	if strings.HasSuffix(parsed.Description, "%") && len(parsed.Values) > 0 {
		parsed.Notes = "%"
	}

	return parsed, nil
}

// Scrape initiates the scraping process for all champions.
func Scrape(championNames map[string]datadragon.DDragonChampion) {
	allowedDomains := "leagueoflegends.fandom.com"
	allChampionsData := make(map[string]lolbuilder.ChampionSupplementalAbilities)
	re := regexp.MustCompile(`\[edit\]`)

	var mutex = &sync.Mutex{}
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains),
		colly.Async(true),
	)

	c.SetRequestTimeout(30 * time.Second)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		Delay:       1 * time.Second,
	})

	fmt.Println("Running wiki scrape")

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("championData", &lolbuilder.ChampionSupplementalAbilities{})
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		mutex.Lock()
		fmt.Println("Visiting: ", r.URL)
		mutex.Unlock()
	})

	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		fmt.Println("Request URL:", r.Request.URL, "failed with response code:", r.StatusCode, "\nError:", err)
		mutex.Unlock()
	})

	c.OnHTML("div.skill_innate", func(e *colly.HTMLElement) {
		currentChampion := e.Request.Ctx.GetAny("championData").(*lolbuilder.ChampionSupplementalAbilities)
		var passive lolbuilder.SupplementalPassive

		h3Selection := e.DOM.Find("h3").First()
		h3Clone := h3Selection.Clone()
		h3Clone.Find("aside").Remove()
		passive.Name = strings.TrimSpace(re.ReplaceAllString(h3Clone.Text(), ""))

		var descriptionBuilder strings.Builder
		e.DOM.Find("p").Each(func(_ int, p *goquery.Selection) {
			if p.Closest("dl.skill-tabs").Length() == 0 {
				pClone := p.Clone()
				pClone.Find("span.ll-item.navbox").Remove()
				paragraphText := strings.TrimSpace(pClone.Text())
				spaceCleanupRegex := regexp.MustCompile(`\s{2,}`)
				cleanedText := spaceCleanupRegex.ReplaceAllString(paragraphText, " ")
				descriptionBuilder.WriteString(cleanedText)
				descriptionBuilder.WriteString(" ")
			}
		})
		passive.Description = strings.TrimSpace(descriptionBuilder.String())
		currentChampion.Passive = passive
	})

	c.OnHTML("div.skill", func(e *colly.HTMLElement) {
		currentChampion := e.Request.Ctx.GetAny("championData").(*lolbuilder.ChampionSupplementalAbilities)
		if e.DOM.HasClass("skill_innate") {
			return
		}

		var spell lolbuilder.SupplementalSpell

		h3Selection := e.DOM.Find("h3").First()
		h3Clone := h3Selection.Clone()
		h3Clone.Find("aside").Remove()
		spell.Name = strings.TrimSpace(re.ReplaceAllString(h3Clone.Text(), ""))

		var descriptionBuilder strings.Builder
		e.DOM.Find("p").Each(func(_ int, p *goquery.Selection) {
			if p.Closest("dl.skill-tabs").Length() == 0 {
				pClone := p.Clone()
				pClone.Find("span.ll-item.navbox").Remove()
				paragraphText := strings.TrimSpace(pClone.Text())
				spaceCleanupRegex := regexp.MustCompile(`\s{2,}`)
				cleanedText := spaceCleanupRegex.ReplaceAllString(paragraphText, " ")
				descriptionBuilder.WriteString(cleanedText)
				descriptionBuilder.WriteString(" ")
			}
		})
		spell.Description = strings.TrimSpace(descriptionBuilder.String())

		e.ForEach(".pi-item.pi-data", func(_ int, el *colly.HTMLElement) {
			label := strings.ToLower(strings.TrimSpace(el.ChildText("h3")))
			value := strings.TrimSpace(el.ChildText(".pi-data-value"))
			switch label {
			case "cost:":
				spell.Cost = cleanAndSplitValues(value)
			case "cooldown:":
				spell.Cooldown = cleanAndSplitValues(value)
			case "cast time:":
				spell.CastTime = value
			case "range:", "target range:":
				spell.Range = cleanAndSplitValues(value)
			case "speed:":
				spell.Speed = cleanAndSplitValues(value)
			case "width:":
				spell.Width = cleanAndSplitValues(value)
			}
		})

		var stats []lolbuilder.ParsedStat
		e.ForEach("dl.skill-tabs", func(_ int, dl *colly.HTMLElement) {
			dl.ForEach("dt", func(i int, dtEl *colly.HTMLElement) {
				ddEl := dtEl.DOM.Next()
				if ddEl.Is("dd") {
					statType := strings.TrimSuffix(strings.TrimSpace(dtEl.Text), ":")
					parsedStat, err := parseStatValue(statType, ddEl)
					if err != nil {
						log.Printf("Could not parse stat for %s: %v", spell.Name, err)
						return
					}
					stats = append(stats, parsedStat)
				}
			})
		})
		spell.Stats = stats

		if e.DOM.HasClass("skill_q") {
			currentChampion.Q = spell
		} else if e.DOM.HasClass("skill_w") {
			currentChampion.W = spell
		} else if e.DOM.HasClass("skill_e") {
			currentChampion.E = spell
		} else if e.DOM.HasClass("skill_r") {
			currentChampion.R = spell
		}
	})

	c.OnScraped(func(r *colly.Response) {
		championNameKey := r.Ctx.Get("championName")
		data := r.Ctx.GetAny("championData").(*lolbuilder.ChampionSupplementalAbilities)
		mutex.Lock()
		allChampionsData[championNameKey] = *data
		fmt.Println(r.Request.URL, " scraped!")
		mutex.Unlock()
	})

	for championName := range championNames {
		originalKey := championName
		urlName := strings.ReplaceAll(championName, " ", "_")

		switch championName {
		case "AurelionSol":
			urlName = "Aurelion_Sol"
		case "BelVeth":
			urlName = "Bel'Veth"
		case "ChoGath":
			urlName = "Cho'Gath"
		case "DrMundo":
			urlName = "Dr._Mundo"
		case "JarvanIV":
			urlName = "Jarvan_IV"
		case "KaiSa":
			urlName = "Kai'Sa"
		case "KhaZix":
			urlName = "Kha'Zix"
		case "KogMaw":
			urlName = "Kog'Maw"
		case "KSante":
			urlName = "K'Sante"
		case "LeeSin":
			urlName = "Lee_Sin"
		case "MasterYi":
			urlName = "Master_Yi"
		case "MissFortune":
			urlName = "Miss_Fortune"
		case "MonkeyKing":
			urlName = "Wukong"
		case "Nunuwillump":
			urlName = "Nunu_%26_Willump"
		case "RekSai":
			urlName = "Rek'Sai"
		case "Renata":
			urlName = "Renata_Glasc"
		case "TahmKench":
			urlName = "Tahm_Kench"
		case "TwistedFate":
			urlName = "Twisted_Fate"
		case "VelKoz":
			urlName = "Vel'Koz"
		case "XinZhao":
			urlName = "Xin_Zhao"
		}

		url := fmt.Sprintf("https://%s/wiki/%s/LoL", allowedDomains, urlName)
		ctx := colly.NewContext()
		ctx.Put("championName", originalKey)
		c.Request("GET", url, nil, ctx, nil)
	}

	c.Wait()

	file, err := json.MarshalIndent(allChampionsData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal data to JSON: %v", err)
	}

	if _, err := os.Stat("Champion"); os.IsNotExist(err) {
		os.Mkdir("Champion", 0755)
	}

	err = os.WriteFile("Champion/supplemental_abilities.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("\nSuccessfully scraped and parsed data, created Champion/supplemental_abilities.json")
}
