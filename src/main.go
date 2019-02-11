package main

import (
	_ "io/ioutil"
	"os"
	_ "path/filepath"
	"regexp"
	"strings"
)
import "fmt"

func main() {

	keys := extractKeys(os.Args)
	parameters := extractParameters(os.Args)
	options := extractOptions(os.Args)
	_, isDebug := options["X"]

	config := readConfig()

	if isDebug {
		fmt.Printf("Config is %s\n", config)
		fmt.Printf("Parameters are %s\n", parameters)
		fmt.Printf("Keys are %s\n", keys)
	}

	if len(keys) == 0 {
		fmt.Println("please provide any parameter")
		return
	}

	var link = config.Links[keys[0]]
	if isEmpty(link) {
		fmt.Printf("shortcut %s not found\n", keys[0])
	}
	for _, arg := range keys[1:] {
		if val, ok := link.Links[arg]; ok {
			link = val
		} else {
			fmt.Printf("shortcut %s not found\n", arg)
			link = Link{}
			break
		}
	}

	if !isEmpty(link) {
		self := resolveParameters(link, parameters)
		if isDebug {
			fmt.Printf("open %s\n", self)
		}
		openbrowser(self)
	}

}

func extractParameters(args []string) []ParameterValue {
	re := regexp.MustCompile("^\\-\\-([^-]+)=([^-]+)$")
	parameters := []ParameterValue{}
	for _, arg := range os.Args[1:] {
		match := re.FindStringSubmatch(arg)
		if match != nil {
			parameters = append(parameters, ParameterValue{
				Name:  match[1],
				Value: match[2]})
		}
	}
	return parameters
}

func extractKeys(args []string) []string {
	keys := []string{}
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			keys = append(keys, arg)
		}
	}
	return keys
}

func extractOptions(args []string) map[string]struct{} {
	re := regexp.MustCompile("^\\-([^-]+)$")
	options := make(map[string]struct{})
	for _, arg := range os.Args[1:] {
		match := re.FindStringSubmatch(arg)
		if match != nil {
			options[match[1]] = struct{}{}
		}
	}
	return options
}

func resolveParameters(link Link, parameters []ParameterValue) string {
	self := link.Self

	for _, par := range parameters {
		fmt.Println(fmt.Sprintf("${%s}", par.Name))
		fmt.Println(self)
		self = strings.Replace(self, fmt.Sprintf("${%s}", par.Name), par.Value, -1)
		fmt.Println(self)
	}

	return self
}
