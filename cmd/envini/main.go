package main

import (
	"fmt"
	"os"
	"path"

	"github.com/runeimp/envini"
)

const (
	appName    = "EnvINI"
	appLabel   = "EnvINI v1.0.0"
	appVersion = "1.0.0"
	usage      = `%s

USAGE: %s [OPTIONS] INIFILE

OPTIONS:
  -h, -help, --help    Display this help info
  -v, -ver, --version  Display app version info

`
)

var configPath = "./env.ini"

var config struct {
	ProjectName string  `env:"PROJECT_NAME" ini:"Project Name"`
	LuckyAgent  float64 `ini:"lucky_agent" default:"12"`
	SecondBool  bool    `ini:"second_bool" default:"false"`
	TrueBool    bool    `ini:"true_bool" default:"true"`
	Context     struct {
		SectionText string `ini:"section_text" default:"Quoth the Raven “Nevermore.”"`
	} `ini:"Context"`
	BookOfNumbers struct {
		FloatTest float32 `ini:"float_test"`
		TheAnswer uint8   `ini:"the_answer" env:"THE_ANSWER"`
	} `ini:"Book of Numbers"`
}

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "-help", "--help":
			fmt.Printf(usage, appLabel, path.Base(os.Args[0]))
			os.Exit(0)
		case "-v", "-ver", "-version", "--version":
			fmt.Println(appLabel)
			os.Exit(0)
		default:
			configPath = os.Args[1]
		}
	} else {
		fmt.Printf(usage, appLabel, path.Base(os.Args[0]))
		os.Exit(0)
	}

	err := envini.GetConfig(configPath, &config)
	if err != nil {
		panic(err)
	}
	jsonStr, err := envini.GetConfigJSON(configPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s as JSON:\n%s\n", configPath, jsonStr)
	// fmt.Printf("config.LuckyAgent: %.3f\n", config.LuckyAgent)
	// fmt.Printf("config.ProjectName: %q\n", config.ProjectName)
	// fmt.Printf("config.SecondBool: %t\n", config.SecondBool)
	// fmt.Printf("config.TrueBool: %t\n", config.TrueBool)
	// fmt.Printf("config.Context.SectionText: %q\n", config.Context.SectionText)
	// fmt.Printf("config.BookOfNumbers.FloatTest: %.3f\n", config.BookOfNumbers.FloatTest)
	// fmt.Printf("config.BookOfNumbers.TheAnswer: %d\n", config.BookOfNumbers.TheAnswer)
}
