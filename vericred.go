package main

import (
	"os/exec"
	"strings"
)

func searchUniqueZips(fileName string, data chan string) {
	cmd := "rg -o '\"zip\": \"\\d{5}\"' " + fileName + " | parsort -u | wc -l"
	zips := exec.Command("bash", "-c", cmd)

	output, err := zips.Output()
	if err != nil {
		data <- err.Error()
		return
	}
	msg := strings.Repeat("-", 50)

	msg += "\nTotal unique zips found: " + string(output) + "The unique zips are saved in unique_zips.txt (To be implemented).\n"

	msg += strings.Repeat("-", 50)

	data <- msg
}

func searchUniqueNetworks(fileName string, data chan string) {
	// cmd := "rg -o '\"networks\": \\[\\{\"name\": \".+\", \"tier\": \".*\"\\}\\],' " + fileName + " | parsort -u"
	cmd := "jq '.networks[].name' " + fileName + " | parsort -u"
	zips := exec.Command("bash", "-c", cmd)

	output, err := zips.Output()
	if err != nil {
		data <- err.Error()
		return
	}
	msg := strings.Repeat("-", 50)

	msg += "\nThe unique networks are:\n\n"
	msg += string(output)

	msg += "\n" + strings.Repeat("-", 50)

	data <- msg
}

func searchUniqueSpecialties(fileName string, data chan string) {
	cmd := "jq '.specialties[].name' " + fileName + " | parsort -u "
	zips := exec.Command("bash", "-c", cmd)

	output, err := zips.Output()
	if err != nil {
		data <- err.Error()
		return
	}
	msg := strings.Repeat("-", 50)

	msg += "\nThe unique specialties are:\n\n"
	msg += string(output)

	msg += "\n" + strings.Repeat("-", 50)

	data <- msg
}
