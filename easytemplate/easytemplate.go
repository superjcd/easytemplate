package easytemplate

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

func fetchTemplates(templateSearch string) map[string]string {
	allTemplates := make(map[string]string)
	// add template: repository pair here
	allTemplates["vue3-vite-ts"] = "https://github.com/superjcd/vue3-vite-ts-template"
	allTemplates["fastapi-peewee-jwt-scheduler"] = "https://github.com/superjcd/fastapi-peewee-jwt-scheduler"

	if templateSearch != "" {
		filteredTemplates := make(map[string]string)

		for t, r := range allTemplates {
			if strings.Contains(t, templateSearch) {
				filteredTemplates[t] = r
			}
		}
		return filteredTemplates

	}

	return allTemplates
}

func fetchTemplateOptions(templateName, templateRepo string) ([]string, error) {
	var err error
	cmd := exec.Command("git", "ls-remote", "--heads", templateRepo)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	optionsRaw := string(output)
	optionsLines := strings.Split(optionsRaw, "\n")

	templateOptions := make([]string, 0)

	for _, line := range optionsLines {
		if line != "" {
			re := regexp.MustCompile(`refs/heads/(\w+)`)
			matches := re.FindStringSubmatch(line)
			if matches != nil {
				repoName := matches[1]
				if repoName != "main" && repoName != "master" {
					templateOptions = append(templateOptions, repoName)
				}
			}
		}
	}

	if len(templateOptions) == 0 {
		err = fmt.Errorf("oops, no options for template %s", templateName)
	}

	return templateOptions, err
}

func downloadTempalate(repo, option, dirName string) {
	var err error
	cmd := exec.Command("git", "clone", "--branch", option, repo, dirName)
	_, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("download template failed: %v\n", err)
		return
	}
}

func Run() {
	var templateSearch string
	var err error

	inputArgs := os.Args
	if len(inputArgs) == 2 {
		templateSearch = inputArgs[1]
	} else if len(inputArgs) > 2 {
		fmt.Println("wrong input arguments, you can provide only one template name")
		return
	}

	templates := fetchTemplates(templateSearch)
	if len(templates) == 0 {
		fmt.Printf("oops, no templates found")
		return
	}
	possibleTemplates := make([]string, 0)

	for t, _ := range templates {
		possibleTemplates = append(possibleTemplates, t)
	}

	promptTemplate := promptui.Select{
		Label: "select a template",
		Items: possibleTemplates,
	}

	_, templateName, err := promptTemplate.Run()

	if err != nil {
		fmt.Printf("prompt failed %v\n", err)
		return
	}

	templateRepo := templates[templateName]

	options, err := fetchTemplateOptions(templateName, templateRepo)

	if err != nil {
		fmt.Printf("Get template options failed %v\n", err)
		return
	}

	promptTemplateOptions := promptui.Select{
		Label: "select a version",
		Items: options,
	}

	_, templateOption, err := promptTemplateOptions.Run()

	if err != nil {
		fmt.Printf("prompt failed %v\n", err)
		return
	}

	// fetch the template
	var direName string
	promptName := promptui.Prompt{
		Label: "your project name?[myproject]",
	}

	direName, err = promptName.Run()

	if err != nil {
		fmt.Printf("prompt failed %v\n", err)
		return
	}

	if direName == "" {
		direName = "myproject"
	}
	downloadTempalate(templateRepo, templateOption, direName)

}
