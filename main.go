package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	deviceID, err := getDeviceID()
	if err != nil {
		log.Fatal("[Couldn't get deiceID]", err.Error())
	}
	propID, err := getTappingPropID(deviceID)
	if err != nil {
		log.Fatal("[Couldn't get propID]", err.Error())
	}
	err = enableTapping(deviceID, propID)
	if err != nil {
		log.Fatal("[Couldn't enable tapping", err.Error())
	}
	fmt.Println("Succeed to enable tapping!")
}

func enableTapping(deviceID int, propID int) error {
	cmd := exec.Command("xinput", "set-prop", strconv.Itoa(deviceID), strconv.Itoa(propID), "1")
	fmt.Printf("try: %s %s\n", cmd.Path, strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// libinput Tapping Enabled (285):	1
func getTappingPropID(deviceID int) (int, error) {
	cmd := exec.Command("xinput", "list-props", strconv.Itoa(deviceID))
	bytes, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(bytes), "\n")
	matchedLine, err := find(lines, func(s string) bool {
		return strings.Contains(s, "Tapping Enabled (")
	})
	if err != nil {
		return 0, err
	}
	// e.g. afterTokens := "285):	0"
	afterTokens := strings.Split(matchedLine, "Tapping Enabled (")[1]
	propsIDAsString := strings.Split(afterTokens, ")")[0]
	return strconv.Atoi(propsIDAsString)
}

func find(strArray []string, predicate func(string) bool) (string, error) {
	for _, s := range strArray {
		if predicate(s) {
			return s, nil
		}
	}
	return "", errors.New("Not found")
}

func getDeviceID() (int, error) {
	cmd := exec.Command("xinput")
	bytes, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("[cmd: %v] %s", cmd.Args, err.Error())
	}
	var touchpadLine string
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Touchpad") {
			touchpadLine = line
			break
		}
	}
	if touchpadLine == "" {
		return 0, errors.New("touchpadLine: was not found")
	}

	columns := strings.Split(touchpadLine, "\t")
	for _, tokens := range columns {
		if strings.HasPrefix(tokens, "id=") {
			start := len("id=")
			deviceID, err := strconv.Atoi(tokens[start:])
			if err != nil {
				return 0, err
			}
			return deviceID, nil
		}
	}
	return 0, errors.New("Unexpected Exception")
}
