package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const version = "0.3.0-dev"

const baseSongsFile = "songs.txt"
const baseAssetFile = "music.asset"

const warSuffix = "__war"

const songTemplate = `
song = {
	name = "%v"
	
	chance = {
%v
	}
}
`
const modifierTemplate = `		modifier = {
			factor = %v
		}`
const modifierWarTemplate = `		modifier = {
			factor = %v
			is_at_war = %v
		}`
const assetTemplate = `
music = {
	name = "%v"
	file = "%v"
}`

// Factor modifier (default: 1)
var Factor = "1"

// War modifier (default: "")
var War = ""

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
	rootCmd.Flags().StringVarP(&Factor, "factor", "f", "1", "global 'factor' modifier to set")
	rootCmd.Flags().StringVarP(&War, "war", "w", "", "global 'is_at_war' modifier to set (yes / no)")
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
	for _, file := range entries {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".ogg") {
			modifier := ""
			songName, atWar := songNameAndIsAtWar(strings.TrimSuffix(file.Name(), ".ogg"))
			if atWar == "" {
				modifier = fmt.Sprintf(modifierTemplate, Factor)
			} else {
				modifier = fmt.Sprintf(modifierWarTemplate, Factor, atWar)
			}
			songs += fmt.Sprintf(songTemplate, songName, modifier)
			asset += fmt.Sprintf(assetTemplate, songName, file.Name())
		}
	}

	if songs == "" {
		fmt.Printf("No .ogg files found at %v\n", dir)
		return
	}

	songsPath := filepath.Join(dir, baseSongsFile)
	ioutil.WriteFile(songsPath, []byte(songs), 0775)
	assetPath := filepath.Join(dir, baseAssetFile)
	ioutil.WriteFile(assetPath, []byte(asset), 0775)

	fmt.Printf(`Following output files generated:
  %v
  %v

Add the contents to the original EU4 music/songs.txt and music/music.asset files.
`, songsPath, assetPath)
}

func songNameAndIsAtWar(fileName string) (string, string) {
	songName := fileName
	atWar := War
	if strings.HasSuffix(fileName, warSuffix) {
		songName = strings.TrimSuffix(songName, warSuffix)
		atWar = "yes"
	}
	return songName, atWar
}

// Execute runs root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
