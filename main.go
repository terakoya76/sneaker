package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/terakoya76/sneaker/parser"
)

func main() {
	cobra.OnInitialize()
	rootCmd.DisableSuggestions = false

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "sneaker",
	Short: "sneaker parse crontab output and visualize it to find execution intervals",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := getInput()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		exps := parser.ParseCrontab(input)
		schedule := parser.InitSchedule()

		for _, exp := range exps {
			schedule, err = exp.Evaluate(schedule)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			}
		}

		fmt.Println(schedule)
	},
}

func getInput() (string, error) {
	var filename string
	flag.Parse()
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	var r io.Reader
	switch filename {
	case "", "-":
		r = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		r = f
	}

	t, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(t), nil
}
