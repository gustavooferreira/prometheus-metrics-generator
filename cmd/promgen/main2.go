package main

// import (
// 	"errors"
// 	"fmt"
// 	"os"
//
// 	"github.com/pterm/pterm"
//
// 	"github.com/gustavooferreira/prometheus-metrics-generator/cmd/promgen/cli"
// )
//
// func main() {
// 	err := cli.NewRootCmd().Execute()
// 	if err != nil {
// 		if errors.Is(err, cli.ErrValidation) {
// 			os.Exit(1)
// 		} else if errors.Is(err, cli.ErrProgram) {
// 			os.Exit(2)
// 		}
//
// 		msg := pterm.Error.Sprintfln("%s", err)
// 		_, _ = fmt.Fprint(os.Stderr, msg)
// 		os.Exit(128)
// 	}
// }
