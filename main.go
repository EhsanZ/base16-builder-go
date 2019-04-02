package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Sources struct {
	Schemes   string
	Templates string
}

// TODO: Pull repos that exist instead of cloning them
func main() {
	sources, err := readSources()
	if err == nil {
		cloneSourceRepos(sources)
		cloneAllSchemesAndTemplates()
	}
}

// TODO: seperate to a new package for reading sources
func readSources() (Sources, error) {
	data, err := ioutil.ReadFile("./sources.yaml")

	if err != nil {
		return Sources{}, err
	}

	var sources Sources
	yaml.Unmarshal([]byte(data), &sources)
	return sources, nil
}

// TODO: seperate to a new package for cloning repos
func cloneAllSchemesAndTemplates() {
	schemes, _ := readYamlFileIntoMap("./sources/schemes/list.yaml")

	var wg sync.WaitGroup
	for slug, url := range schemes {
		wg.Add(1)
		go clone(url, fmt.Sprintf("schemes/%s", slug), &wg)
	}

	templates, _ := readYamlFileIntoMap("./sources/templates/list.yaml")

	for slug, url := range templates {
		wg.Add(1)
		go clone(url, fmt.Sprintf("templates/%s", slug), &wg)
	}

	wg.Wait()
}

func readYamlFileIntoMap(path string) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var parsedMap map[string]string
	yaml.Unmarshal([]byte(data), &parsedMap)
	return parsedMap, nil
}

func cloneSourceRepos(sources Sources) {
	fmt.Println("cloneSourceRepos")
	var wg sync.WaitGroup

	wg.Add(2)
	go clone(sources.Schemes, "sources/schemes", &wg)
	go clone(sources.Templates, "sources/templates", &wg)
	wg.Wait()
}

func clone(url string, directory string, wg *sync.WaitGroup) {
	git.PlainClone(directory, false, &git.CloneOptions{
		URL:   url,
		Depth: 1,
	})

	wg.Done()
}
