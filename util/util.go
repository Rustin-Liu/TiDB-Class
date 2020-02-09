package util

import (
	"log"
	"os"
)

type FileType string

const (
	Cpu FileType = "Cpu"
	Mem FileType = "Mem"
)

func CreateProfile(benchType string, fileType FileType, version string) (os.File, error) {
	f, err := os.Create("prof/" + version + "/" + benchType + string(fileType) + ".prof")
	if err != nil {
		log.Fatalf("could not create %s sort %s profile: %s", benchType, fileType, err.Error())
		return *f, err
	}
	return *f, nil
}
