package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const (
	leader = "mcgo!"
)

var (
	restartFlag     = true
	nextRestartTime time.Time
)

// サーバーを起動
func (c Config) RunServer() error {
	for restartFlag {
		restartFlag = c.Server.Restart
		if err := c.DownloadServer(); err != nil {
			return err
		}
		if err := c.DownloadPlugins(); err != nil {
			return err
		}
		if err := run(&c); err != nil {
			return err
		}
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
			if scanner.Text() == leader+"stop" {
				io.WriteString(stdin, "stop\n")
				restartFlag = false
				fmt.Printf("再起動を停止しました\n")
				break
			}
			if scanner.Text() == leader+"restart" {
				io.WriteString(stdin, "stop\n")
				fmt.Printf("再起動します\n")
				break
			}
			if scanner.Text() == leader+"next" {
				_, err := fmt.Printf("次の再起動時刻: %s\n", nextRestartTime.Format("15:04"))
				if err != nil {
					fmt.Printf("次の再起動時刻を取得できませんでした\n")
				}
				continue
			}
			io.WriteString(stdin, scanner.Text()+"\n")
		}
	}()
}

// サーバー起動
func run(c *Config) error {
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

	// 再起動設定
	var timer *time.Timer
	if c.Server.Restart {
		restartTimes, err := getRestartTimes(c)
		if err != nil {
			return err
		}
		// 次の再起動時刻を取得
		nextRestartTime = getNextRestartTime(restartTimes)
		fmt.Printf("次の再起動時刻: %s\n", nextRestartTime.Format("15:04"))
		// 次の再起動時刻になったら再起動
		d := nextRestartTime.Sub(time.Now())
		timer = time.AfterFunc(d, func() {
			fmt.Println("再起動します")
			// stdinにstopを書き込む
			io.WriteString(stdin, "stop\n")
		})
	}

	// サーバーが終了したら終了
	if err := cmd.Wait(); err != nil {
		restartFlag = false
		fmt.Printf("サーバーが異常終了しました\n")
		return err
	}
	if c.Server.Restart {
		timer.Stop()
	}
	return nil
}

// 時刻をパース
func parseTimes(times []string) ([]time.Time, error) {
	var parsedTimes []time.Time
	for _, t := range times {
		parsedTime, err := time.Parse("15:04", t)
		if err != nil {
			return nil, err
		}
		parsedTimes = append(parsedTimes, parsedTime)
	}
	return parsedTimes, nil
}

// 再起動時刻を取得
func getRestartTimes(c *Config) ([]time.Time, error) {
	for i, t := range c.Server.RestartTime {
		fmt.Printf("再起動時刻[%d]: %s\n", i, t)
	}
	return parseTimes(c.Server.RestartTime)
}

// 次の再起動時刻を取得
func getNextRestartTime(restartTimes []time.Time) time.Time {
	now := time.Now()
	current_hour := now.Hour()
	current_minute := now.Minute()

	var closest time.Time

	for _, t := range restartTimes {
		hour := t.Hour()
		minute := t.Minute()
		if hour > current_hour {
			closest = t
			break
		}
		if hour == current_hour && minute > current_minute {
			closest = t
			break
		}
	}
	return closest
}
