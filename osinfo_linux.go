/* Copyright 2016 Miguel Moll
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

package osinfo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

var uname = "/usr/bin/uname"
var osReleaseFile string
var distroInfo = make(map[string]string)
var releaseFile string

func init() {
	releaseFiles := []string{}

	files, _ := ioutil.ReadDir("/etc")
	for _, f := range files {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), "-release") {
				releaseFiles = append(releaseFiles, f.Name())
			}

			if f.Name() == "os-release" {
				osReleaseFile = "/etc/os-release"
			}

		}
	}

	if osReleaseFile != "" {
		parseOSRelease()
	}

}

func parseOSRelease() {

	bytes, err := ioutil.ReadFile(osReleaseFile)
	if err != nil {
		log.Println("Unable to parse os release file")
		return
	}

	dataLines := strings.Split(string(bytes), "\n")
	for _, line := range dataLines {
		if line == "" {
			continue
		}

		entries := strings.Split(line, "=")
		distroInfo[entries[0]] = entries[1]
	}
}

func platform() string {
	_, err := exec.Command(uname).Output()
	if err != nil {
		return ""
	}

	return runtime.GOOS
}

func kernel() string {
	out, err := exec.Command(uname, "-r").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

func arch() string {
	out, err := exec.Command(uname, "-p").Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

func distribution() string {
	if val, ok := distroInfo["NAME"]; ok {
		return val
	}

	return ""
}

func version() string {
	if val, ok := distroInfo["VERSION_ID"]; ok {
		return val
	}

	return ""
}

func pretty() string {
	if val, ok := distroInfo["PRETTY_NAME"]; ok {
		return strings.TrimSuffix(strings.TrimPrefix(val, "\""), "\"")
	}

	return ""
}
