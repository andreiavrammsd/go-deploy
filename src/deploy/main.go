package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"os"
	"os/exec"
	"flag"
	"regexp"
	"path"
)

const configFile = "config.yml"
const remoteDestinationPattern = "(.+)@(.+):(.+)"
const hostsFile = "~/.ssh/known_hosts"

type Config struct {
	Projects map[string]Project
}

type Project struct {
	Repository string
	Branch string
	Destinations []string
}

type Destination struct {
	Full string
	Auth string
	Host string
	Path string
	Filename string
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getConfig(path string) Config {
	yamlFile, err := ioutil.ReadFile(path)
	checkError(err)

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	checkError(err)
	
	return config
}

func updateCode(codePath string, definition Project) {
	if _, err := os.Stat(codePath); err == nil {
		out, err1 := exec.Command("git", "-C", codePath, "reset", "--hard").Output()
		checkError(err1)
		fmt.Print("-- " + string(out))

		out, err2 := exec.Command("git", "-C", codePath, "pull", "--rebase").Output()
		checkError(err2)
		fmt.Print("-- " + string(out))
	} else {
		err := exec.Command("git", "clone", "-b", definition.Branch, definition.Repository, codePath).Run()
		checkError(err)
	}
}

func getDependencies() []byte {
	out, err := exec.Command("go", "get", "-d").Output()
	checkError(err)
	
	return out
}

func runTests() []byte {
	out, err := exec.Command("go", "test").Output()
	checkError(err)
	
	return out
}

func install() {
	err := exec.Command("go", "install").Run()
	checkError(err)
}

func build(binaryPath string) {
	err := exec.Command("go", "build", "-o", binaryPath).Run()
	checkError(err)
}

func getDestination(destinationString string) Destination {
	match, _ := regexp.MatchString(remoteDestinationPattern, destinationString)
	if !match {
		fmt.Printf(
			"\nInvalid destination pattern: %s\nExpected: %s\n\n",
			destinationString,
			remoteDestinationPattern,
		)
		os.Exit(0)
	}

	r, _ := regexp.Compile(remoteDestinationPattern)
	segments := r.FindStringSubmatch(destinationString)

	return Destination {
		Full: destinationString,
		Auth: segments[1],
		Host: segments[2],
		Path: path.Dir(segments[3]),
		Filename: path.Base(segments[3]),
	}
}

func sync(localBinaryPath string, destination Destination) {
	exec.Command("sh", "-c", "ssh-keyscan " + destination.Host + " >> " + hostsFile).Run()
	
	rsync := "rsync -aq --rsync-path='mkdir -p " + destination.Path + " && rsync' " + localBinaryPath + " " +
		destination.Full
	exec.Command("sh", "-c", rsync).Run()
}

func main() {
	requiredProjectPointer := flag.String("p", "", "-p project_name")
	flag.Parse()

	goPath := os.Getenv("GOPATH")
	config := getConfig(goPath + "/" + configFile)

	requiredProject := *requiredProjectPointer
	if len(requiredProject) > 1{
		_, isset := config.Projects[requiredProject]
		if !isset {
			fmt.Printf("Project %s not defined\n.", requiredProject)
			os.Exit(0)
		}
	}

	var result []byte;
	
	for project, definition := range config.Projects {
		if len(requiredProject) > 0 && project != requiredProject {
			continue
		}
		
		fmt.Printf("\nProject: %s\n", project)

		codePath := goPath + "/src/" + project

		fmt.Println("- Update code")
		updateCode(codePath, definition)

		os.Chdir(codePath)

		fmt.Println("- Dependencies")
		result = getDependencies()
		if out := string(result); len(out) > 0 {
			fmt.Println(string(result))
		}
		install()

		fmt.Println("- Run tests")
		result = runTests()
		fmt.Print("-- " + string(result))
		
		fmt.Println("- Build")
		localBinaryPath := goPath + "/bin/" + project
		build(localBinaryPath)

		fmt.Println("- Sync")
		for _, destination := range definition.Destinations {
			fmt.Printf("-- %s\n", destination)
			sync(localBinaryPath, getDestination(destination))
		}

		fmt.Println()
	}
}
