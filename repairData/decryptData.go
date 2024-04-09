package repairdata

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	. "reg/repairData/everyTime"

	. "github.com/p0n41k/goUseful"
)

func RepairSettings(args []string, pathOfMain string) error {
	var err error

	files, err := findCryptedFilesFromPATH(pathOfMain)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		log.Println("cryptedDATA finded")
		err = dataDecrypt(args, pathOfMain, files)
		if err != nil {
			return err
		}

		err = mvDecryptedDATA()
		if err != nil {
			return err
		}

		EveryTime(pathOfMain)

	} else {
		log.Print("WARNING\nWARNING User`s DATA not found in working repository\nWARNING If it your first start WRITE confog.json right data and contunue\nWARNING Are you edit comfig.json? (y/n) ")
		a := InputScaner()
		if a == "y" {
			EveryTime(pathOfMain)
		} else {
			fmt.Println("Exited")
			return err
		}
	}

	return err
}

func findCryptedFilesFromPATH(PATH string) ([]string, error) {
	filesFS, err := os.ReadDir(PATH)
	if err != nil {
		return make([]string, 0), err
	}

	files := make([]string, 0, len(filesFS))
	for i := 0; i < len(filesFS); i++ {
		if !filesFS[i].IsDir() {
			files = append(files, filesFS[i].Name())
		}
	}

	cryptedDATA := make([]string, 0)
	for i := 0; i < len(files); i++ {
		if len(files[i]) > 11 && files[i][:11] == "cryptedDATA" {
			cryptedDATA = append(cryptedDATA, files[i])
		}
	}

	return cryptedDATA, nil
}

func dataDecrypt(args []string, pathOfMain string, cryptedFiles []string) error {
	var err error

	pass := ""
	if len(args) > 1 {
		pass = args[1]
	} else {
		log.Print("Enter your PASSWORD: ")
		log.Print()
		pass = AnonimInput()
	}

	if !passChecer(pass) {
		return errors.New("Password not corect")
	}

	passForDecript := sha256.Sum256([]byte(pass))
	passSTR := hex.EncodeToString(passForDecript[:])

	err = sobiratelAllInOne(cryptedFiles)

	err = decryptCryptedFile(passSTR)
	return err
}

func passChecer(pass string) bool {
	log.Println("Checking password")
	time.Sleep(2 * time.Second)

	upper, sim, passLen := false, false, false
	for i := 0; i < len(pass); i++ {
		if !upper && (pass[i] >= 'A' && pass[i] <= 'Z') {
			upper = true
		}

		if !sim && (pass[i] < 'a' || pass[i] > 'z') && (pass[i] < 'A' || pass[i] > 'Z') {
			sim = true
		}

		if !passLen && len(pass) >= 8 {
			passLen = true
		}

	}

	if !upper || !sim || !passLen {
		return false
	} else {
		return true
	}
}

func sobiratelAllInOne(cryptedFiles []string) error {
	var err error

	allCryptedDATA := make([]byte, 0)
	for i := 0; i < len(cryptedFiles); i++ {
		file, err := os.ReadFile(cryptedFiles[i])
		if err != nil {
			return err
		}

		allCryptedDATA = append(allCryptedDATA, file...)
	}

	err = os.WriteFile("cryptedDATA.tar.gz.enc", allCryptedDATA, 0777)

	return err
}

func decryptCryptedFile(passSTR string) error {
	var err error

	_, err = Bash("openssl enc -d -aes256 -pbkdf2 -in cryptedDATA.tar.gz.enc -out uncrupted.tar.gz -k " + passSTR)
	if err != nil {
		log.Println("Password not corect")
		return err
	}

	_, err = Bash("gunzip uncrupted.tar.gz")
	if err != nil {
		log.Println("Somthing wrong with gunzip")
		return err
	}

	_, err = Bash("tar -xf uncrupted.tar")
	if err != nil {
		log.Println("Somthing wrong with tar unarchiver")
		return err
	}

	_, err = Bash("rm -f uncrupted.tar")
	if err != nil {
		log.Println("Somthing wrong with deliting uncrupted.tar")
		return err
	}

	return err
}

func mvDecryptedDATA() error {
	var err error

	err = sshDirMove()

	err = browserDirMove()

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

func closeAllFirefox() {
	if windows, err := Bash("xdotool search --name firefox"); err == nil {

		windows = windows[:len(windows)-1]
		windowArr := strings.Split(windows, "\n")

		for len(windowArr) != 0 {
			fmt.Println(Bash("xdotool windowkill " + windowArr[0]))
			if windows, err := Bash("xdotool search --name firefox"); err == nil {
				windows = windows[:len(windows)-1]
				windowArr = strings.Split(windows, "\n")
			} else {
				break
			}

		}
	}
}

func openFirefox() {
	Bash("xdotool mousemove 50 50 click 1")
	time.Sleep(2 * time.Second)
}

func sshDirMove() error {
	_, err := Bash([]string{
		"rm -r -f /home/student/.ssh",
		"cp -r .ssh /home/student/.",
		"ssh-keyscan 01.alem.school >> ~/.ssh/known_hosts",
		"ssh-keyscan github.com >> ~/.ssh/known_hosts",
		"ssh-agent",
		"rm -r -f .ssh",
	})
	return err
}

func browserDirMove() error {
	var err error
	usersRepo := browserDirParse()

	if usersRepo == "" {
		closeAllFirefox()
		openFirefox()
		closeAllFirefox()

		usersRepo = browserDirParse()
		if usersRepo == "" {
			return errors.New("Script can`t parse browser info. Reboot PC and try again, or contact with author")
		}

	}

	_, err = Bash("rm -r -f /home/student/.mozilla/firefox/" + usersRepo)
	if err != nil {
		return err
	}

	_, err = Bash("cp -r browserDATA /home/student/.mozilla/firefox/.")
	if err != nil {
		return err
	}

	_, err = Bash("mv /home/student/.mozilla/firefox/browserDATA /home/student/.mozilla/firefox/" + usersRepo)
	if err != nil {
		return err
	}

	err = exec.Command("rm", "-r", "-f", "browserDATA").Run()
	if err != nil {
		return err
	}
	return nil
}
