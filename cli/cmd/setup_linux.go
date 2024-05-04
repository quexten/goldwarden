//go:build linux || freebsd

package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/quexten/goldwarden/cli/agent/systemauth/biometrics"
	"github.com/quexten/goldwarden/cli/browserbiometrics"
	"github.com/spf13/cobra"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func setupPolkit() {
	if isRoot() {
		fmt.Println("Do not run this command as root!")
		return
	}

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

//go:embed goldwarden.service
var systemdService string

func setupSystemd() {
	if isRoot() {
		fmt.Println("Do not run this command as root!")
		return
	}

	file, err := os.Create("/tmp/goldwarden.service")
	if err != nil {
		panic(err)
	}

	path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(strings.ReplaceAll(systemdService, "@BINARY_PATH@", path))
	if err != nil {
		panic(err)
	}
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
		if isRoot() {
			fmt.Println("Do not run this command as root!")
			return
		}

		setupSystemd()
	},
}

var browserbiometricsCmd = &cobra.Command{
	Use:   "browserbiometrics",
	Short: "Sets up browser biometrics",
	Long:  "Sets up browser biometrics",
	Run: func(cmd *cobra.Command, args []string) {
		if isRoot() {
			fmt.Println("Do not run this command as root!")
			return
		}

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
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.AddCommand(polkitCmd)
	setupCmd.AddCommand(systemdCmd)
	setupCmd.AddCommand(browserbiometricsCmd)
}
