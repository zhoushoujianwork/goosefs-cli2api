/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"goosefs-cli2api/cmd"
)

var version string

// @title           GooseFS-CLI2API
// @version         v1
// @termsOfService  http://swagger.io/terms/
// @host            localhost:8080
// @BasePath
func main() {
	cmd.Execute(version)
}
