package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dataRepo    = "Kengxxiao/ArknightsGameData"
	dataBranch  = "master"
	dataFolder  = "zh_CN/gamedata/excel"
	avatarRepo  = "fexli/ArknightsResource"
	avatarBranch = "main"
	avatarFolder = "avatar/ASSISTANT"
)

var dataFiles = []string{
	"character_table.json",
	"skill_table.json",
	"building_data.json",
	"char_meta_table.json",
	"range_table.json",
	"release_dates.json",
}

// CheckResult holds the result of an update check.
type CheckResult struct {
	HasUpdate bool
	Message   string
}

// CheckUpdate compares local and remote data versions.
func CheckUpdate(dataRoot string) (*CheckResult, error) {
	local := readVersionRaw(dataRoot)
	remote, err := httpGetString(dataURL("data_version.txt"))
	if err != nil {
		return &CheckResult{HasUpdate: false, Message: "无法连接更新服务器"}, nil
	}
	remote = strings.TrimSpace(remote)

	if local == "" {
		return &CheckResult{HasUpdate: true, Message: "本地还没有数据，点击更新开始首次下载"}, nil
	}
	if local == remote {
		return &CheckResult{HasUpdate: false, Message: "已是最新版本"}, nil
	}
	return &CheckResult{HasUpdate: true, Message: "发现新版本，点击更新下载"}, nil
}

// Update downloads game data files from GitHub and removes the cache.
func Update(dataRoot string, progress func(string)) error {
	gamedataDir := filepath.Join(dataRoot, "gamedata")
	avatarDir := filepath.Join(dataRoot, "avatars")
	os.MkdirAll(gamedataDir, 0755)
	os.MkdirAll(avatarDir, 0755)

	for i, name := range dataFiles {
		progress(fmt.Sprintf("下载数据 %d/%d %s %d%%", i+1, len(dataFiles), name, (i+1)*15/len(dataFiles)))
		tmp := filepath.Join(gamedataDir, name+".tmp")
		target := filepath.Join(gamedataDir, name)
		if err := httpGetToFile(dataURL(name), tmp); err != nil {
			return fmt.Errorf("下载 %s 失败: %w", name, err)
		}
		if err := os.Rename(tmp, target); err != nil {
			return err
		}
	}

	// Download version file last (marks successful update)
	versionTmp := filepath.Join(gamedataDir, "data_version.txt.tmp")
	versionTarget := filepath.Join(gamedataDir, "data_version.txt")
	if err := httpGetToFile(dataURL("data_version.txt"), versionTmp); err != nil {
		return fmt.Errorf("下载版本文件失败: %w", err)
	}
	os.Rename(versionTmp, versionTarget)

	// Remove operator cache so it gets regenerated
	cachePath := filepath.Join(dataRoot, cacheFile)
	os.Remove(cachePath)

	// Download missing avatars
	ids, _ := ReadOperatorIDs(filepath.Join(gamedataDir, "character_table.json"))
	missing := 0
	for _, id := range ids {
		if _, err := os.Stat(filepath.Join(avatarDir, id+".png")); os.IsNotExist(err) {
			missing++
		}
	}
	if missing > 0 {
		progress(fmt.Sprintf("正在下载头像（共 %d 名干员，缺 %d 个）…", len(ids), missing))
		DownloadAvatars(ids, avatarDir, progress)
	}

	progress("下载完成 100%")
	return nil
}

func dataURL(file string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", dataRepo, dataBranch, dataFolder, file)
}

func readVersionRaw(dataRoot string) string {
	path := filepath.Join(dataRoot, "gamedata", "data_version.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

var httpClient = &http.Client{Timeout: 5 * time.Minute}

func httpGetString(url string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func httpGetToFile(url string, path string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// ReadOperatorIDs extracts all operator IDs from character_table.json (for avatar downloads).
func ReadOperatorIDs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data map[string]json.RawMessage
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(data))
	for id := range data {
		ids = append(ids, id)
	}
	return ids, nil
}

// DownloadAvatars downloads missing operator avatars from fexli/ArknightsResource.
func DownloadAvatars(ids []string, avatarDir string, progress func(string)) {
	done := 0
	missing := 0
	for _, id := range ids {
		target := filepath.Join(avatarDir, id+".png")
		if _, err := os.Stat(target); err == nil {
			continue
		}
		missing++
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s.png", avatarRepo, avatarBranch, avatarFolder, id)
		if err := httpGetToFile(url, target); err != nil {
			continue // skip missing avatars silently
		}
		done++
		if done%50 == 0 || done == missing {
			progress(fmt.Sprintf("头像下载进度：%d/%d", done, missing))
		}
	}
	if done > 0 {
		progress(fmt.Sprintf("头像下载完成：%d 个", done))
	}
}
