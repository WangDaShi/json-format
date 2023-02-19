package tool

import (
    "runtime"
    "path"
    "bufio"
    "os"
   "io/ioutil" 
)

func GetProjectFolder() string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        panic("Could not get current file path")
    }
    return path.Dir(path.Dir(filename))
}

func ReadString(file string) string {
    content, err := ioutil.ReadFile(file)
    if err != nil {
        panic(err)
    }
    return string(content)
}

func ReadLine(filePath string) []string {
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var line []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line = append(line,scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
    return line
}
