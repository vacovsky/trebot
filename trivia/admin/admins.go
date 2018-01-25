package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AdminModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LoadAdmins loads the contents of the admins file into memory
func LoadAdmins(adminsFile string) map[string]AdminModel {
	adminsLocal := []AdminModel{}
	rAdmins := map[string]AdminModel{}
	file, err := os.Open(adminsFile)
	defer file.Close()
	tried := false
	if err != nil && !tried {
		tried = true
		SaveAdmins(adminsFile, rAdmins)
		LoadAdmins(adminsFile)
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&adminsLocal); err != nil {
		fmt.Println(err)
	}
	for _, user := range adminsLocal {
		tmp := AdminModel{
			Name: user.Name,
			ID:   user.ID,
		}
		rAdmins[user.ID] = tmp
	}
	return rAdmins
}

// SaveAdmins writes new bans to the bans file
func SaveAdmins(adminsFile string, admins map[string]AdminModel) {
	adminsModel := []AdminModel{}
	for _, u := range admins {
		adminsModel = append(adminsModel, u)
	}
	adminsJSON, _ := json.MarshalIndent(adminsModel, "", "    ")
	err := ioutil.WriteFile(adminsFile, adminsJSON, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
