package updateChecker

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	gouseful "github.com/p0n41k/goUseful"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

func UpdateCheck(scriptName, scriptVersoin, pathOfMain string) error {
	reposURLAPI := "https://api.github.com/repos/ynurmakh/FastReg/releases/latest"
	reposURLCloning := "https://github.com/ynurmakh/FastReg"

	version, news, err := githubRepoParse(reposURLAPI)
	if err != nil {
		return errors.New("Somthing wrong with get request in updater func" + err.Error())
	}

	if version == scriptVersoin {
		log.Println("You have The latest version of The sctipt. Continue...")
		return nil
	} else {
		log.Printf("Finded new version of script.\n\nNews:\n%v\n\n\n%v >>> %v\nPress enter for download. ", news, scriptVersoin, version)
		gouseful.InputScaner()
		err := downloadInstall(reposURLCloning)
		if err != nil {
			return errors.New("Can`t download new version. You can download there: " + reposURLCloning)
		}

		return errors.New("")
	}
}

func githubRepoParse(reposURL string) (string, string, error) {
	resp, err := http.Get(reposURL)
	if err != nil {
		return "", "", errors.New("Error with get request\n" + err.Error())
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return "", "", errors.New("The script can`t update. Get requests status not 200\n" + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", errors.New("Error with get requests budy\n" + err.Error())
	}

	var releaseInfo ReleaseInfo
	if err := json.Unmarshal(body, &releaseInfo); err != nil {
		return "", "", errors.New("Error with unmarshall\n" + err.Error())
	}

	return releaseInfo.TagName, releaseInfo.Body, nil
}

func StartTheScript() {
	cmd := exec.Command("chmod", "+x", "*")
	_, _ = cmd.CombinedOutput()
	cmd = exec.Command("./reg")
	_, _ = cmd.CombinedOutput()
}

func downloadInstall(url string) error {
	cmd, err := gouseful.Bash([]string{
		"rm -r -f ./FastReg/",
		"git clone " + url,
	})
	if err != nil {
		return err
	}
	log.Println(cmd)

	cmd, err = gouseful.Bash([]string{
		"mv ./FastReg/reg .",
		"mv ./FastReg/README.md .",
		"rm -r -f ./FastReg",
	})
	if err != nil {
		return err
	}
	log.Println(cmd)

	return err
}
