package main

import (
	_ "io/ioutil"
	"log"
	"os"
	"os/exec"
	_ "path/filepath"
	"regexp"
	"runtime"
	"strings"
)
import "fmt"



func main() {

	re := regexp.MustCompile("^\\-\\-([^-]+)=([^-]+)$")
	config := readConfig()

	fmt.Printf("Config is %s\n", config)

	parameters := []ParameterValue{}
	args := []string{}
	for _, arg := range os.Args[1:] {
		match := re.FindStringSubmatch(arg)
		if match != nil {
			parameters = append(parameters, ParameterValue{
				Name: match[1],
				Value: match[2]})
		} else {
			args = append(args, arg)
		}
	}

	fmt.Printf("Parameters are %s\n", parameters)
	fmt.Printf("Arguments are %s\n", args)

	var link = config.Links[args[0]]
	if isEmpty(link) {
		fmt.Printf("shortcut %s not found\n", os.Args[1])
	}
	for _, arg := range args[1:] {
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

		fmt.Printf("open %s\n", self)
		openbrowser(self)
	}


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

func isEmpty(link Link) bool {
	return link.Self == "" && link.Links == nil
}



func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
