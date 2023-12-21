//go:build linux || freebsd

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/quexten/goldwarden/agent/systemauth/biometrics"
	"github.com/quexten/goldwarden/browserbiometrics"
	"github.com/spf13/cobra"
)

func setupPolkit() {
	file, err := os.Create("/tmp/goldwarden-policy")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(biometrics.POLICY)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

	command := exec.Command("pkexec", "mv", "/tmp/goldwarden-policy", "/usr/share/polkit-1/actions/com.quexten.goldwarden.policy")
	err = command.Run()
	if err != nil {
		panic(err)
	}

	command2 := exec.Command("pkexec", "chown", "root:root", "/usr/share/polkit-1/actions/com.quexten.goldwarden.policy")
	err = command2.Run()
	if err != nil {
		panic(err)
	}

	command3 := exec.Command("sudo", "chcon", "system_u:object_r:usr_t:s0", "/usr/share/polkit-1/actions/com.quexten.goldwarden.policy")
	err = command3.Run()
	if err != nil {
		fmt.Println("failed setting selinux context")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Set selinux context successfully")
		fmt.Println("Might require a reboot to take effect!")
	}

	fmt.Println("Polkit setup successfully")
}

func IsPolkitSetup() bool {
	_, err := os.Stat("/usr/share/polkit-1/actions/com.quexten.goldwarden.policy")
	return !os.IsNotExist(err)
}

var polkitCmd = &cobra.Command{
	Use:   "polkit",
	Short: "Sets up polkit",
	Long:  "Sets up polkit",
	Run: func(cmd *cobra.Command, args []string) {
		setupPolkit()
	},
}

const SYSTEMD_SERVICE = `[Unit]
Description="Goldwarden daemon"

[Service]
ExecStart=BINARY_PATH daemonize
Environment="DISPLAY=:0"

[Install]
WantedBy=graphical-session.target`

func setupSystemd() {
	file, err := os.Create("/tmp/goldwarden.service")
	if err != nil {
		panic(err)
	}

	path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	file.WriteString(strings.ReplaceAll(SYSTEMD_SERVICE, "BINARY_PATH", path))
	file.Close()

	userDirectory := os.Getenv("HOME")
	//ensure user systemd dir exists
	command0 := exec.Command("mkdir", "-p", userDirectory+"/.config/systemd/user/")
	err = command0.Run()
	if err != nil {
		fmt.Println("failed creating systemd user dir")
		fmt.Println(err.Error())
		panic(err)
	}

	command := exec.Command("mv", "/tmp/goldwarden.service", userDirectory+"/.config/systemd/user/goldwarden.service")
	err = command.Run()
	if err != nil {
		fmt.Println("failed moving goldwarden service file to systemd dir")
		fmt.Println(err.Error())
		panic(err)
	}

	command2 := exec.Command("systemctl", "--now", "--user", "enable", "goldwarden.service")
	command2.Stdout = os.Stdout
	command2.Stderr = os.Stderr
	err = command2.Run()
	if err != nil {
		fmt.Println("failed enabling systemd service")
		panic(err)
	}

	fmt.Println("Systemd setup successfully")
}

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Sets up systemd autostart",
	Long:  "Sets up systemd autostart",
	Run: func(cmd *cobra.Command, args []string) {
		setupSystemd()
	},
}

var browserbiometricsCmd = &cobra.Command{
	Use:   "browserbiometrics",
	Short: "Sets up browser biometrics",
	Long:  "Sets up browser biometrics",
	Run: func(cmd *cobra.Command, args []string) {
		err := browserbiometrics.DetectAndInstallBrowsers()
		if err != nil {
			fmt.Println("Error: " + err.Error())
		} else {
			fmt.Println("Done.")
		}
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up Goldwarden integrations",
	Long:  "Sets up Goldwarden integrations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.AddCommand(polkitCmd)
	setupCmd.AddCommand(systemdCmd)
	setupCmd.AddCommand(browserbiometricsCmd)
}
