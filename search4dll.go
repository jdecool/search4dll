package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var dumpbinPath string

func main() {
	if runtime.GOOS != "windows" {
		log.Fatalln("This application can only run on Windows")
	}

	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatalln("Missing arguments: search4dll.exe PATH FILENAME")
	}

	path := flag.Arg(0)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Path '%s' not exists", path)
	}

	dllFilename := flag.Arg(1)
	log.Printf("Looking for '%s' file in '%s'\n", dllFilename, path)

	dumpbinPath = "dumpbin.exe"
	if _, err := os.Stat(dumpbinPath); os.IsNotExist(err) {
		dumpbinPath = "bin/dumpbin.exe"
		if _, err := os.Stat(dumpbinPath); os.IsNotExist(err) {
			log.Fatalln("dumpbin.exe not found")
		}
	}

	start := time.Now()

	err := filepath.Walk(path, searchExecutableFile(dllFilename))
	if err != nil {
		log.Fatalln(err.Error())
	}

	elapsed := time.Since(start)
	log.Printf("Scan took %s", elapsed)
}

func searchExecutableFile(dllFilename string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if f.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".exe" {
			return nil
		}

		searchDll(path, dllFilename)

		return nil
	}
}

func searchDll(file string, dll string) {
	cmdResult, err := exec.Command(dumpbinPath, "/IMPORTS", file).Output()
	if err != nil {
		log.Printf("dumpbin error '%s' on file '%s'\n", err.Error(), file)
		return
	}

	output := string(cmdResult)
	if strings.Contains(output, dll) {
		fmt.Printf("Used in '%s'\n", file)
	}
}
