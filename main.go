package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const path = "/Library/Developer/Xcode/iOS DeviceSupport/"

type M struct {
	UUID string `json:"uuid"`
	Arch string `json:"arch"`
	Path string `json:"path"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	root := homeDir + path
	libs := make([]M, 0, 9000)
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		// 打印文件路径
		cmd := exec.Command("dwarfdump", "--uuid", path)
		output, err := cmd.Output()
		if err == nil {
			s := strings.SplitN(string(output), " ", 4)
			arch := strings.TrimPrefix(s[2], "(")
			arch = strings.TrimSuffix(arch, ")")
			m := M{
				UUID: s[1],
				Arch: arch,
				Path: strings.TrimSuffix(s[3], "\n"),
			}
			libs = append(libs, m)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(libs)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("map.json", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
