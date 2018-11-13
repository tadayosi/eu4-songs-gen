package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const version = "0.1"
const songsFile = "songs.txt"
const assetFile = "music.asset"

var rootCmd = &cobra.Command{
	Use:   "eu4-songs-gen [path to .ogg files]",
	Short: "Generate EU4 songs.txt from .ogg files",
	Long: `Europa Universalis IV - Song List Generator
This command enables you to add your own song list to EU4
by generating songs.txt from local .ogg files.`,
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	Run:     run,
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("EU4 Song List Generator %s\n", rootCmd.Version))
}

func run(cmd *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	songs := ""
	asset := ""
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".ogg") {
			songs += fmt.Sprintf(`
song = {
	name = "%v"
	
	chance = {
		modifier = {
			factor = 1
		}
	}
}
`, strings.TrimSuffix(entry.Name(), ".ogg"))
			asset += fmt.Sprintf(`
music = {
	name = "%v"
	file = "%v"
}`, strings.TrimSuffix(entry.Name(), ".ogg"), entry.Name())
		}
	}

	if songs == "" {
		fmt.Printf("No .ogg files found at %v\n", dir)
		return
	}

	songsPath := filepath.Join(dir, songsFile)
	ioutil.WriteFile(songsPath, []byte(songs), 0775)
	assetPath := filepath.Join(dir, assetFile)
	ioutil.WriteFile(assetPath, []byte(asset), 0775)

	fmt.Printf(`Following output files generated:
  %v
  %v
Add the contents to the original EU4 music/songs.txt and music/music.asset files.
`, songsPath, assetPath)
}

// Execute runs root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
