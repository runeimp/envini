package main

import (
	"log"

	"github.com/runeimp/envini"
)

const configPath = "./env.ini"

var config struct {
	BoolTest    bool   `ini:"bool_test"`
	TrueBool    bool   `ini:"true_bool"`
	ProjectName string `ini:"Project Name"`
	Context     struct {
		SectionText string `ini:"section_text"`
	} `ini:"Context"`
	BookOfNumbers struct {
		FloatTest  float32 `ini:"float_test"`
		LuckyAgent float64 `ini:"lucky_agent"`
		TheAnswer  uint8   `ini:"the_answer"`
	} `ini:"Book of Numbers"`
}

func main() {
	envini.GetConfig(configPath, &config)
	log.Printf("CmdEnvINI.main() | config.BoolTest: %t\n", config.BoolTest)
	log.Printf("CmdEnvINI.main() | config.TrueBool: %t\n", config.TrueBool)
	log.Printf("CmdEnvINI.main() | config.ProjectName: %q\n", config.ProjectName)
	log.Printf("CmdEnvINI.main() | config.Context.SectionText: %q\n", config.Context.SectionText)
	log.Printf("CmdEnvINI.main() | config.BookOfNumbers.FloatTest: %.3f\n", config.BookOfNumbers.FloatTest)
	log.Printf("CmdEnvINI.main() | config.BookOfNumbers.LuckyAgent: %.3f\n", config.BookOfNumbers.LuckyAgent)
	log.Printf("CmdEnvINI.main() | config.BookOfNumbers.TheAnswer: %d\n", config.BookOfNumbers.TheAnswer)
}
