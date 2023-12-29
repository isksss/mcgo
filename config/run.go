package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// サーバーを起動
func (c Config) RunServer() error {
	serverName := fmt.Sprintf("%s-%s.jar", c.Server.Project, c.Server.Version)
	cmdPath, err := GetCmdPath("java")
	if err != nil {
		return err
	}

	// サーバーを起動
	cmd := exec.Command(cmdPath, "-Xmx"+c.Server.Memory, "-Xms"+c.Server.Memory, "-jar", serverName, "nogui")

	// 標準入出力を取得
	stdin, stdout, stderr, err := getStd(cmd)
	if err != nil {
		return err
	}

	// サーバーを起動
	if err := cmd.Start(); err != nil {
		return err
	}

	// 標準入力を受け付ける
	handleInput(stdin)

	// 標準出力を表示
	printOutput(stdout)

	// 標準エラー出力を表示
	printOutput(stderr)

	// サーバーが終了したら終了
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// 標準入出力を取得
func getStd(c *exec.Cmd) (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	stdin, err := c.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	return stdin, stdout, stderr, nil
}

// 出力を表示
func printOutput(reader io.Reader) {
	go func() {
		for {
			var b = make([]byte, 1024)
			n, err := reader.Read(b)
			if err != nil {
				break
			}
			fmt.Print(string(b[:n]))
		}
	}()
}

// 標準入力を受け付ける
func handleInput(stdin io.WriteCloser) {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			io.WriteString(stdin, scanner.Text()+"\n")
		}
	}()
}
