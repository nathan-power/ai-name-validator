package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Entry struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Format      string `json:"format"`
	Stream      bool   `json:"stream"`
	Temperature int    `json:"temperature"`
}

type Response struct {
	Response string `json:"response"`
}

func main() {
	processRecords("users.csv")
}

func processRecords(csvFileName string) {
	file, err := os.Open(csvFileName)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %s", err)
	}
	defer file.Close()

	totalRecords, err := lineCounter(csvFileName)
	if err != nil {
		log.Fatalf("Failed to count lines in CSV file: %s", err)
	}

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read headers from CSV file: %s", err)
	}

	firstNameIndex, lastNameIndex := findIndex(headers, "firstName"), findIndex(headers, "lastName")
	if firstNameIndex == -1 || lastNameIndex == -1 {
		log.Fatalf("CSV does not contain required 'firstName' or 'lastName' columns")
	}

	processEachRecord(reader, firstNameIndex, lastNameIndex, totalRecords)
}

func processEachRecord(reader *csv.Reader, firstNameIndex, lastNameIndex, totalRecords int) {
	spinner := []string{"|", "/", "-", "\\"}
	var processedRecords int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read a record from CSV file: %s", err)
		}

		processedRecords++
		displayProgress(processedRecords, totalRecords, spinner)

		fullName := record[firstNameIndex] + " " + record[lastNameIndex]
		responseData := queryModel(fullName)

		if strings.Contains(responseData, "no") {
			fmt.Printf("\r\033[K%s is not likely a person's name.\n", fullName)
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("\r\033[KProcessing complete.\n")
}

func displayProgress(current, total int, spinner []string) {
	spinnerIndex := current % len(spinner)
	fmt.Printf("\r\033[KProcessing... %d%% complete %s ", (current*100)/total, spinner[spinnerIndex])
}

func queryModel(fullName string) string {
	entry := Entry{
		Model:       "phi3",
		Prompt:      fmt.Sprintf("Using your knowledge-base of possible person names or what person names simply look like, determine if '%s' could likely be a person's name. Answer only with a single 'yes' or 'no' in lowercase. Keep in mind they may or may not have hyphens. Anything that might be a person's name answer 'yes', but if it wouldn't likely be a person's a name respond with 'no'. If your response contains anything other than just a single 'yes', 'no', or 'maybe' you'll break the script you're outputting to, so never reply with anything but a single word.", fullName),
		Temperature: 0,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("Failed to encode entry to JSON: %s", err)
	}

	return postData(jsonData)
}

func lineCounter(fileName string) (int, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	return bytes.Count(data, []byte("\n")), nil
}

func findIndex(headers []string, column string) int {
	for i, header := range headers {
		if header == column {
			return i
		}
	}
	return -1
}

func postData(data []byte) string {
	url := "http://localhost:11434/api/generate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %s", err)
	}
	defer resp.Body.Close()

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Fatalf("Failed to decode response: %s", err)
	}

	return res.Response
}
