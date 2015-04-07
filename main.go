package main

import (
	"encoding/json"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/andlabs/ui"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type ItemData []struct {
	Name             string `json:"name"`
	Shown            bool   `json:"shown,omitempty"`
	Tears            string `json:"tears,omitempty"`
	Delay            string `json:"delay,omitempty"`
	DelayMultiplier  string `json:"delayx,omitempty"`
	Text             string `json:"text,omitempty"`
	Range            string `json:"range,omitempty"`
	Height           string `json:"height,omitempty"`
	Damage           string `json:"damage,omitempty"`
	DamageMultiplier string `json:"damagex,omitempty"`
	Speed            string `json:"speed,omitempty"`
	Health           string `json:"health,omitempty"`
	Space            bool   `json:"space,omitempty"`
	SoulHearts       string `json:"soulhearts,omitempty"`
	ID               string `json:"id"`
}

type RunData struct {
	Seed             string
	Tears            float64
	Delay            float64
	DelayMultiplier  float64
	Range            float64
	Height           float64
	Damage           float64
	DamageMultiplier float64
	Speed            float64
	ItemIDs          []string
}

type LastItemStat struct {
	Name             string
	Tears            float64
	Delay            float64
	DelayMultiplier  float64
	Range            float64
	Height           float64
	Damage           float64
	DamageMultiplier float64
	Speed            float64
	Description      string
}

var SeedLabel, TearsLabel, DelayLabel, DelayMultLabel, RangeLabel, DamageLabel, DamageMultLabel, SpeedLabel, LastItemLabel ui.Label
var Data ItemData
var Run RunData
var LastItem LastItemStat

func main() {

	go ui.Do(gui)

	go processData()
	err := ui.Go()
	if err != nil {
		log.Fatal(err)
	}

}

func addToRunData(id string) {
	for _, element := range Data {
		if element.ID == id {
			LastItem.resetLastItem()
			if element.Tears != "" {
				tears, err := strconv.ParseFloat(element.Tears, 64)
				LastItem.Tears = tears
				if err != nil {
					log.Fatal(err)
				}
				Run.Tears += tears
			}
			if element.Delay != "" {
				delay, err := strconv.ParseFloat(element.Delay, 64)
				LastItem.Delay = delay
				if err != nil {
					log.Fatal(err)
				}
				Run.Delay += delay
			}
			if element.DelayMultiplier != "" {
				delayMultiplier, err := strconv.ParseFloat(element.DelayMultiplier, 64)
				LastItem.DelayMultiplier = delayMultiplier
				if err != nil {
					log.Fatal(err)
				}
				Run.DelayMultiplier += delayMultiplier
			}
			if element.Range != "" {
				rangey, err := strconv.ParseFloat(element.Range, 64)
				LastItem.Range = rangey
				if err != nil {
					log.Fatal(err)
				}
				Run.Range += rangey
			}
			if element.Height != "" {
				height, err := strconv.ParseFloat(element.Height, 64)
				LastItem.Height = height
				if err != nil {
					log.Fatal(err)
				}
				Run.Height += height
			}
			if element.Damage != "" {
				damage, err := strconv.ParseFloat(element.Damage, 64)
				LastItem.Damage = damage
				if err != nil {
					log.Fatal(err)
				}
				Run.Damage += damage
			}
			if element.DamageMultiplier != "" {
				damageMultiplier, err := strconv.ParseFloat(element.DamageMultiplier, 64)
				LastItem.DamageMultiplier = damageMultiplier
				if err != nil {
					log.Fatal(err)
				}
				Run.DamageMultiplier += damageMultiplier
			}
			if element.Speed != "" {
				speed, err := strconv.ParseFloat(element.Speed, 64)
				LastItem.Speed = speed
				if err != nil {
					log.Fatal(err)
				}
				Run.Speed += speed
			}
			LastItem.Name = element.Name
			break
		}
	}
}

