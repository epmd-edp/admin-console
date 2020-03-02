/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import "strings"

func Contains(array []string, e string) bool {
	for _, element := range array {
		if element == e {
			return true
		}
	}
	return false
}

func RemoveElByValue(array []string, val string) []string {
	for i, v := range array {
		if strings.ToLower(v) == strings.ToLower(val) {
			return append(array[:i], array[i+1:]...)
		}
	}
	return array
}
