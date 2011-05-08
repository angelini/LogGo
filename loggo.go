/**
 *
 *  LogGo
 *  Author, alex.louis.angelini@gmail.com
 *
 *  Copyright (c) Alex Angelini <alex.louis.angelini@gmail.com>
 *  View LICENSE
 *
 */

/**
 *
 *  Loggo set up, this file is in charge of parsing the config
 *  and properly setting up the work server
 *
 */

package main

import (
    "lparse"
    "os"
    "fmt"
    "time"
    "strconv"
)

func main() {
    parser, e := lparse.NewLparse()
    if e != nil {
        fmt.Printf("Cannot create new Parser: %v\n", e)
    }

    e = parser.AddFile("./logs/test.log")
    if e != nil {
        fmt.Printf("Cannot add file: %v\n", e)
    }

    go StorageLoop(parser)

    // Event.Mask == 2, File was modified
    for {
    select {
    case ev := <-parser.Watcher.Event:
            if ev.Mask == 2 {
                parser.Parse(ev.Name)
            }
    case err := <-parser.Watcher.Error:
      fmt.Printf("Watcher error caught: %s\n", err)
    }
  }
}

func StorageLoop(parser *lparse.Lparse) {
    ticker := time.NewTicker(50e9)
    key := int64(0)
    for {
        <-ticker.C

        time, _, e := os.Time()
        if e != nil {
            fmt.Printf("Time Error: %v\n", e)
        }

        if (time / 100) != key {           
            key = (time / 100)

            fmt.Printf("Store ran with key: %v\n", key)
            fmt.Printf("Index: %v\n", parser.Storage.Unit.Index)
            fmt.Printf("Buffer: %v\n", parser.Storage.Unit.Buffer)

            parser.Storage.Store(strconv.Itoa64(key))

            fmt.Printf("Left: %v\n", parser.Storage.Unit.Buffer)
        }
    }
}
