package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, password := cloud.Login()
		password = cloud.HashPassword(password)

		viper.Set("username", username)
		viper.Set("password", password)

		saveConfig()

		err := pushlogin(username, password)
		if err != nil {
			log.Printf("error pushing login: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}

func pushlogin(us, ps string) error {
	data := map[string]string{
		"username": us,
		"password": ps,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json encoding error: %v", err)
	}

	req, err := http.NewRequest("POST", URL+"/adduser", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request error: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
