/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package archiver

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CreateApplicationArchive(folder string) (string, error) {
	tarball, err := ioutil.TempFile(os.TempDir(), "blob")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer tarball.Close()
	gz := gzip.NewWriter(tarball)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	if _, err := os.Stat(filepath.Join(folder, "run.sh")); os.IsNotExist(err) {
		fmt.Println("run.sh does not exist")
		fmt.Println("Create a script with commands how to install required dependencies offline and run your application.")
		return "", err
	}

	err = filepath.Walk(folder, walkAndCompress(folder, tw))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tarball.Name(), nil
}

func walkAndCompress(baseDir string, tw *tar.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == baseDir {
			return nil
		}

		relativePath := strings.TrimPrefix(path, baseDir+"/")

		if (info.Mode() & os.ModeSymlink) == 0 {
			err = saveFileToTar(tw, info, relativePath)
		} else {
			err = saveSymlinkToTar(tw, info, relativePath)
		}

		if err != nil {
			return err
		}

		fmt.Printf("Added to archive: %v\n", relativePath)
		return nil
	}
}

func saveFileToTar(tw *tar.Writer, info os.FileInfo, path string) error {
	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return err
	}
	header.Name = path

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(tw, file)

	return err
}

func saveSymlinkToTar(tw *tar.Writer, info os.FileInfo, path string) error {
	symlink, err := os.Readlink(path)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, symlink)
	if err != nil {
		return err
	}
	header.Name = path

	return tw.WriteHeader(header)
}
