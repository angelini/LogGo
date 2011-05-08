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
 *	Live parser, this object contains all the information and methods
 *  for the files being parsed
 *
 */

package lparse

import (
    "fmt"
    "strings"
    "os"
    "os/inotify"
    "web"
)

type Lparse struct {
    Files   []*MonitFile
    Watcher *inotify.Watcher
    Storage *StorageController
    Server  *web.Server
}

var NonIndex = map[string] bool {
    " ": true,    
    ",": true,
    "'": true,
    ":": true,
    "=": true,
}

var mainParser *Lparse

func NewLparse() (*Lparse, os.Error) {
    parser := new(Lparse)

    watch, e := inotify.NewWatcher() 
    if e != nil {
        return nil, e
    }

    server, e := StartServer(8080) 
    if e != nil {
        return nil, e
    }

    parser.Watcher = watch
    parser.Storage = NewStorageController()
    parser.Server = server    

    mainParser = parser

    return parser, nil
}

func (l *Lparse) AddFile(filePath string) os.Error {
    monit, e := NewMonitFile(filePath)
    if e != nil {
        return e
    }

    _, e = monit.Read()

    l.Files = append(l.Files, monit)
    l.Watcher.Watch(monit.File.Name())

    return nil
}

func (l *Lparse) RemoveFile(filePath string) os.Error {
    for i := 0; i < len(l.Files); i++ {
        if l.Files[i].File.Name() == filePath {
            l.Files = append(l.Files[:i], l.Files[i+1:]...)
            
            return nil
        } 
    }

    return os.NewError("File not found in Parser Object")
}

func (l *Lparse) Parse(name string) os.Error{
    var monit *MonitFile    

    for i := 0; i < len(l.Files); i++ {
        if l.Files[i].File.Name() == name {
            monit = l.Files[i]
        }
    }

    if monit == nil {
        return os.NewError("Event caught for unknown File")
    }

    lines, e := monit.Read()
    if e != nil {
        return e
    }

    for i := 0; i < len(lines); i++ {
        index := l.Index(lines[i])
        l.Storage.Unit.AddLine(lines[i], index)
        fmt.Printf("Indices for Line: %s\n Are %v\n", lines[i], index)
    }

    return nil
}

func (l *Lparse) Index(line string) []string {
    return strings.FieldsFunc(line, testIndex) 
}

func testIndex(rune int) bool {
    str := string(rune)
    _, present := NonIndex[str]

    if present {
        return true
    }

    return false
}
