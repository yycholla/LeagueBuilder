package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	datadragon "github.com/yycholla/LeagueBuilder/DataDragon"
)

// StatDetail holds a single line of an ability's scaling statistics.
type StatDetail struct {
	Type  string `json:"type"`
	Value string `json:"value"`
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
	Stats       []StatDetail `json:"stats,omitempty"`
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

// cleanAndSplitValues takes a raw string, splits it by " / ", and trims whitespace.
func cleanAndSplitValues(input string) []string {
	parts := strings.Split(input, " / ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

func Scrape(championNames datadragon.ChampionFile) {
	AllowedDomains := "leagueoflegends.fandom.com"
	allChampionsData := make(map[string]ChampionSupplementalAbilities)
	re := regexp.MustCompile(`\[edit\]`)

	// Mutex to safely write to the allChampionsData map and the console from multiple goroutines
	var mutex = &sync.Mutex{}

	c := colly.NewCollector(
		colly.AllowedDomains(AllowedDomains),
		colly.Async(true),
	)

	// Set a longer timeout for requests to handle slow server responses.
	c.SetRequestTimeout(30 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		Delay:       1 * time.Second,
	})

	fmt.Println("Running wiki scrape")

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("championData", &ChampionSupplementalAbilities{})
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")

		// Lock the mutex for console output to prevent race conditions on stdout.
		mutex.Lock()
		fmt.Println("Visiting: ", r.URL)
		mutex.Unlock()
	})

	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		fmt.Println("Request URL:", r.Request.URL, "failed with response code:", r.StatusCode, "\nError:", err)
		mutex.Unlock()
	})

	// --- Scraper for the Passive Ability ---
	c.OnHTML("div.skill_innate", func(e *colly.HTMLElement) {
		currentChampion := e.Request.Ctx.GetAny("championData").(*ChampionSupplementalAbilities)

		var passive SupplementalPassive
		h3Clone := e.DOM.Find("h3").Clone()
		h3Clone.Find("aside").Remove()
		passive.Name = strings.TrimSpace(re.ReplaceAllString(h3Clone.Text(), ""))

		var descriptionBuilder strings.Builder
		e.ForEach("div.skill_header + div > div > p", func(_ int, el *colly.HTMLElement) {
			descriptionBuilder.WriteString(strings.TrimSpace(el.Text) + " ")
		})
		passive.Description = strings.TrimSpace(descriptionBuilder.String())

		currentChampion.Passive = passive
	})

	// --- Scraper for Q, W, E, and R abilities ---
	c.OnHTML("div.skill", func(e *colly.HTMLElement) {
		currentChampion := e.Request.Ctx.GetAny("championData").(*ChampionSupplementalAbilities)
		if e.DOM.HasClass("skill_innate") {
			return
		}

		var spell SupplementalSpell
		h3Clone := e.DOM.Find("h3").Clone()
		h3Clone.Find("aside").Remove()
		spell.Name = strings.TrimSpace(re.ReplaceAllString(h3Clone.Text(), ""))

		var descriptionBuilder strings.Builder
		e.ForEach("div.ability-info-container + div > div > p", func(_ int, el *colly.HTMLElement) {
			descriptionBuilder.WriteString(strings.TrimSpace(el.Text) + " ")
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

		var stats []StatDetail
		e.ForEach("dl.skill-tabs", func(_ int, dl *colly.HTMLElement) {
			dl.ForEach("dt", func(i int, dtEl *colly.HTMLElement) {
				ddEl := dtEl.DOM.Next()
				if ddEl.Is("dd") {
					stat := StatDetail{
						Type:  strings.TrimSuffix(strings.TrimSpace(dtEl.Text), ":"),
						Value: strings.TrimSpace(ddEl.Text()),
					}
					stats = append(stats, stat)
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
		data := r.Ctx.GetAny("championData").(*ChampionSupplementalAbilities)

		mutex.Lock()
		allChampionsData[championNameKey] = *data
		fmt.Println(r.Request.URL, " scraped!")
		mutex.Unlock()
	})

	// --- Loop through champions and visit their pages ---
	for championName := range championNames.Data {
		// Store the original key before modifying it for the URL
		originalKey := championName

		// Handle special champion names for URL formatting
		switch championName {
		case "DrMundo":
			championName = "Dr._Mundo"
		case "RekSai":
			championName = "Rek'Sai"
		case "JarvanIV":
			championName = "Jarvan_IV"
		case "KSante":
			championName = "K'Sante"
		case "TahmKench":
			championName = "Tahm_Kench"
		case "AurelionSol":
			championName = "Aurelion_Sol"
		case "MasterYi":
			championName = "Master_Yi"
		case "KogMaw":
			championName = "Kog'Maw"
		case "XinZhao":
			championName = "Xin_Zhao"
		case "MonkeyKing":
			championName = "Wukong"
		case "LeeSin":
			championName = "Lee_Sin"
		case "TwistedFate":
			championName = "Twisted_Fate"
		case "MissFortune":
			championName = "Miss_Fortune"
		}

		url := "https://" + AllowedDomains + "/wiki/" + championName + "/LoL"

		// Use a new context for each request to pass the original champion key
		ctx := colly.NewContext()
		ctx.Put("championName", originalKey)
		c.Request("GET", url, nil, ctx, nil)
	}

	c.Wait()

	// --- Write the final data to a JSON file ---
	file, err := json.MarshalIndent(allChampionsData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal data to JSON: %v", err)
	}

	if _, err := os.Stat("Champion"); os.IsNotExist(err) {
		os.Mkdir("Champion", 0755)
	}

	err = ioutil.WriteFile("Champion/supplemental_abilities.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("\nSuccessfully scraped data and created Champion/supplemental_abilities.json")
}
