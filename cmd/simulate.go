/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"lucy/simulate"

	"github.com/spf13/cobra"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulate called", args)

		simulate.NewService().Handle(simulate.Input{
			Method: Method,
			Count:  Count,
		})
	},
}

var (
	Method int
	Count  int
)

func init() {
	simulateCmd.Flags().IntVarP(&Count, "count", "c", 1, "simulate user count, default 1")
	simulateCmd.Flags().IntVarP(&Method, "method", "m", simulate.StressMoveMethod, "1-多人移动压测")

	rootCmd.AddCommand(simulateCmd)
}
