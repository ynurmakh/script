package everytime

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	. "github.com/p0n41k/goUseful"

	. "reg/cfgStruct"
)

func EveryTime(pathOfMain string) {
	cfg := GetConfig()

	log.Println("START: WIFI")
	go wifiHotspot(cfg)

	log.Println("START: git username and email write")
	go setGitConfigs(cfg)

	log.Println("START: hotkey for block create")
	go hotkeyCreate(cfg, pathOfMain)

	log.Println("START: user repos cloning")
	go cloneUsersRepo(cfg, pathOfMain)

	log.Println("START: russian keyboard add")
	exec.Command("gsettings", "set", "org.gnome.desktop.input-sources", "sources", "[('xkb', 'us'), ('xkb', 'ru')]").Run()

	v1, v2, v3, v4, v5 := 0, 0, 0, 0, 0
	for true {

		if runtime.NumGoroutine() == 1 {
			break
		}

		v2, v3, v4, v5 = v1, v2, v3, v4
		v1 = runtime.NumGoroutine()

		if v1 == v2 && v2 == v3 && v3 == v4 && v4 == v5 {
			break
		}

		time.Sleep(1 * time.Second)
	}
	setPassPC(cfg)

	fmt.Println()
	fmt.Println("SUCESSFUL")
	fmt.Println()
}

func hotkeyCreate(cfg Config, mainpath string) error {
	var err error

	cmd := exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "name", "'Custom_Lock'")
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "command", "'./block.sh'")
	cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "binding", "'Pause'")
	cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func wifiHotspot(cfg Config) {
	if cfg.WifiHotspotName == "" || len(cfg.WifiHotspotPass) < 8 {
		log.Println("Wi-fi hotspot not started. Not correct config")
	}

	wifiCmd := exec.Command("nmcli", "device", "wifi", "hotspot", "ifname", "wlp0s20f3", "ssid", cfg.WifiHotspotName, "password", cfg.WifiHotspotPass)
	err := wifiCmd.Run()
	if err != nil {
		log.Println("Wi-fi hotspot not started. Not correct config")
	}
}

func setGitConfigs(cfg Config) {
	gitConfigCmd := exec.Command("git", "config", "--global", "user.name", cfg.GitGlobalUser_NAME)
	gitConfigCmd.Run()

	gitConfigCmd = exec.Command("git", "config", "--global", "user.email", cfg.GitGlobalUser_EMAIL)
	gitConfigCmd.Run()
}

func cloneUsersRepo(cfg Config, pathOfMain string) {
	os.Chdir("/home/student/")

	for i := 0; i < len(cfg.GitLinks); i++ {
		log.Println("All repos:", len(cfg.GitLinks), ": Cloned:", i+1)

		_, err := Bash("git clone " + cfg.GitLinks[i])
		if err != nil {
			log.Println("Link", i+1, "not cloned"+err.Error())
		}

	}

	os.Chdir(pathOfMain)
}

func setPassPC(cfg Config) {
	openTerminal := func() {
		Bash("gnome-terminal")
	}
	go openTerminal()
	time.Sleep(2 * time.Second)
	Bash([]string{
		"xdotool mousemove 2000 700 click 1",
		"xdotool type \"passwd\"",
		"xdotool key Return",
		"xdotool type \"" + cfg.PassOfPC + "\"",
		"sleep 0.1",
		"xdotool key Return",
		"sleep 0.1",
		"xdotool type \"" + cfg.PassOfPC + "\"",
		"sleep 0.1",
		"xdotool key Return",
		"sleep 0.1",
		"xdotool type \"exit\"",
		"sleep 0.1",
		"xdotool key Return",
	})
}
