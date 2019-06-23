package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/logrusorgru/aurora"
)

func main() {
	log.SetFlags(0)
	_, envPresent := os.LookupEnv("FF_NOCOLOR")
	au := aurora.NewAurora(!envPresent)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get current directory: %v", err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Can't get home directory: %v", err)
	}

	// you'd be surprised how easily this happens...
	if usr.HomeDir == currentDir {
		log.Fatal(au.BrightRed("Yeah, you're in your home directory, genius.  Might not want to fuck around here.").Bold())
	}

	ws := regexp.MustCompile("[[:space:]]+")
	re := regexp.MustCompile("[^a-zA-Z0-9.-_]+")

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newName := strings.ToLower(ws.ReplaceAllString(info.Name(), "_"))
		newName = re.ReplaceAllLiteralString(newName, "x")
		if newName == info.Name() {
			return nil
		}

		if _, err := os.Stat(newName); err == nil {
			log.Printf("Not renaming: %s (would collide)", au.Yellow(info.Name()))
			return nil
		}

		log.Printf("Renaming '%s' to '%s'", au.Yellow(info.Name()), au.Green(newName))
		return os.Rename(info.Name(), newName)
	})
	if err != nil {
		log.Printf("Error encountered: %v", au.BrightRed(err))
	}
}
