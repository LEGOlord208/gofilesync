package main

import (
	"encoding/json"
	"os"

	"github.com/legolord208/stdutil"
)

type location struct {
	Src string
	Dst string
}

var data struct {
	Schedule  int
	Locations []location
}

const dataFile = ".gofilesync"

func loadData() {
	f, err := os.Open(dataFile)
	if err != nil {
		if !os.IsNotExist(err) {
			stdutil.PrintErr("Could not load data file", err)
		}
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		stdutil.PrintErr("Could not decode data file", err)
		return
	}
}

func saveData() string {
	f, err := os.Create(dataFile)
	if err != nil {
		stdutil.PrintErr("Could not load data file", err)
		return "Failed to create file"
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(&data)
	if err != nil {
		stdutil.PrintErr("Could not encode data file", err)
		return "Failed to encode file"
	}
	return "Saved data!"
}
