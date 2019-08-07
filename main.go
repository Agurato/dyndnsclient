package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&nested.Formatter{
		ShowFullLevel: true,
		TrimMessages:  true,
		NoColors:      true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	dynHosts, updateHost := GetDynHosts(os.Args[1])

	var (
		ip    string
		nochg uint
	)

	ip = GetIP()

	for i, host := range dynHosts {
		if updateHost[i] && UpdateDynHost(host, ip) {
			nochg += 1
		}
	}
	if nochg == uint(len(dynHosts)) {
		log.Info("All IPs were up to date")
	}
}

func GetIP() string {

	ipGivers := [...]string{
		"https://ipinfo.io/ip",
		"http://checkip.dyndns.com/",
		"https://api.ipify.org/?format=text",
		"https://myexternalip.com/raw",
		"http://whatismyip.akamai.com/"}

	regex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

	failedCount := 0
	for _, ipGiver := range ipGivers {
		resp, err := http.Get(ipGiver)
		if err != nil {
			log.Debug("Could not get IP from " + ipGiver)
			failedCount++
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		return regex.FindString(string(body))
	}

	if failedCount == len(ipGivers) {
		log.Fatal("Could not get public IP")
	}

	return ""
}

func GetUrl(urlTemplate, hostname, ip string) (url string) {
	url = urlTemplate
	url = strings.ReplaceAll(url, TEMPLATE_HOST, hostname)
	url = strings.ReplaceAll(url, TEMPLATE_IP, ip)
	return
}

func SendRequest(url, login, passwd string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
	}
	req.SetBasicAuth(login, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s
}

func UpdateDynHost(host DynamicHost, ip string) (nochg bool) {
	url := GetUrl(host.UrlTemplate, host.Hostname, ip)
	resp := SendRequest(url, host.Login, host.Password)

	nochg = false
	if strings.Contains(resp, "good") {
		log.Info("Update successful for " + host.Hostname + " : " + resp)
	} else if strings.Contains(resp, "nochg") {
		log.Debug("Update didn't do anything for " + host.Hostname + " : " + resp)
		nochg = true
	} else {
		log.Error("Update failed for " + host.Hostname + " : " + resp)
	}
	return
}
