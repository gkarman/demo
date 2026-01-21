package main

import "github.com/gkarman/demo/internal/config"

func main() {
	cfg := config.MustLoad("configs/config.yaml")
	_ = cfg
}
