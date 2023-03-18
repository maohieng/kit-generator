package cmd

import (
	"github.com/GrantZheng/kit/generator"
	"github.com/GrantZheng/kit/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var methods []string
var initserviceCmd = &cobra.Command{
	Use:     "service",
	Short:   "Initiate a service",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide a name for the service")
			return
		}
		if viper.GetString("g_s_transport") == "grpc" {
			if !checkProtoc() {
				return
			}
		}
		// pbPath,pbImportPath only work when transport is grpc
		var pbPath, pbImportPath string
		pbPath = viper.GetString("g_s_pb_path")
		pbImportPath = viper.GetString("g_s_pb_import_path")
		if pbPath != "" {
			exist := utils.IsExist(pbPath)
			if !exist {
				logrus.Errorf("You must provide a existed pb path to put pb dir, given path:<%s> does not exist", pbPath)
				return
			}
			// The pbImportPath validity is not checked here, and it is the responsibility of the developer
			if pbImportPath == "" {
				logrus.Error("You must provide pb import path by --pb_import_path or -i, because you have provided pb_path")
				return
			}
		}

		var emw, smw bool
		if viper.GetBool("g_s_dmw") {
			emw = true
			smw = true
		} else {
			emw = viper.GetBool("g_s_endpoint_mdw")
			smw = viper.GetBool("g_s_svc_mdw")
		}
		var routerName string
		if viper.GetBool("g_s_fastrouter") {
			routerName = generator.FastRouter
		} else if viper.GetBool("g_s_httprouter") {
			routerName = generator.HttpRouter
		} else if viper.GetBool("g_s_gorilla") {
			routerName = generator.Gorilla
		}
		g := generator.NewGenerateService(
			args[0],
			viper.GetString("g_s_transport"),
			routerName,
			pbPath,
			pbImportPath,
			smw,
			emw,
			methods,
		)
		if err := g.Generate(); err != nil {
			logrus.Error(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(initserviceCmd)
	initserviceCmd.Flags().StringP("transport", "t", "http", "The transport you want your service to be initiated with")
	initserviceCmd.Flags().StringP("pb_path", "p", "", "Specify path to store pb dir")
	initserviceCmd.Flags().StringP("pb_import_path", "i", "", "Specify path to import pb")
	initserviceCmd.Flags().BoolP("dmw", "w", false, "Generate default middleware for service and endpoint")
	initserviceCmd.Flags().Bool(generator.Gorilla, false, "Generate http using gorilla mux")
	initserviceCmd.Flags().Bool(generator.HttpRouter, false, "Generate http using julienschmidt/httprouter")
	initserviceCmd.Flags().Bool(generator.FastRouter, false, "Generate http using fasthttp/router")
	initserviceCmd.Flags().StringArrayVarP(&methods, "methods", "m", []string{}, "Specify methods to be generated")
	initserviceCmd.Flags().Bool("svc-mdw", false, "If set a default Logging and Instrumental middleware will be created and attached to the service")
	initserviceCmd.Flags().Bool("endpoint-mdw", false, "If set a default Logging and Tracking middleware will be created and attached to the endpoint")
	viper.BindPFlag("g_s_transport", initserviceCmd.Flags().Lookup("transport"))
	viper.BindPFlag("g_s_pb_path", initserviceCmd.Flags().Lookup("pb_path"))
	viper.BindPFlag("g_s_pb_import_path", initserviceCmd.Flags().Lookup("pb_import_path"))
	viper.BindPFlag("g_s_dmw", initserviceCmd.Flags().Lookup("dmw"))
	viper.BindPFlag("g_s_gorilla", initserviceCmd.Flags().Lookup(generator.Gorilla))
	viper.BindPFlag("g_s_httprouter", initserviceCmd.Flags().Lookup(generator.HttpRouter))
	viper.BindPFlag("g_s_fastrouter", initserviceCmd.Flags().Lookup(generator.FastRouter))
	viper.BindPFlag("g_s_svc_mdw", initserviceCmd.Flags().Lookup("svc-mdw"))
	viper.BindPFlag("g_s_endpoint_mdw", initserviceCmd.Flags().Lookup("endpoint-mdw"))
}