func gui() {
	SeedLabel = ui.NewLabel("Seed")
	TearsLabel = ui.NewLabel("Tears")
	DelayLabel = ui.NewLabel("Delay")
	DelayMultLabel = ui.NewLabel("Delay Multiplier")
	RangeLabel = ui.NewLabel("Range")
	DamageLabel = ui.NewLabel("Damage")
	DamageMultLabel = ui.NewLabel("Damage Multiplier")
	SpeedLabel = ui.NewLabel("Speed")
	LastItemLabel = ui.NewLabel("Last Item Info")

	stack := ui.NewVerticalStack(
		SeedLabel,
		TearsLabel,
		DelayLabel,
		DelayMultLabel,
		RangeLabel,
		DamageLabel,
		DamageMultLabel,
		SpeedLabel,
		LastItemLabel)

	w := ui.NewWindow("Binding of Isaac Item Tracker", 1080, 350, stack)

	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})

	w.Show()
}

func processData() {
	ItemsJson, err := ioutil.ReadFile("items.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(ItemsJson, &Data); err != nil {
		log.Fatal(err)
	}

	logTail, err := tail.TailFile("/Users/danny/Library/Application Support/Binding of Isaac Rebirth/log.txt", tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range logTail.Lines {
		ui.Do(fixGui)
		if strings.HasPrefix(line.Text, "RNG Start Seed") {
			splitSeed := strings.Split(line.Text, " ")
			Run.Seed = strings.Join(splitSeed[3:5], " ")
			Run.Tears = 0
			Run.Delay = 0
			Run.DelayMultiplier = 0
			Run.Range = 0
			Run.Height = 0
			Run.Damage = 0
			Run.DamageMultiplier = 0
			Run.Speed = 0
			Run.ItemIDs = make([]string, 0)
		}

		if strings.HasPrefix(line.Text, "Adding collectible") {
			collectibleSplit := strings.Split(line.Text, " ")
			collectibleID := collectibleSplit[2]
			Run.ItemIDs = append(Run.ItemIDs, collectibleID)
			addToRunData(collectibleID)
		}
	}
}

func fixGui() {
	SeedLabel.SetText(Run.Seed)
	TearsLabel.SetText(fmt.Sprintf("Tear Rate: %.2f", Run.Tears))
	DelayLabel.SetText(fmt.Sprintf("Tear Delay: %.2f", Run.Delay))
	DelayMultLabel.SetText(fmt.Sprintf("Tear Delay Multiplier: %.2f", Run.DelayMultiplier))
	RangeLabel.SetText(fmt.Sprintf("Range: %.2f", Run.Range))
	DamageLabel.SetText(fmt.Sprintf("Damage: %.2f", Run.Damage))
	DamageMultLabel.SetText(fmt.Sprintf("Damage Multiplier: %.2f", Run.DamageMultiplier))
	SpeedLabel.SetText(fmt.Sprintf("Speed: %.2f", Run.Speed))
	LastItemLabel.SetText(fmt.Sprintf("Last Item was %s and gave you: %.2f Tear Rate, %.2f Tear Delay, %.2f Tear Delay Multiplier, %.2f Range,%.2f Height, %.2f Damage, %.2f Damage Multiplier, %.2f Speed",
		LastItem.Name,
		LastItem.Tears,
		LastItem.Delay,
		LastItem.DelayMultiplier,
		LastItem.Range,
		LastItem.Height,
		LastItem.Damage,
		LastItem.DamageMultiplier,
		LastItem.Speed),
	)
}

func (lastItem LastItemStat) resetLastItem() {
	lastItem.Tears = 0
	lastItem.Delay = 0
	lastItem.DelayMultiplier = 0
	lastItem.Range = 0
	lastItem.Height = 0
	lastItem.Damage = 0
	lastItem.DamageMultiplier = 0
	lastItem.Speed = 0
	lastItem.Description = ""
	lastItem.Name = ""
}
