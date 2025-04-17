package main

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
	"path/filepath"
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
		// fmt.Println("Checking all domains.")
		CheckAllDomains = true
	}

	if CheckAllDomains {
		configFiles := getListOfConfigs(fmt.Sprintf("%s/renewal", *letsencryptdir))
		// fmt.Println("configFiles:", configFiles)
		for _, config := range configFiles {
			config = fmt.Sprintf("%s/renewal/%s", *letsencryptdir, config)
			cert := getCertificatePathFromConfig(config)
			certX509 := getCertificateX509(cert)
			days := getRemainingDays(certX509)
			domain := getDomainNameFromConfigFile(config)
			printRemainingDays(domain, days)
		}
	} else {
		configFile := fmt.Sprintf("%s/renewal/%s.conf", *letsencryptdir, *domain)
		certificate := getCertificatePathFromConfig(configFile)
		// fmt.Println("certificate:", certificate)
		certX509 := getCertificateX509(certificate)
		days := getRemainingDays(certX509)
		printRemainingDays(*domain, days)

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

func getCertificateX509(certificate string) *x509.Certificate {
	file, err := os.Open(certificate)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file: %v", certificate))
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	// 8192 is a good number?
	text := make([]byte, 8192)
	counter, err := rd.Read(text)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file: %v", certificate))
	}
	if counter == 0 {
		panic(fmt.Sprintf("Failed to read file (counter is zero): %v", certificate))
	}
	certPEM, _ := pem.Decode(text)
	if certPEM == nil {
		panic(fmt.Sprintf("Failed to decode PEM from file: %v", certificate))
	}
	cert, err := x509.ParseCertificate(certPEM.Bytes)
		if certPEM == nil {
		panic(fmt.Sprintf("Failed to parse as X509 from file: %v", certificate))
	}

	return cert
}

func getRemainingDays(cert *x509.Certificate) string {
	expiration := cert.NotAfter
	remaining := expiration.Sub(time.Now())
	return fmt.Sprintf("%d", int(remaining.Hours()/24))
}

func printRemainingDays(domain , days string) {
	fmt.Printf("%s=%s\n", domain, days)
}

func getDomainNameFromConfigFile(config string) string {
	fileName := filepath.Base(config)
	return sed(".conf", "", fileName)
}
