package config

import (
	"fmt"
	"os"
	"path"
	"sync"
)

// GetLatestBuild は指定されたプロジェクトとバージョンの最新のビルド情報を取得します。
// もしプロジェクトやバージョンに関するエラーがあれば、エラーを返します。
func GetLatestBuild(project, version string) (int, error) {
	// プロジェクトの確認
	if err := checkProject(project); err != nil {
		return 0, err
	}

	// バージョンの確認
	if err := checkVersion(project, version); err != nil {
		return 0, err
	}

	// ビルド番号の取得
	j, err := GetJson(apiUrl + "/" + project + "/versions/" + version)
	if err != nil {
		return 0, err
	}

	return j.Builds[len(j.Builds)-1], nil
}

// DownloadServer は構成に基づいてサーバーのJARファイルをダウンロードします。
// もし最新のビルド情報の取得やJARファイルのダウンロードにエラーがあれば、エラーを返します。
func (c Config) DownloadServer() error {
	// 最新のビルドを取得
	build, err := GetLatestBuild(c.Server.Project, c.Server.Version)
	if err != nil {
		return err
	}

	// サーバーのダウンロード
	jarName := fmt.Sprintf("%s-%s-%d.jar", c.Server.Project, c.Server.Version, build)
	url := fmt.Sprintf("%s/%s/versions/%s/builds/%d/downloads/%s", apiUrl, c.Server.Project, c.Server.Version, build, jarName)

	fmt.Println("URL: " + url + "")
	jar, err := GetBody(url)
	if err != nil {
		return err
	}

	// JARファイルの書き込み
	fileName := fmt.Sprintf("%s-%s.jar", c.Server.Project, c.Server.Version)
	if err := writeJar(fileName, jar); err != nil {
		return err
	}

	return nil
}

// DownloadPlugins は構成に基づいてプラグインをダウンロードします。
// もしプラグインディレクトリの作成やプラグインのダウンロードにエラーがあれば、エラーを返します。
func (c Config) DownloadPlugins() error {
	// プラグインディレクトリの作成
	if _, err := os.Stat("plugins"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("plugins", 0777); err != nil {
				return err
			}
		}
	}

	var wg sync.WaitGroup
	for _, plugin := range c.Plugins {
		wg.Add(1)
		go downloadPlugin(plugin, &wg)
	}

	wg.Wait()
	return nil
}

// downloadPlugin は構成に基づいてプラグインをダウンロードします。
// もしプラグインのダウンロードやJARファイルの書き込みにエラーがあれば、エラーを出力します。
func downloadPlugin(p Plugin, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("プラグインのダウンロード: %s\n", p.Name)
	fmt.Println("URL: " + p.Url + "")
	body, _ := GetBody(p.Url)
	// JARファイルの書き込み
	if err := writeJar(path.Join("plugins", p.Name), body); err != nil {
		fmt.Printf("プラグイン %s のダウンロードエラー: %s\n", p.Name, err)
		return
	}
	fmt.Printf("プラグインのダウンロード完了: %s\n", p.Name)
}

// writeJar は指定されたパスにJARファイルを書き込みます。
// もしファイルの作成や書き込みにエラーがあれば、エラーを返します。
func writeJar(filePath string, body []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}
	return nil
}
