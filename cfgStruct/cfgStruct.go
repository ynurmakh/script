package cfgstruct

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	PassOfScript         string   `json: "passOfScript"`
	PassOfPC             string   `json: "passOfPC"`
	CastomLockScreenPath string   `json: "castomLockScreenPath"`
	GitLinks             []string `json: "gitLinks"`
	WifiHotspotName      string   `json: "wifiHotspotName"`
	WifiHotspotPass      string   `json: "wifiHotspotPass"`
	GitGlobalUser_NAME   string   `json: "gitGlobalUser_NAME"`
	GitGlobalUser_EMAIL  string   `json: "gitGlobalUser_EMAIL"`
}

func GetConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("not found config.json")
	}
	defer file.Close()

	var res Config

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		log.Fatalln("error on decode config file")
	}

	return res
}
