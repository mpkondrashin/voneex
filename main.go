package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/mpkondrashin/vone"
)

var (
	csvFilename     = "iana-ipv4-special-registry-1.csv"
	specialRegistry = "https://www.iana.org/assignments/iana-ipv4-special-registry/" + csvFilename
)

func DownloadFile(url, filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println("Special registry file already exists, skipping download")
		return nil
	}
	resp, err := http.Get(url + filename)
	if err != nil {
		return fmt.Errorf("failed to download Special registry file: %w", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save Special registry content: %w", err)
	}
	log.Println("Special registry file downloaded and saved successfully")
	return nil
}

func IterateCSV(filename string, callback func(string) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the first line (header)
	_, err = reader.Read()
	if err != nil {
		return err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) > 0 {
			err := callback(record[0])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	apiKey := flag.String("apikey", "", "The Vision One API key")

	flag.Parse()

	if *apiKey == "" {
		fmt.Println("Error: --apikey is required")
		flag.Usage()
		return
	}

	//	fmt.Printf("The provided API key is: %s\n", *apiKey)

	if err := DownloadFile(specialRegistry, csvFilename); err != nil {
		log.Fatal(err)
	}

	domain, err := vone.DetectVisionOneDomain(context.TODO(), *apiKey, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Detected domain: %s", domain)
	vOne := vone.NewVOne(domain, *apiKey)
	addException := vOne.AddExceptions()

	err = IterateCSV(csvFilename, func(ip string) error {
		for _, s := range strings.Split(ip, ",") {
			// Remove suffix with square brackets and number
			s = regexp.MustCompile(`\s*\[\d+\]$`).ReplaceAllString(s, "")
			s = strings.TrimSpace(s)
			addException.AddIP(s, "IANA IPv4 Special-Purpose Address Registry")
			log.Printf("Add %s", s)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Run API call")
	response, err := addException.Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range *response {
		if r.Status/100 != 2 {
			log.Printf("%s: %s", r.Body.Error.Code, r.Body.Error.Message)
		}
	}
	log.Println("done")
}
