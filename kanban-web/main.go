package main

import (
	"os"
	"strings"

	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/server"
	"irpl.com/kanban-web/services"
)

var (
	Version string
	Build   string
)

const DefaultRestHost string = "0.0.0.0"     // Default port if not set in env
const DefaultRestPort string = "4200"        // Default port if not set in env
const DefaultDBHelperHost string = "0.0.0.0" // Default port if not set in env
const DefaultDBHelperPort string = "4100"    // Default port if not set in env

func init() {
	// Initialize dbHelperHost and dbHelperPort with value from environment variable or fallback to default
	utils.RestHost = os.Getenv("RESTSRV_HOST")
	if strings.TrimSpace(utils.RestHost) == "" {
		utils.RestHost = DefaultRestHost
	}

	utils.RestPort = os.Getenv("RESTSRV_PORT")
	if strings.TrimSpace(utils.RestPort) == "" {
		utils.RestPort = DefaultRestPort
	}

	utils.DBHelperHost = os.Getenv("DBHELPER_HOST")
	if strings.TrimSpace(utils.DBHelperHost) == "" {
		utils.DBHelperHost = DefaultDBHelperHost
	}

	utils.DBHelperPort = os.Getenv("DBHELPER_PORT")
	if strings.TrimSpace(utils.DBHelperPort) == "" {
		utils.DBHelperPort = DefaultDBHelperPort
	}

	utils.DBURL = utils.JoinStr("http://", utils.DBHelperHost, ":", utils.DBHelperPort)

	utils.RestURL = utils.JoinStr("http://", utils.RestHost, ":", utils.RestPort)
}

func main() {

	utils.Version = utils.SetVersion(Version, Build)

	utils.WaitForHTTPServer(utils.DBURL)
	utils.WaitForHTTPServer(utils.RestURL)

	utils.DefaultsMap, _ = services.GetAllDefaultsHandler("000001")

	server.Web()

}