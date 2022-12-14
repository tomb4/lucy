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
	Short: "simulate agent move",
	Long:  `simulate agent move`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start simulating", Method)

		simulate.NewService(Cfg, Count).Handle(Method)
	},
}

var (
	Method int
	Count  int
	Cfg    string
)

func init() {
	simulateCmd.Flags().IntVarP(&Count, "count", "c", 1, "simulate agent count")
	simulateCmd.Flags().IntVarP(&Method, "method", "m", simulate.StressMoveMethod, "1-Multi-person mobile pressure measurement")
	simulateCmd.Flags().StringVarP(&Cfg, "cfg", "f", "./simulate/config/test.yml", "config file")

	rootCmd.AddCommand(simulateCmd)
}
