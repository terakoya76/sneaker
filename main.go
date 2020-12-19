package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/terakoya76/sneaker/parser"
)

var filename string
var dayFilter int

func main() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVar(&filename, "filename", "", "filename of crontab")
	rootCmd.PersistentFlags().IntVar(&dayFilter, "day", 0, "filter for displaying the specific day of execution schedule")
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
		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		input, err := getInput(filename)
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

		day, err := cmd.Flags().GetInt("day")
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		if day > 0 {
			fmt.Println(schedule[day])
		} else {
			fmt.Println(schedule.String())
		}
	},
}

func getInput(filename string) (string, error) {
	var r io.Reader
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	t, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(t), nil
}
