package main

import (
	"os"

	"github.com/recruit-tech/RISUCON2022Summer/assets/generator"
)

func main() {
	generator.GenerateIcon(os.Stdout)
}
