package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

func getEmail(c string) (string, string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://wallet.dosi.world/api/v1/user", nil)
	cookie := &http.Cookie{
		Name:  "DOSI_SES",
		Value: c,
	}
	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return "", ""
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	obj := &Resp{}
	if err := json.Unmarshal(body, &obj); err != nil {
		fmt.Println("JSON Error")
	}

	return obj.RespCode, obj.RespData.Email
}

func getNftCount(c string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://citizen.dosi.world/api/citizen/v1/membership", nil)
	cookie := &http.Cookie{
		Name:  "DOSI_SES",
		Value: c,
	}
	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	obj := &Resp{}
	if err := json.Unmarshal(body, &obj); err != nil {
		fmt.Println("JSON Error")
	}

	return strconv.Itoa(obj.NftCount)

}

func sendTele(t string, i string, txt string) bool {
	client := &http.Client{}
	data := url.Values{}
	data.Set("chat_id", i)
	data.Set("text", txt)
	data.Set("parse_mode", "Markdown")
	req, _ := http.NewRequest("POST", "https://api.telegram.org/bot"+t+"/sendMessage", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return false
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	obj := &Resp{}
	if err := json.Unmarshal(body, &obj); err != nil {
		fmt.Println("JSON Error")
	}
	return obj.Ok
}

var (
	ConfigFile string
	Tg         string
	config     Elements
	outputFile string
)

func init() {
	flag.StringVar(&ConfigFile, "c", "config.yaml", "Configuration File")
	flag.StringVar(&ConfigFile, "conf", "config.yaml", "Configuration File")
	flag.StringVar(&Tg, "tg", "", "Sending Notification to Telegram")
	flag.StringVar(&Tg, "telegram", "", "Sending Notification to Telegram")
	flag.StringVar(&outputFile, "o", "", "File Output")
	flag.StringVar(&outputFile, "output", "", "File Output")
	flag.Usage = func() {
		h := []string{
			"",
			"Dosee (Dosi Checker)",
			"",
			"Is a simple GO tool for Check Dosi Account by Session Lists",
			"",
			"Coded By : github.com/vsec7",
			"",
			"Basic Usage :",
			" ▶ cat session_list.txt | dosee",
			" ▶ dosee < session_list.txt",
			"Advanced Usage :",
			" ▶ cat session_list.txt | dosee -tg all -o result.txt",
			"",
			"Options :",
			"  -c, --conf <config.yaml>		Set file config.yaml (default: config.yaml)",
			"  -tg, --telegram <all|active|expired>	Set Notification to Telegram",
			"  -o, --output <file>	        	Set Output File",
			"",
			"",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
	flag.Parse()
}

func main() {
	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			for s := range jobs {
				stat, e := getEmail(s)
				if stat == "200" {
					nft := getNftCount(s)
					fmt.Printf("Email: %s | NFT: %s\n", e, nft)

					if len(Tg) != 0 {
						yamlFile, err := ioutil.ReadFile(ConfigFile)
						if err != nil {
							fmt.Printf("[ERROR] File %s not found!\n", ConfigFile)
							os.Exit(0)
						}
						err = yaml.Unmarshal(yamlFile, &config)

						if (Tg == "all") || (Tg == "active") {
							sendTele(config.Token, config.Chat_id, "✅ Email: `"+e+"` | NFT: "+nft)
						}
					}

					if outputFile != "" {
						file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Printf("[!] Failed Creating File: %s", err)
						}
						buf := bufio.NewWriter(file)
						buf.WriteString("[LIVE] Email: " + e + " | NFT: " + nft + "\n")
						buf.Flush()
						file.Close()
					}
				} else {
					fmt.Printf("SESSION EXPIRED: %s\n", s)

					if len(Tg) != 0 {
						yamlFile, err := ioutil.ReadFile(ConfigFile)
						if err != nil {
							fmt.Printf("[ERROR] File %s not found!\n", ConfigFile)
							os.Exit(0)
						}
						err = yaml.Unmarshal(yamlFile, &config)
						if (Tg == "all") || (Tg == "expired") {
							sendTele(config.Token, config.Chat_id, "🔴 SESSION EXPIRED: `"+s+"`")
						}
					}

					if outputFile != "" {
						file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Printf("[!] Failed Creating File: %s", err)
						}
						buf := bufio.NewWriter(file)
						buf.WriteString("[SESSION EXPIRED]: " + s + "\n")
						buf.Flush()
						file.Close()
					}
				}
			}
			wg.Done()
		}()
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)
	wg.Wait()
}

type Resp struct {
	RespCode string   `json:"responseCode"`
	RespData Elements `json:"responseData"`
	NftCount int      `json:"nftCount"`
	Ok       bool     `json:"ok"`
}

type Elements struct {
	Email   string `json:"email"`
	Token   string `yaml:"BOT_TOKEN"`
	Chat_id string `yaml:"CHAT_ID"`
}
