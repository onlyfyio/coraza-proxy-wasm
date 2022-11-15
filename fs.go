// Copyright The OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io/fs"
	"strings"
)

type rulesFS struct {
	fs           fs.FS
	filesMapping map[string]string
	dirsMapping  map[string]string
}

func (r rulesFS) Open(name string) (fs.File, error) {
	return r.fs.Open(r.mapPath(name))
}

func (r rulesFS) ReadDir(name string) ([]fs.DirEntry, error) {
	for a, dst := range r.dirsMapping {
		if a == name {
			return fs.ReadDir(r.fs, dst)
		}

		prefix := a + "/"
		if strings.HasPrefix(name, prefix) {
			return fs.ReadDir(r.fs, fmt.Sprintf("%s/%s", dst, name[len(prefix):]))
		}
	}
	return fs.ReadDir(r.fs, name)
}

func (r rulesFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(r.fs, r.mapPath(name))
}

func (r rulesFS) mapPath(p string) string {
	if strings.IndexByte(p, '/') != -1 {
		// is not in root, hence we can do dir mapping
		for a, dst := range r.dirsMapping {
			prefix := a + "/"
			if strings.HasPrefix(p, prefix) {
				return fmt.Sprintf("%s/%s", dst, p[len(prefix):])
			}
		}
	}

	for a, dst := range r.filesMapping {
		if a == p {
			return dst
		}
	}

	return p
}
