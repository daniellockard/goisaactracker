package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
var Window ui.Window

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
			LastItem = LastItemStat{}
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
			if element.Text != "" {
				LastItem.Description = element.Text 
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

	Window = ui.NewWindow("Binding of Isaac Item Tracker", 1200, 175, stack)

	Window.OnClosing(func() bool {
		ui.Stop()
		return true
	})

	Window.Show()
}

func processData() {
	ItemsJson, err := ioutil.ReadFile("items.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(ItemsJson, &Data); err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for _ = range ticker.C {
			lines, err := readLines()
			if err != nil {
				log.Fatal(err)
			}
			for _, element := range lines {
				processLine(element)
			}
			ui.Do(fixGui)
		}
	}()

}

func readLines() ([]string, error) {
	file, err := os.Open("log.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func fixGui() {
	SeedLabel.SetText(fmt.Sprintf("Seed: %s", Run.Seed))
	TearsLabel.SetText(fmt.Sprintf("Tear Rate: %.2f", Run.Tears))
	DelayLabel.SetText(fmt.Sprintf("Tear Delay: %.2f", Run.Delay))
	DelayMultLabel.SetText(fmt.Sprintf("Tear Delay Multiplier: %.2f", Run.DelayMultiplier))
	RangeLabel.SetText(fmt.Sprintf("Range: %.2f", Run.Range))
	DamageLabel.SetText(fmt.Sprintf("Damage: %.2f", Run.Damage))
	DamageMultLabel.SetText(fmt.Sprintf("Damage Multiplier: %.2f", Run.DamageMultiplier))
	SpeedLabel.SetText(fmt.Sprintf("Speed: %.2f", Run.Speed))
	LastItemLabel.SetText(buildItemDescription())
}

func processLine(line string) {
	if strings.HasPrefix(line, "RNG Start Seed") {
		splitSeed := strings.Split(line, " ")
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
		LastItem = LastItemStat{}
	}

	if strings.HasPrefix(line, "Adding collectible") {
		collectibleSplit := strings.Split(line, " ")
		collectibleID := collectibleSplit[2]
		Run.ItemIDs = append(Run.ItemIDs, collectibleID)
		addToRunData(collectibleID)
	}
}

func buildItemDescription() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("The last item was %s ", LastItem.Name))
	if LastItem.Description != "" {
		buffer.WriteString(fmt.Sprintf("(%s) ", LastItem.Description))
	}
	buffer.WriteString("and it gave you ")
	if LastItem.Tears != 0 {
		buffer.WriteString(fmt.Sprintf("Tears: %.2f, ", LastItem.Tears))
	}
	if LastItem.Delay != 0 {
		buffer.WriteString(fmt.Sprintf("Tear Delay: %.2f, ", LastItem.Delay))
	}
	if LastItem.DelayMultiplier != 0 {
		buffer.WriteString(fmt.Sprintf("Tear Delay Multiplier: %.2f, ", LastItem.DelayMultiplier))
	}
	if LastItem.Range != 0 {
		buffer.WriteString(fmt.Sprintf("Range: %.2f, ", LastItem.Range))
	}
	if LastItem.Height != 0 {
		buffer.WriteString(fmt.Sprintf("Height: %.2f, ", LastItem.Height))
	}
	if LastItem.Damage != 0 {
		buffer.WriteString(fmt.Sprintf("Damage: %.2f, ", LastItem.Damage))
	}
	if LastItem.DamageMultiplier != 0 {
		buffer.WriteString(fmt.Sprintf("Damage Multiplier: %.2f, ", LastItem.DamageMultiplier))
	}
	if LastItem.Speed != 0 {
		buffer.WriteString(fmt.Sprintf("Speed: %.2f, ", LastItem.Speed))
	}

	constructedString := buffer.String()
	if strings.HasSuffix(constructedString, ", ") {
		constructedString = constructedString[:len(constructedString)-2]
	}

	if strings.HasSuffix(constructedString, "and it gave you ") {
		constructedString = constructedString + "nothing"
	}
	return constructedString

}
