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
	Arch string `json:"arch"`
	Path string `json:"path"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	root := homeDir + path
	libMap := make(map[string]M)
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		cmd := exec.Command("dwarfdump", "--uuid", path)
		output, err := cmd.Output()
		if err == nil {
			s := strings.SplitN(string(output), " ", 4)
			arch := strings.TrimPrefix(s[2], "(")
			arch = strings.TrimSuffix(arch, ")")
			m := M{
				Arch: arch,
				Path: strings.TrimSuffix(s[3], "\n"),
			}
			libMap[s[1]] = m
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(libMap)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("map.json", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
