// Copyright © 2019 S. van der Baan <steven@vdbaan.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"encoding/hex"
	"os"
)

//
func printOutputS(reply string) string {
	return printOutput([]byte(reply), len(reply))
}

//
func printOutput(reply []byte, length int) string {
	log.Debugf("Processing %d bytes from a total %d",length, len(reply))
	if noHex {
		return fmt.Sprintf("%q",string(reply[:length]))
	} else {
		hasHex := false
		for _,b := range reply[:length] {
			if b < 9 {
				hasHex = true
			} else if b > 13 && b < 32 {
				hasHex = true
			} else if b > 126 && b < 160 {
				hasHex = true
			}
		}
		if hasHex {
			return hex.Dump(reply[:length])
		} else {
			return fmt.Sprintf("%s",string(reply[:length]))
		}
	}
}

func ifErrorMessageStop(err error, message string) {
    if err != nil {
        log.Error(message)
        os.Exit(1)
    }
}

func printBanner() {
    log.Info(programBanner)
}