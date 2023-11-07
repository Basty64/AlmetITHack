package util

import (
	"fmt"
	"os/exec"
)

func GetTxt(file string) {
	b := exec.Command("pdf_script.py " + file)
	fmt.Println(b.Output())
}
