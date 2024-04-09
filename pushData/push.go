package pushdata

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	. "reg/cfgStruct"

	. "github.com/p0n41k/goUseful"
)

func PushAndBlock(cfg Config) {
	cfg = GetConfig()

	err := moveBrowserDirToMain()
	if err != nil {
		errHandler(errors.New("Error: func pushAndBlock(\n" + err.Error() + " \n)"))
		return
	}

	err = creatCryptedDATA(cfg)
	if err != nil {
		errHandler(errors.New("Error: func pushAndBlock(\n" + err.Error() + " \n)"))
		return
	}

	err = splitDATA()
	if err != nil {
		errHandler(errors.New("Cant`t crypt DATA" + err.Error()))
	}

	_, err = Bash("rm -f config.json")

	err = pushDATA()
	if err != nil {
		errHandler(errors.New("Cant`t crypt DATA" + err.Error()))
	}
}

func moveBrowserDirToMain() error {
	var err error
	log.Println("moveBrowserDirToMain() start")

	defaultDirName := browserDirParse()

	if defaultDirName == "" {
		return errors.New("Error: can`t parse dir name. func moveBrowserDirToMain(\n" + err.Error() + " \n)")
	}
	log.Println("dir name", defaultDirName)

	_, err = Bash("cp -r /home/student/.mozilla/firefox/" + defaultDirName + " .")
	if err != nil {
		return errors.New("Error: can`t copy Browser dir to work dir. func moveBrowserDirToMain(\n" + err.Error() + " \n)")
	}

	log.Println("dir", defaultDirName, "moved to main")

	_, err = Bash("mv " + defaultDirName + " browserDATA")
	if err != nil {
		return errors.New("Error: can`t rename Browser dir to browseDATA. func moveBrowserDirToMain(\n" + err.Error() + " \n)")
	}
	log.Println("dir", defaultDirName, "renamed to browserDATA")

	_, err = Bash("cp -r /home/student/.ssh .")
	if err != nil {
		return errors.New("Error: can`t copy .ssh to working dir. func moveBrowserDirToMain(\n" + err.Error() + " \n)")
	}
	log.Println("dir", defaultDirName, "renamed to browserDATA")

	return err
}

func browserDirParse() string {
	browserUserInfoFile := "/home/student/.mozilla/firefox/profiles.ini"
	infoBytes, err := os.ReadFile(browserUserInfoFile)
	if err != nil {
		return ""
	}

	infoSting := string(infoBytes)
	index := strings.Index(infoSting, "Default=")
	if index == -1 {
		return ""
	}
	def := infoSting[index:]
	defSplited := strings.Split(def, "\n")
	if len(defSplited) == 0 {
		return ""
	}

	return defSplited[0][8:]
}

func creatCryptedDATA(cfg Config) error {
	log.Println("creatCryptedDATA() start")

	filesForCrypt := []string{
		".ssh",
		"browserDATA",
		"config.json",
	}

	tarCmd := exec.Command("tar", "czf", "-", filesForCrypt[0], filesForCrypt[1], filesForCrypt[2])
	passBytes := sha256.Sum256([]byte(cfg.PassOfScript))
	pass := hex.EncodeToString(passBytes[:])
	cryptedFinalName := "cryptedDATA.tar.gz.enc"

	opensslCmd := exec.Command("openssl", "enc", "-e", "-aes256", "-pbkdf2", "-out", cryptedFinalName, "-k", pass)

	// Захватываем вывод tar и устанавливаем его ввод для openssl
	tarOutput, err := tarCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	opensslCmd.Stdin = tarOutput

	// Запускаем tar и openssl
	if err := tarCmd.Start(); err != nil {
		panic(err)
	}
	if err := opensslCmd.Start(); err != nil {
		panic(err)
	}

	// Ждем завершения обеих команд
	if err := tarCmd.Wait(); err != nil {
		panic(err)
	}
	if err := opensslCmd.Wait(); err != nil {
		panic(err)
	}

	return err
}

func errHandler(err error) {
	errStr := err.Error()
	os.WriteFile("/home/student/.err.log", []byte(errStr+"\n\n    Contacts\nTG: @nur_erbol\nEmail: nur_erbol_2002@mail.ru\nGithub: p0n41k"), 0777)
	Bash("bash .errHandler.sh")
}

func splitDATA() error {
	_, err := Bash("split -d -b 50M cryptedDATA.tar.gz.enc cryptedDATA")
	if err != nil {
		return errors.New("Error: can`t split cryptedDATA func splitDATA(\n" + err.Error() + " \n)")
	}

	_, err = Bash("rm -r -f cryptedDATA.tar.gz.enc")
	if err != nil {
		return errors.New("Error: remove cryptedDATA func splitDATA(\n" + err.Error() + " \n)")
	}
	return err
}

func pushDATA() error {
	Bash("rm -r -f .ssh")

	_, err := Bash("git add .")
	if err != nil {
		return errors.New("Error: git add . func pushDATA(\n" + err.Error() + " \n)")
	}

	_, err = Bash("git commit -m \"AutoSave\"")
	if err != nil {
		return errors.New("Error: git commit -m func pushDATA(\n" + err.Error() + " \n)")
	}

	_, err = Bash("git push --force")
	if err != nil {
		return errors.New("Error: git push --force func pushDATA(\n" + err.Error() + " \n)")
	}
	return err
}
