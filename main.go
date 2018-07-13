package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"sort"
	"time"
)

type config struct {
	Zones         map[string]string `json:"zones"`
	ReferenceTime string            `json:"reference_time"`
}
type result struct {
	name string
	zone string
	time time.Time
}

func getConfig() (config, error) {
	locationList := config{}

	user, err := user.Current()
	if err != nil {
		return locationList, fmt.Errorf("error getting current user: %s", err)
	}

	locData, err := ioutil.ReadFile(user.HomeDir + "/.worldclock.json")
	if err != nil {
		return locationList, fmt.Errorf("error reading worldclock.json file: %s", err)
	}

	if err := json.Unmarshal(locData, &locationList); err != nil {
		return locationList, fmt.Errorf("error parsing location map: %s", err)
	}
	return locationList, nil
}

func main() {
	ref := flag.String("a", "", "print time using reference date instead of now")
	flag.Parse()

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	var rt time.Time
	if *ref != "" {
		timeLayout := config.ReferenceTime

		var err error
		rt, err = time.Parse(timeLayout, *ref)
		if err != nil {
			log.Fatal("error parsing ref time: ", err)
		}
	} else {
		rt = time.Now()
	}

	output := make([]result, 0)
	for name, l := range config.Zones {
		loc, err := time.LoadLocation(l)
		if err != nil {
			log.Fatal("error loading timezone: ", err)
		}

		output = append(output, result{name, l, rt.In(loc).Truncate(1 * time.Second)})
	}

	sort.Slice(output, func(i, j int) bool {
		_, oi := output[i].time.Zone()
		_, oj := output[j].time.Zone()
		return oi < oj
	})
	for _, t := range output {
		fmt.Printf("%-10s %-20s %s\n", t.name, t.zone, t.time)
	}
}
