package main

import (
	"log"
	"os"
	"path/filepath"
	"reg/updateChecker"
	"strings"

	. "github.com/p0n41k/goUseful"

	. "reg/cfgStruct"
	. "reg/gitConfigs"
	. "reg/pushData"
	. "reg/repairData"
)

var (
	scriptName    string
	scriptVersoin string
)

func init() {
	scriptName = getNameOfScrypt()
	scriptVersoin = "1.0"
}

func getNameOfScrypt() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	return filepath.Base(exePath)
}

func main() {
	if err := fromMain(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func fromMain(args []string) error {
	var err error

	// Печатаем версию скрипта
	log.Println("script v" + scriptVersoin)

	// Парсим путь до скрипта и переходим по этой папке
	pathOfMain := GetPathOfCallerFile(args[0])
	if err := os.Chdir(pathOfMain); err != nil {
		return err
	}

	// Проверяем обновления
	err = updateChecker.UpdateCheck(scriptName, scriptVersoin, pathOfMain)
	if err != nil {
		return err
	}

	os.Exit(1)

	GitConfigs()
	blockScriptCreate(pathOfMain)
	createErrHandler()
	createGitIgnore()
	createMainPathTXT(pathOfMain)

	if len(args) > 1 && args[1] == "--push-block" {
		cfg := GetConfig()
		PushAndBlock(cfg)
		return nil
	} else if (len(args) == 1) || (len(args) > 1 && args[1] != "--push-block") {
		err = RepairSettings(args, pathOfMain)
		if err != nil {
			return err
		}
	}

	return err
}

func blockScriptCreate(mainpath string) {
	os.WriteFile("/home/student/block.sh", []byte("# !/bin/bash\n\npath=$(cat .pathOfMain.txt)\npathToPNG=$(cat $path/config.json | grep Path | cut -d '\"' -f 4)\ni3lock -i $pathToPNG\n# $path'reg' --push-block"), 0777)
}

func createErrHandler() {
	errHandler := "# !/bin/bash\n\nerr=$(cat ~/.err.log)\ngnome-terminal.real -t \"Error\" --geometry=100x10 --wait -- bash -c \"echo '$err'; read test; echo \\$test > /tmp/phrase:\""
	os.WriteFile(".errHandler.sh", []byte(errHandler), 0777)
}

func createGitIgnore() {
	ignorList := []string{
		"uncrypted*",
		"*default-release*",
		".errHandler.sh",
		"browserDATA",
		".ssh",
	}

	ignorStr := strings.Join(ignorList, "\n")
	os.WriteFile(".gitignore", []byte(ignorStr), 0777)
}

func createMainPathTXT(pathOfMain string) {
	os.WriteFile("/home/student/.pathOfMain.txt", []byte(pathOfMain), 0777)
}
