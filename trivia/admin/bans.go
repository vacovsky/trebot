package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// BanModel stores the bans information
type BanModel struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	BannedSince  time.Time `json:"banned_since"`
	Reason       string    `json:"reason"`
	BannedByID   string    `json:"banned_by_id"`
	BannedByName string    `json:"banned_by_name"`
}

// LoadBans loads the contents of the Bans file into memory
func LoadBans(bansFile string) map[string]BanModel {
	bannedLocal := []BanModel{}
	rBans := map[string]BanModel{}
	file, err := os.Open(bansFile)
	defer file.Close()
	tried := false
	if err != nil && !tried {
		tried = true
		SaveBans(bansFile, rBans)
		LoadBans(bansFile)
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&bannedLocal); err != nil {
		fmt.Println(err)
	}
	for _, user := range bannedLocal {
		tmp := BanModel{
			Name: user.Name,
			ID:   user.ID,
		}
		rBans[user.ID] = tmp
	}
	return rBans
}

// SaveBans writes new bans to the bans file
func SaveBans(bansFile string, bans map[string]BanModel) {
	bansModel := []BanModel{}
	for _, u := range bans {
		bansModel = append(bansModel, u)
	}
	bansJSON, _ := json.MarshalIndent(bansModel, "", "    ")
	err := ioutil.WriteFile(bansFile, bansJSON, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
