package gofilesync

import (
	"os"
	"path/filepath"
)

// ForceSync will completely delete dst,
// and re-clone from src.
func ForceSync(src string, dst string) error {
	if isUsed(src, dst) {
		return ErrUsed
	}
	use(src, dst)
	defer unuse(src, dst)

	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Faster and more efficient to just clone directly
	if !info.IsDir() {
		err := os.Remove(dst)
		if err != nil {
			return err
		}
		err = clonefile(src, dst, info)
		if err != nil {
			return err
		}
		OnSuccess()
		return nil
	}

	err = os.RemoveAll(dst)
	if err != nil {
		return err
	}

	data := make(map[string]int64)
	err = clone(src, dst, info, data)
	if err != nil {
		return err
	}
	err = WriteSyncData(dst, data)
	if err != nil {
		return err
	}

	OnSuccess()
	return nil
}

func clone(src, dst string, info os.FileInfo, data map[string]int64) error {
	if info.Name() == SyncData {
		return nil
	}

	if info.IsDir() {
		err := os.Mkdir(dst, 0755)
		if err != nil {
			return err
		}

		dir, err := os.Open(src)
		if err != nil {
			return err
		}
		files, err := dir.Readdir(-1)
		dir.Close()
		if err != nil {
			return err
		}

		for _, file := range files {
			pathSrc := filepath.Join(src, file.Name())
			pathDst := filepath.Join(dst, file.Name())
			err := clone(pathSrc, pathDst, file, data)
			if err != nil {
				return err
			}
		}
	} else {
		data[dst] = info.ModTime().UnixNano()
		clonefile(src, dst, info)
	}
	return nil
}
