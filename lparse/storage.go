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
    "redis"
    "math"
    "os"
)

type StorageController struct {
    Redis   redis.Client
    Unit    *StorageUnit
}

type StorageUnit struct {
    Index   map[string] []bool
    Buffer  []string
}

func NewStorageController() *StorageController {
    controller := new(StorageController)
    controller.Unit = NewStorageUnit()

    controller.Redis.Addr = "127.0.0.1:6379"
    controller.Redis.Db = 1

    return controller
}

func (sc *StorageController) Store(key string) os.Error {
    for index, bin_array := range sc.Unit.Index {
        _, e := sc.Redis.Zadd(key, []byte(index), binToInt(bin_array))
        if e != nil {
            return e
        }
    }

    for i := 0; i < len(sc.Unit.Buffer); i++ {
        sc.Redis.Rpush(key + ":result", []byte(sc.Unit.Buffer[i]))
    }

    sc.Unit = NewStorageUnit()

    return nil
}

func NewStorageUnit() *StorageUnit {
    unit := new(StorageUnit)
    
    unit.Index = make(map[string] []bool)

    return unit
}

func (s *StorageUnit) AddLine(line string, index []string) {
    s.Buffer = append(s.Buffer, line)
    line_index := len(s.Buffer) - 1

    for i := 0; i < len(index); i++ {
        _, present := s.Index[index[i]]

        if present {
            s.Index[index[i]] = binAdd(line_index, s.Index[index[i]])
        } else {
            s.Index[index[i]] = binEncode(line_index)
        }
    }
} 

func binEncode(index int) []bool {
    result := make([]bool, index + 1)

    for i := 0; i < index; i++ {
        result[i] = false
    }
    result[index] = true
    
    return result
}

func binAdd(index int, bin_array []bool) []bool {
    for len(bin_array) < index {
        bin_array = append(bin_array, false)
    }
    bin_array = append(bin_array, true)

    return bin_array
}

func binToInt(bin_array []bool) float64 {
    result := 0.0

    for i := 0; i < len(bin_array); i++ {
        if bin_array[i] {
            result += math.Pow(2.0, float64(i))
        }
    }

    return result
}
