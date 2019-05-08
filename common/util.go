package common

import (
	"log"
	"os"
	"path/filepath"
)

func ExecutableDir() string {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Printf("util: find executableDir err: %v", err)
		return ""
	}
	return filepath.Dir(pathAbs)
}


func CheckFileExist(filename string) bool {
	if _,err:=os.Stat(filename);os.IsNotExist(err){
		return  false
	}
	return  true
}