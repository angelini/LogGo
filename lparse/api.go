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
 *	Web API, to search the storage backend
 *
 */

package lparse

import (
    "web"
    "os"
    "json"
    "strconv"
    "fmt"
)

func StartServer(port int) (*web.Server, os.Error) {
    server := new(web.Server)

    wd, err := os.Getwd()
    if err != nil {
        return nil, err
    }
    
    server.Config = new(web.ServerConfig)
    server.Config.StaticDir = wd + "/web/"
    server.Post("/(.*)", pageRouter)

    server.Run("0.0.0.0:" + strconv.Itoa(port))

    return server, nil
}

func pageRouter(ctx *web.Context, route string) {
    switch route {
    case "search":
        json, e := Search(ctx)
        if e != nil {
            ctx.WriteString(e.String())
        }
        ctx.Write(json)
    default:
        
    }
}

func Search(ctx *web.Context) ([]byte, os.Error) {
    var result []string    
    var key int64
    var err os.Error

    int_key, present := ctx.Request.Params["current"]
    if present {
        key, err = strconv.Atoi64(int_key)
        if err != nil {
            return nil, err
        }
    } else {
        if time, _, e := os.Time(); e != nil {
            return nil, e
        } else {
            key = time / 100
        }
    }

    member := []byte(ctx.Request.Params["search"])
    for len(result) < 25 {
        score, e := mainParser.Storage.Redis.Zscore(strconv.Itoa64(key), member)

        if e == nil {
            if res_array, err := getStrings(score, strconv.Itoa64(key)); err != nil {
                return nil, err
            } else {
                result = append(result, res_array...)
            }
        } else {
            _, e := mainParser.Storage.Redis.Exists(strconv.Itoa64(key))
            if e != nil {
                break
            }
        }

        key--
    }

    return json.Marshal(result)       
}

func getStrings(score float64, key string) ([]string, os.Error) {
    var bin_array []bool
    var results []string

    // Int to Binary Conversion
    int_score := int(score)
    for int_score > 0 {
        if remainder := int_score % 2; remainder != 0 {
            bin_array = append(bin_array, true)    
        } else {
            bin_array = append(bin_array, false)
        }
    
        int_score = int_score / 2    
    }

    for i := 0; i < len(bin_array); i++ {
        if bin_array[i] {
            str, err := mainParser.Storage.Redis.Lrange(key + ":result", i, i)
            if err != nil {
                return nil, err
            }

            results = append(results, string(str[0]))
        }
    }

    return results, nil
}
