package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	LETSENCRYPTDIR = "/etc/letsencrypt"
)

var (
	Version         = "development"
	CheckAllDomains = false
)

func main() {
	version := flag.Bool("version", false, "version")
	letsencryptdir := flag.String("letsencryptdir", LETSENCRYPTDIR, "Root directory for lestencrypt installation")
	domain := flag.String("domain", "", "Domain to be checked (otherwise all domains will be checked)")
	flag.Parse()

	if *version {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

	if *letsencryptdir != LETSENCRYPTDIR {
		fmt.Println("Using directory:", *letsencryptdir)
	}

	if *domain == "" {
		fmt.Println("Checking all domains.")
		CheckAllDomains = true
	}

	if CheckAllDomains {
		configFiles := getListOfConfigs(fmt.Sprintf("%s/renewal", *letsencryptdir))
		fmt.Println("configFiles:", configFiles)
	} else {
		configFile := fmt.Sprintf("%s/renewal/%s.conf", *letsencryptdir, *domain)
		certificate := getCertificatePathFromConfig(configFile)
		fmt.Println("certificate:", certificate)

	}

}

func getListOfConfigs(directory string) []string {
	files, err := os.ReadDir(directory)
	if err != nil {
		panic(fmt.Sprintf("Failed to read from directory: %v", directory))
	}

	var configFiles []string
	for _, file := range files {
		if grep(".conf", file.Name()) {
			configFiles = append(configFiles, file.Name())
		}
	}
	return configFiles
}

func grep(pattern, text string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

func getCertificatePathFromConfig(configFile string) string {
	file, err := os.Open(configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file: %v", configFile))
	}
	defer file.Close()

	var certificate = ""
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		strLine := string(line)
		if grep("^cert", strLine) {
			strLine = sed(" ", "", strLine)
			certificate = strings.Split(strLine, "=")[1]
		}
	}
	return certificate
}

func sed(before, after, text string) string {
	return strings.Replace(text, before, after, -1)
}
