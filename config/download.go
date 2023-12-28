package config

import (
	"fmt"
	"os"
	"path"
	"sync"
)

// Download
func (c Config) DownloadServer() error {
	// check project
	if err := checkProject(c.Server.Project); err != nil {
		return err
	}

	// check version
	if err := checkVersion(c.Server.Project, c.Server.Version); err != nil {
		return err
	}

	// get build bumber
	j, err := GetJson(apiUrl + "/" + c.Server.Project + "/versions/" + c.Server.Version)
	if err != nil {
		return err
	}
	build := j.Builds[len(j.Builds)-1]

	// download server
	jarName := fmt.Sprintf("%s-%s-%d.jar", c.Server.Project, c.Server.Version, build)
	url := fmt.Sprintf("%s/%s/versions/%s/builds/%d/downloads/%s", apiUrl, c.Server.Project, c.Server.Version, build, jarName)

	fmt.Println("URL: " + url + "")
	jar, err := GetBody(url)
	if err != nil {
		return err
	}

	// Write jar
	f, err := os.Create(jarName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(jar)
	if err != nil {
		return err
	}
	return nil
}

func (c Config) DownloadPlugins() error {
	// create plugins dir
	if err := os.Mkdir("plugins", 0777); err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, plugin := range c.Plugins {
		wg.Add(1)
		path := path.Join("plugins", plugin.Name)
		url := plugin.Url
		go func() {
			defer wg.Done()
			fmt.Printf("Download plugin: %s\n", plugin.Name)
			fmt.Println("URL: " + url + "")
			body, _ := GetBody(url)
			// Write jar
			f, _ := os.Create(path)
			defer f.Close()

			_, _ = f.Write(body)
		}()

	}

	wg.Wait()
	return nil
}
