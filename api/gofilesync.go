package gofilesync

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var used []string

// ErrUsed is thrown if you attempt to sync
// a file that is already being synced from/to.
var ErrUsed = errors.New("file already being synced from/to")

// OnCopy gets called before any attempt to copy.
// By default logs to console.
var OnCopy = func(src, dst string) {
	fmt.Println("Syncing " + src + " to " + dst)
}

// OnDelete gets called before any attempt to delete a deleted file.
// By default logs to console.
var OnDelete = func(src, dst string) {
	fmt.Println("Deleting " + dst)
}

// SyncData is the sync data file name.
const SyncData = ".syncData"

// ReadSyncData reads the .syncData file in dir
func ReadSyncData(dir string) (map[string]int64, error) {
	f, err := os.Open(filepath.Join(dir, SyncData))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make(map[string]int64)
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// WriteSyncData writes the .syncData file in dir
func WriteSyncData(dir string, data map[string]int64) error {
	f, err := os.Create(filepath.Join(dir, SyncData))
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	return encoder.Encode(data)
}

func clonefile(src, dst string, info os.FileInfo) error {
	OnCopy(src, dst)

	if link, err := os.Readlink(src); err == nil {
		os.Symlink(link, dst)
	} else {
		from, err := os.Open(src)
		if err != nil {
			return err
		}
		defer from.Close()

		to, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, info.Mode())
		if err != nil {
			return err
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			return err
		}
	}
	return nil
}

func use(dirs ...string) {
	used = append(used, dirs...)
}
func isUsed(dirs ...string) bool {
	for _, dir := range used {
		for _, dir2 := range dirs {
			if dir == dir2 {
				return true
			}
		}
	}
	return false
}
func unuse(dirs ...string) {
	var used2 []string
	for _, dir := range used {
		contains := false
		for _, dir2 := range dirs {
			if dir == dir2 {
				contains = true
			}
		}

		if !contains {
			used2 = append(used2, dir)
		}
	}
	used = used2
}
