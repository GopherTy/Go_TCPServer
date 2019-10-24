package cmd

import (
	"log"
	"server/cmd/client"
	"server/configure"
	"server/logger"
	"server/utils"

	"github.com/spf13/cobra"
)

func init() {
	var filename string
	basePath := utils.BasePath()
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run client test",
		Run: func(cmd *cobra.Command, args []string) {
			// load configure
			cnf := configure.Single()
			e := cnf.Load(filename)
			if e != nil {
				log.Fatalln(e)
			}
			e = cnf.Format(basePath)
			if e != nil {
				log.Fatalln(e)
			}
			// init logger
			e = logger.Init(basePath, &cnf.Logger)
			if e != nil {
				log.Fatalln(e)
			}
			// run
			client.Run()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, "config",
		"c",
		utils.Abs(basePath, "server.jsonnet"),
		"configure file",
	)
	rootCmd.AddCommand(cmd)
}
