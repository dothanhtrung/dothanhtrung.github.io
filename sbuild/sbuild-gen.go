/*
 Copyright (C) 2017 Toshiba Corporation

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 2 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful, but
 WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see
 <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	sbuildLog            = "https://raw.githubusercontent.com/dothanhtrung/build-log/master/"
	debian_cross_patches = "https://github.com/meta-debian/debian-cross-patches/tree/master/"
)

func main() {
	statusFile, _ := ioutil.ReadFile("sbuild-status")
	lines := strings.Split(string(statusFile), "\n")

	var remarks map[string]interface{}
	remarkFile, _ := ioutil.ReadFile("remark.json")
	json.Unmarshal(remarkFile, &remarks)

	f, _ := os.Create("index.html")
	defer f.Close()

	html := "<html><head><title>Sbuild Status</title>\n"
	html += "<script src=\"https://kryogenix.org/code/browser/sorttable/sorttable.js\"></script>\n"
	html += "<script src=\"../js/default.js\"></script>\n"
	html += "<link rel=\"stylesheet\" type=\"text/css\" id=\"table_row_counter\" href=\"../css/table_row_counter.css\"/>\n"
	html += "</head><body>\n"
	html += "<h1>Debian cross-build state</h1>\n"
	html += "Build Architecture: amd64<br/>Host Architecture: armhf<br/>##Summary##<br/>\n"
	html += "<br/><input type=\"checkbox\"/ onclick=\"tableRemoveCounter(this);\"> Disable row counter. This helps table sort faster. (Re-enabling will take time)\n"
	html += "<br/><table class=\"sortable\" id=\"sortable\">\n"
	html += "<tr bgcolor=\"#bdc3c7\">" +
		"<th class=\"sorttable_nosort\"></th>" +
		"<th>Source Name</th>" +
		"<th class=\"sorttable_nosort\">Version</th>" +
		"<th width=\"80\">Status</th>" +
		"<th width=\"150\">Build At</th>" +
		"<th>Remark</th></tr>\n"

	successCount := 0

	for i := 0; i < len(lines)-1; i++ {
		pkgInfo := strings.Fields(lines[i])
		name := pkgInfo[0]
		version := pkgInfo[1]
		status := pkgInfo[2]
		t, _ := time.Parse(time.RFC3339, pkgInfo[3])
		timestmp := t.Format("Jan _2, 2006 3:04PM")
		customTimestmp := t.Format("20060102150405")

		bgcolor := "#df2029"
		if status == "attempted" {
			bgcolor = "#e74c3c"
		} else if status == "skipped" {
			bgcolor = "#ecf0f1"
		} else if status == "successful" {
			bgcolor = "#2ecc71"
			successCount++
		} else if status == "given-back" {
			bgcolor = "#e67e22"
		}

		remark := ""
		if remarks[name] != nil {
			remark = remarks[name].(string)
		}
		if _, err := os.Stat("./debian-cross-patches/"+name); !os.IsNotExist(err) {
			if remark != "" {
				remark += "</br>"
			}
			remark += "\n (<a href=\"" + debian_cross_patches + name + "\">debian-cross-patches</a>)"
		}


		logFile := name + "_" + version + "_armhf.build"
		row := "<tr bgcolor=\"" + bgcolor + "\"><td></td>\n" +
			"<td><a href=\"" + sbuildLog + logFile + "\">" + name + "</a></td>\n" +
			"<td>" + version + "</td>\n" +
			"<td>" + status + "</td>\n" +
			"<td sorttable_customkey=\"" + customTimestmp + "\">" + timestmp + "</td>\n" +
			"<td>" + remark + "</td></tr>\n"
		html += row
	}

	html = strings.Replace(html, "##Summary##", "Success: "+strconv.Itoa(successCount)+"<br/>Total: "+strconv.Itoa(len(lines)-1), -1)
	html += "</table></body></html>"

	f.WriteString(html)
	f.Sync()
	w := bufio.NewWriter(f)
	w.Flush()
}
