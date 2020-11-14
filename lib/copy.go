package redux

import (
	"log"
	"strings"

	cp "github.com/otiai10/copy"
)

// CopyDir moves over the dir and skips git dirs
func CopyDir(src, dst string) (err error) {
	err = cp.Copy(src, dst,
		cp.Options{
			Skip: func(src string) (bool, error) {
				isGit := strings.Contains(src, "/.git")
				if isGit {
					log.Printf("Skip: %s", src)
				}
				return isGit, nil
			},
		},
	)
	return
}
