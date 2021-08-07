package main

import (
	"log"

	"github.com/runeimp/envini"
)

const configPath = "./env.ini"

var config struct {
	BoolTest    bool    `ini:"bool_test"`
	LuckyAgent  float64 `ini:"lucky_agent"`
	TheAnswer   int     `ini:"the_answer"`
	ProjectName string  `ini:"project_name"`
	Context     struct {
		SectionText string `ini:"section_text"`
	} `ini:"Context"`
}

func main() {
	envini.GetConfig(configPath, &config)
	log.Printf("CmdEnvINI.main() | config.BoolTest: %t\n", config.BoolTest)
	log.Printf("CmdEnvINI.main() | config.LuckyAgent: %.3f\n", config.LuckyAgent)
	log.Printf("CmdEnvINI.main() | config.TheAnswer: %d\n", config.TheAnswer)
	log.Printf("CmdEnvINI.main() | config.ProjectName: %q\n", config.ProjectName)
	log.Printf("CmdEnvINI.main() | config.Context.SectionText: %q\n", config.Context.SectionText)
}
