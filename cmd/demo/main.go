package main

import (
	"fmt"

	"github.com/gkarman/demo/internal/config"
)

func main() {
	cfg := config.MustLoad("configs/config.yaml")
	_ = cfg
	fmt.Printf("%+v\n", cfg)
}
