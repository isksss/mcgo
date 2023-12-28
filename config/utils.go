package config

import (
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
)

const (
	apiUrl = "https://api.papermc.io/v2/projects"
)

type Project struct {
	Error         bool     `json:"error"`
	Projects      []string `json:"projects"`
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
	Version       string   `json:"version"`
	Builds        []int    `json:"builds"`
}

// コマンドpathを返す
func GetCmdPath(c string) (string, error) {

	path, err := exec.LookPath(c)
	if err != nil {
		return "", err
	}

	return path, nil
}

// httpRequestを行いjsonを返す
func GetJson(url string) (Project, error) {
	var project Project

	body, err := GetBody(url)
	if err != nil {
		return project, err
	}

	err = json.Unmarshal(body, &project)
	if err != nil {
		return project, err
	}

	return project, nil
}

// httpRequestを行いbodyを返す
func GetBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return body, nil
}
