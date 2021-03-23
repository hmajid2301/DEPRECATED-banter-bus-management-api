package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type arrayFlags []string

const (
	usage = `usage: %s
Create devcontainer files

Options:
`
)

func main() {
	var extension arrayFlags
	shellPtr := flag.String("shell", "", "The shell to use within the devcontainer.")
	flag.Var(&extension, "fisher", "List of fisher extensions to install.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	fmt.Printf("devcontainer will use Shell: %s \n", *shellPtr)

	shell := getShell(*shellPtr)
	err := createTemplate("docker-compose.dev.yml", shell)
	if err != nil {
		fmt.Println("failed to create docker-compose.yml")
		panic(err)
	}

	err = createTemplate("devcontainer.json", shell)
	if err != nil {
		fmt.Println("failed to create devcontainer.json")
		panic(err)
	}

	if *shellPtr == "fish" && extension != nil {
		fmt.Printf("Adding fisher extensions to postInst %s \n", extension)
		err = createFishPostInst(extension)
		if err != nil {
			fmt.Println("failed to create postInst.fish")
			panic(err)
		}
	}
}

func getShell(shellType string) interface{} {
	type Shell struct {
		BinaryLocation string
		HistoryPath    string
		PostCreate     string
	}

	switch shellType {
	case "fish":
		return Shell{
			BinaryLocation: "/usr/bin/fish",
			HistoryPath:    "/home/vscode/.local/share/fish/fish_history",
			PostCreate:     "fish /home/vscode/postInst.fish",
		}
	default:
		return Shell{
			BinaryLocation: "/usr/bin/bash",
			HistoryPath:    "/home/vscode/.bash_history",
			PostCreate:     "bash /home/vscode/postInst.bash",
		}
	}
}

func createTemplate(fileName string, commands interface{}) error {
	filePath := fmt.Sprintf("utils/templates/%s", fileName)
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}

	outputPath := fmt.Sprintf(".devcontainer/%s", fileName)
	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	err = tpl.Execute(output, commands)
	if err != nil {
		return err
	}
	return nil
}

func createFishPostInst(extensions arrayFlags) error {
	err := createTemplate("postInst.fish", extensions)
	return err
}

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}
