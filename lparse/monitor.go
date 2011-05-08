/**
 *
 *	LogGo
 *	Author, alex.louis.angelini@gmail.com
 *
 *	Copyright (c) Alex Angelini <alex.louis.angelini@gmail.com>
 *	View LICENSE
 *
 */

/**
 *
 *	Monitor functions, contains all the functions which have to do with reading
 *	files and parsing them
 *
 */

package lparse

import (
	"os"
	"strings"
)

type MonitFile struct {
	File     *os.File
	Count    int64
}

func NewMonitFile(filePath string) (*MonitFile, os.Error) {
    monit := new(MonitFile)

    file, e := os.Open(filePath)
    if e != nil {
        return nil, e
    }

    monit.File = file
    monit.Count = 0

    return monit, nil
}

// Reads any new lines from the file, an if toParse is true sends them to the regex
// functions then to the DB or else just ignores them, the ignoring is used when
// setting up files

func (monit *MonitFile) Read() ([]string, os.Error) {
	const BUF = 2056
	var buf [BUF]byte
	var lines []byte

    file_info, e := monit.File.Stat()
    if e != nil {
        return nil, e
    }

    if file_info.Size < monit.Count {
        monit.Count = 0
    }

	for {
		n, _ := monit.File.ReadAt(buf[:], monit.Count)

		lines = append(lines, buf[:n]...)

		monit.Count += int64(n)
		if n != BUF {
			break
		}
	}

	stringArr := strings.Split(string(lines[:]), "\n", -1)
	return stringArr[:(len(stringArr) - 1)], nil
}
