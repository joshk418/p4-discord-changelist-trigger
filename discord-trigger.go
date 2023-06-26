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
	"log"
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

func loadConfiguration(file string) (Config, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return Config{}, nil
	}
	defer configFile.Close()

	var cfg Config
	if err := json.NewDecoder(configFile).Decode(&cfg); err != nil {
		return Config{}, err
	}
	
	return cfg, nil
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
		log.Printf("Non-204 response received: %d\n", resp.StatusCode)
		log.Println(string(body))
		return err
	}

	return nil
}

func getChangelistFiles(change int) (string, error) {
	cmd := exec.Command("p4", "-Ztag", "-F", "%depotFile% - %action%", "files", fmt.Sprintf("@=%d", change))
	
	cmdOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(cmdOutput), nil
}

func getChangelistDescription(change int) (string, error) {
	cmd := exec.Command("p4", "-Ztag", "-F", "%Description%", "change", "-o", fmt.Sprintf("%d", change))
	
	cmdOutput, err := cmd.Output()
	if err != nil {
		return "", nil
	}

	return strings.Trim(string(cmdOutput), "\n"), nil
}

func setupLogFile() error {
	f, err := os.OpenFile("discord-output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Initialized log file")

	return nil
}

func main() {
	username := flag.String("user", "qbag", "username of the change maker")
	changeNumber := flag.Int("change", -1, "changelist number")
	configFileLocation := flag.String("config", "/etc/discord-trigger.conf", "config file location")
	flag.Parse()

	if err := setupLogFile(); err != nil {
		fmt.Println(err)
		return
	}

	config, err := loadConfiguration(*configFileLocation)
	if err != nil {
		log.Println(err.Error())
		return
	}

	changedFiles, err := getChangelistFiles(*changeNumber)
	if err != nil {
		log.Printf("Failed to get files for changelist %d\n", *changeNumber)
		log.Println(err.Error())
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
		log.Println("Failed to POST webhook")
		log.Println(err.Error())
		return
	}
}
