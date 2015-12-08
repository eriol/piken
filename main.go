package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"

	"eriol.xyz/piken/format"
	"eriol.xyz/piken/sql"
)

const (
	unicodeDataUrl      = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"
	pikenHome           = ".piken"
	defaultDatabaseFile = "piken.sqlite3"
	defaultDataFile     = "UnicodeData.txt"
	version             = "0.1a"
)

var (
	baseDir      = path.Join(getHome(), pikenHome)
	databaseFile = path.Join(baseDir, defaultDatabaseFile)
	dataFile     = path.Join(baseDir, defaultDataFile)
	store        sql.Store
)

func main() {

	app := cli.NewApp()
	app.Name = "piken"
	app.Version = version
	app.Author = "Daniele Tricoli"
	app.Email = "eriol@mornie.org"
	app.Usage = "unicode search tool"

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		os.Mkdir(baseDir, 0755)
	}

	if err := store.Open(databaseFile); err != nil {
		logrus.Fatal(err)
	}
	defer store.Close()

	app.Commands = []cli.Command{
		{
			Name:  "update",
			Usage: "Update unicode data",
			Action: func(c *cli.Context) {
				modifiedTime, err := checkLastModified(unicodeDataUrl)
				if err != nil {
					logrus.Fatal(err)
				}

				lastUpdate, err := store.GetLastUpdate(defaultDataFile)
				if err != nil {
					logrus.Fatal(err)
				}

				if lastUpdate.Before(modifiedTime) {
					download(unicodeDataUrl, dataFile)

					records, err := readCsvFile(dataFile)
					if err != nil {
						logrus.Fatal(err)
					}

					if err := store.LoadFromRecords(records); err != nil {
						logrus.Fatal(err)
					}

					if err := store.CreateLastUpdate(defaultDataFile,
						modifiedTime); err != nil {
						logrus.Fatal(err)
					}
				} else {
					logrus.Info("Already up to date.")
				}

			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search for unicode",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "copy, c",
					Usage: "copy glyph to clipboard",
				},
				cli.BoolTFlag{
					Name:  "show-glyph",
					Usage: "show glyph (defaults to true, use --show-glyph=false to disable)",
				},
				cli.StringFlag{
					Name:  "separator",
					Value: " -- ",
					Usage: "separator for unicode fields",
				},
			},
			Action: func(c *cli.Context) {
				args := strings.Join(c.Args(), " ")
				rows, err := store.SearchUnicode(args)
				if err != nil {
					logrus.Fatal(err)
				}

				formatter := format.NewTextFormatter(
					[]string{"CodePoint", "Name"},
					c.String("separator"),
					c.Bool("show-glyph"))

				if c.Bool("copy") && len(rows) > 1 {
					logrus.Warn("Copy to clipboard not allowed for multiple rows.")
				}

				for _, row := range rows {

					b, err := formatter.Format(&row)
					if err != nil {
						logrus.Fatal(err)
					}
					fmt.Println(b)

					// Copy to clipboard only when one row is returned by search.
					if c.Bool("copy") && len(rows) == 1 {
						glyph, err := format.CodePointToGlyph(row.CodePoint)
						if err != nil {
							logrus.Fatalf("Impossible to convert %s to glyph.",
								row.CodePoint)
						}
						if err := clipboard.WriteAll(glyph); err != nil {
							logrus.Fatalf("Copy to clipboard failed: %v", err)
						}
					}
				}
			},
		},
	}

	app.Run(os.Args)

}
