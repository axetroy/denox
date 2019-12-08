package fs

import "os"

func PathExists(path string) (isExist bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
