package main

import (
	"log"

	"github.com/runeimp/envini"
)

const configPath = "./env.ini"

// type Config struct {
// 	BoolTest    bool   `ini:"bool_test"`
// 	IntTest     int    `ini:"int_test"`
// 	ProjectName string `ini:"project_name"`
// 	// Section     struct {
// 	// 	SectionText string `ini:"section_text"`
// 	// } `ini:"section"`
// }

var config struct {
	BoolTest    bool   `ini:"bool_test"`
	IntTest     int    `ini:"int_test"`
	ProjectName string `ini:"project_name"`
	Context     struct {
		SectionText string `ini:"section_text"`
	} `ini:"section"`
}

func main() {
	// config := Config{ProjectName: "Pheonix"}
	envini.GetConfig(configPath, &config)
	log.Printf("CmdEnvINI.main() | config.BoolTest: %t\n", config.BoolTest)
	log.Printf("CmdEnvINI.main() | config.IntTest: %d\n", config.IntTest)
	log.Printf("CmdEnvINI.main() | config.ProjectName: %q\n", config.ProjectName)
}
