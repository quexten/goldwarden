package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/icza/gox/stringsx"
	"github.com/quexten/goldwarden/client"
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/spf13/cobra"
)

var baseLoginCmd = &cobra.Command{
	Use:   "logins",
	Short: "Commands for managing logins.",
	Long:  `Commands for managing logins.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var getLoginCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a login in your vault",
	Long:  `Gets a login in your vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		loginIfRequired()

		uuid, _ := cmd.Flags().GetString("uuid")
		name, _ := cmd.Flags().GetString("name")
		username, _ := cmd.Flags().GetString("username")
		fullOutput, _ := cmd.Flags().GetBool("full")

		resp, err := commandClient.SendToAgent(messages.GetLoginRequest{
			Name:     name,
			Username: username,
			UUID:     uuid,
		})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch resp.(type) {
		case messages.GetLoginResponse:
			response := resp.(messages.GetLoginResponse)
			if fullOutput {
				fmt.Println(response.Result)
			} else {
				fmt.Println(response.Result.Password)
			}
			break
		case messages.ActionResponse:
			fmt.Println("Error: " + resp.(messages.ActionResponse).Message)
			return
		}
	},
}

var listLoginsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all logins in your vault",
	Long:  `Lists all logins in your vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		loginIfRequired()

		logins, err := ListLogins(commandClient)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		var toPrintLogins []map[string]string
		for _, login := range logins {
			data := map[string]string{
				"name":     stringsx.Clean(login.Name),
				"uuid":     stringsx.Clean(login.UUID),
				"username": stringsx.Clean(login.Username),
				"password": stringsx.Clean(strings.ReplaceAll(login.Password, "\"", "\\\"")),
				"totp":     stringsx.Clean(login.TOTPSeed),
				"uri":      stringsx.Clean(login.URI),
			}
			toPrintLogins = append(toPrintLogins, data)
		}
		toPrintJSON, _ := json.Marshal(toPrintLogins)
		fmt.Println(string(toPrintJSON))
	},
}

func ListLogins(client client.Client) ([]messages.DecryptedLoginCipher, error) {
	resp, err := client.SendToAgent(messages.ListLoginsRequest{})
	if err != nil {
		return []messages.DecryptedLoginCipher{}, err
	}

	switch resp.(type) {
	case messages.GetLoginsResponse:
		castedResponse := (resp.(messages.GetLoginsResponse))
		return castedResponse.Result, nil
	case messages.ActionResponse:
		castedResponse := (resp.(messages.ActionResponse))
		return []messages.DecryptedLoginCipher{}, errors.New("Error: " + castedResponse.Message)
	default:
		return []messages.DecryptedLoginCipher{}, errors.New("Wrong response type")
	}
}

func init() {
	rootCmd.AddCommand(baseLoginCmd)
	baseLoginCmd.AddCommand(getLoginCmd)
	getLoginCmd.PersistentFlags().String("name", "", "")
	getLoginCmd.PersistentFlags().String("username", "", "")
	getLoginCmd.PersistentFlags().String("uuid", "", "")
	getLoginCmd.PersistentFlags().Bool("full", false, "")
	baseLoginCmd.AddCommand(listLoginsCmd)
}
