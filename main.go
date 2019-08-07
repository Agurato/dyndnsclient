package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/glendc/go-external-ip"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&nested.Formatter{
		ShowFullLevel: true,
		TrimMessages:  true,
	})
	log.SetOutput(os.Stdout)

	dynHosts, updateHost := GetDynHosts(os.Args[1])

	var (
		ip string
	)

	ip, err := GetIP()
	if err != nil {
		log.Fatal(err)
	}

	for i, host := range dynHosts {
		if updateHost[i] {
			UpdateDynHost(host, ip)
		}
	}
}

func GetIP() (string, error) {
	// Create the default consensus,
	// using the default configuration and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)
	// Get your IP,
	// which is never <nil> when err is <nil>.
	ip, err := consensus.ExternalIP()
	if err == nil {
		return ip.String(), nil
	}
	return "", err
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

func UpdateDynHost(host DynamicHost, ip string) {
	url := GetUrl(host.UrlTemplate, host.Hostname, ip)
	resp := SendRequest(url, host.Login, host.Password)

	if strings.Contains(resp, "good") {
		log.Info("Update successful for " + host.Hostname + " : " + resp)
	} else if strings.Contains(resp, "nochg") {
		log.Warn("Update didn't do anything for " + host.Hostname + " : " + resp)
	} else {
		log.Error("Update failed for " + host.Hostname + " : " + resp)
	}
}
