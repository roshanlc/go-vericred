package main

import (
	"os/exec"
	"strings"
)

func searchUniqueZips(fileName string, data chan string) {
	cmd := "rg -o '\"zip\": \"\\d{5}\"' " + fileName + " | parsort | uniq -c | tee unique_zips.txt | wc -l"
	zips := exec.Command("bash", "-c", cmd)

	output, err := zips.Output()
	if err != nil {
		data <- err.Error()
		return
	}
	msg := strings.Repeat("-", 50)

	msg += "\nTotal unique zips found: " + string(output) + "The unique zips are saved in unique_zips.txt.\n"

	msg += strings.Repeat("-", 50)

	data <- msg
}

func searchUniqueNetworks(fileName string, data chan string) {
	// cmd := "rg -o '\"networks\": \\[\\{\"name\": \".+\", \"tier\": \".*\"\\}\\],' " + fileName + " | parsort -u"
	cmd := "jq -r '.networks[].name' " + fileName + " | parsort | uniq -c | tee unique_networks.txt | wc -l"
	zips := exec.Command("bash", "-c", cmd)

	output, err := zips.Output()
	if err != nil {
		data <- err.Error()
		return
	}
	msg := strings.Repeat("-", 50)

	msg += "\nThe total unique networks found : " + string(output) + "The unique networks are saved in unique_networks.txt\n"

	msg += "\n" + strings.Repeat("-", 50)

	data <- msg
}

func searchUniqueSpecialties(fileName string, data chan string) {
	cmd := "jq -r '.specialties[].name' " + fileName + " | parsort | uniq -c | tee unique_specs.txt | wc -l"
	zips := exec.Command("bash", "-c", cmd)

	output, err := zips.Output()
	if err != nil {
		data <- err.Error()
		return
	}
	msg := strings.Repeat("-", 50)

	msg += "\nThe total unique specialties found : " + string(output) + "The unique specialties are saved in unique_specs.txt\n"

	msg += "\n" + strings.Repeat("-", 50)

	data <- msg
}
