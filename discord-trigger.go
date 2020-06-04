package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type DiscordEmbedBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DiscordWebhookBody struct {
	Embeds []DiscordEmbedBody `json:"embeds"`
}

type Config struct {
	Webhook string `json:"webhook"`
}

func loadConfiguration(file string) (config Config, err error) {
	configFile, err := os.Open(file)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	return
}

func sendDiscordWebhook(url string, embed []DiscordEmbedBody) error {
	data, err := json.Marshal(DiscordWebhookBody{Embeds: embed})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		fmt.Printf("Non-204 response received: %d\n", resp.StatusCode)
		fmt.Println(string(body))
		return err
	}

	return nil
}

func getChangelistFiles(change int) (changedFiles string, err error) {
	cmd := exec.Command("p4", "-Ztag", "-F", "%depotFile% - %action%", "files", fmt.Sprintf("@=%d", change))
	cmdOutput, err := cmd.Output()

	if err != nil {
		return
	}

	changedFiles = string(cmdOutput)
	return
}

func getChangelistDescription(change int) (changelistDescription string, err error) {
	cmd := exec.Command("p4", "-Ztag", "-F", "%Description%", "change", "-o", fmt.Sprintf("%d", change))
	cmdOutput, err := cmd.Output()

	if err != nil {
		return
	}

	changelistDescription = strings.Trim(string(cmdOutput), "\n")
	return
}

func main() {
	username := flag.String("user", "qbag", "username of the change maker")
	changeNumber := flag.Int("change", -1, "changelist number")
	configFileLocation := flag.String("config", "/etc/discord-trigger.conf", "config file location")
	flag.Parse()

	config, err := loadConfiguration(*configFileLocation)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	changedFiles, err := getChangelistFiles(*changeNumber)
	if err != nil {
		fmt.Printf("Failed to get files for changelist %d\n", *changeNumber)
		fmt.Println(err.Error())
		return
	}

	changelistDescription, err := getChangelistDescription(*changeNumber)
	if err != nil {
		fmt.Printf("Failed to get changed files for change %d\n", *changeNumber)
		fmt.Println(err.Error())
		return
	}

	embedTitle := fmt.Sprintf("New change %d submitted by %s: %s", *changeNumber, *username, changelistDescription)
	embedBody := []DiscordEmbedBody{DiscordEmbedBody{Title: embedTitle, Description: changedFiles}}

	err = sendDiscordWebhook(config.Webhook, embedBody)
	if err != nil {
		fmt.Println("Failed to POST webhook")
		fmt.Println(err.Error())
		return
	}
}
