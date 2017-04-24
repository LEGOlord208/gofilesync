package gofilesync

import (
	"os"
	"path/filepath"
)

// LazySync will check for differences between src and dst,
// and update dst as necessary.
func LazySync(src string, dst string) error {
	if isUsed(src, dst) {
		return ErrUsed
	}
	use(src, dst)
	defer unuse(src, dst)

	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Easier and cleaner to just clone the file entirely.
	// I guess I am LazySync is LAZY (ba dum tss)
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

	data, err := ReadSyncData(dst)
	if err != nil {
		if os.IsNotExist(err) {
			data = make(map[string]int64)
		} else {
			return err
		}
	}
	err = lazy(src, dst, info, data)
	if err != nil {
		return err
	}
	err = cleanup(src, dst, info, data)
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

func lazy(src, dst string, info os.FileInfo, data map[string]int64) error {
	if info.Name() == SyncData {
		return nil
	}

	if info.IsDir() {
		err := os.Mkdir(dst, 0755)
		if err != nil && !os.IsExist(err) {
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
			err := lazy(pathSrc, pathDst, file, data)
			if err != nil {
				return err
			}
		}
	} else {
		t, ok := data[dst]
		modtime := info.ModTime().UnixNano()

		if !ok || t != modtime {
			data[dst] = modtime
			clonefile(src, dst, info)
		} else {
			_, err := os.Lstat(dst)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}

				data[dst] = modtime
				clonefile(src, dst, info)
			}
		}
	}
	return nil
}
func cleanup(src, dst string, info os.FileInfo, data map[string]int64) error {
	if info.Name() == SyncData {
		return nil
	}

	if info.IsDir() {
		dir, err := os.Open(dst)
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
			err := cleanup(pathSrc, pathDst, file, data)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := os.Lstat(src)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			OnDelete(src, dst)

			delete(data, dst)
			err = os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
