package main

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"turtle/blueprints"
	"turtle/core/usersApi"
	"turtle/docmining"
	"turtle/themes"
	"turtle/toolsApi"

	"turtle/crm"
	"turtle/security"
	"turtle/turtleio"

	"turtle/agentTools"
	"turtle/api"
	"turtle/apiApp"
	"turtle/auth"
	"turtle/credentials"
	"turtle/db"
	"turtle/documents"
	"turtle/flows"
	"turtle/fn"
	"turtle/forecasting"
	"turtle/knowledgeHub"
	"turtle/llm"
	"turtle/llm/llmApi"
	"turtle/llm/llmCtrl"
	"turtle/nn"
	"turtle/tables"
	"turtle/tags"

	"turtle/lgr"
	"turtle/models"
	"turtle/server"
	"turtle/vfs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite" //Must be kvoli registracii drivera na sqlite3
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	return r
}

func _MainRoute(c *gin.Context) {

	version := Version
	buildTime := BuildTime

	c.JSON(http.StatusOK, gin.H{
		"version":    version,
		"build_time": buildTime,
	})
}

func use_cors(r *gin.Engine) {

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://tauri.localhost", "*"}
	r.Use(cors.New(config))

}
func dev_main() {

	db.DB.InstallJavaScriptFunctions()

	lgr.Ok("Working directory: ", vfs.GetWorkingDirectory())

	r := setupRouter()
	use_cors(r)

	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".xyz", "text/plain")
	mime.AddExtensionType(".gzip", "application/x-gzip-compressed")
	mime.AddExtensionType(".gz", "application/x-gzip-compressed")

	go func() {
		defer func() {
			llmCtrl.InitOllama()
		}()
	}()

	// r.Use()

	r.GET("/api/main", _MainRoute)

	agentTools.InitCoreTools()

	agentTools.InitAgentToolsApi(r)
	documents.InitDocumentsApi(r)
	documents.InitDocumentCollectionsApi(r)
	api.InitApi(r)
	apiApp.InitApiApp(r)
	auth.Init_api_auth0(r)
	nn.InitNNApi(r)
	fn.InitFnApi(r)
	knowledgeHub.InitKnowledgeHubApi(r)
	tags.InitTagsApi(r)
	forecasting.InitForecastingApi(r)
	crm.InitCrmApi(r)
	security.InitSecurityApi(r)

	usersApi.InitUsersApi(r)

	tables.InitTablesApi(r)

	llm.InitIncidentsApi(r)
	llmApi.InitLLMApi(r)
	blueprints.InitBlueprintsApi(r)
	flows.InitFlowsApi(r)

	themes.InitThemesApi(r)

	docmining.InitDocMiningApi(r)

	toolsApi.InitToolsApi(r)

	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	//r.NoRoute(tools.ProxyMiddleware2())

	prefix := "http://"
	port := "8080"

	if credentials.RunHttps() {
		prefix = "https://"
	}

	addr := "0.0.0.0:" + port

	server.RunMyioServer(r)

	turtleio.InitTurtleSocketsApi(r)

	lgr.Info("Running in: ", vfs.GetExeFile())

	srv := &http.Server{
		Addr:         addr, // Change to your desired port
		Handler:      r,
		ReadTimeout:  1000 * time.Second, // StringSet read timeout
		WriteTimeout: 1000 * time.Second, // StringSet write timeout
		IdleTimeout:  1000 * time.Second, // StringSet idle timeout
	}

	if credentials.RunHttps() {
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	lgr.Ok("------------------------------------")
	lgr.Ok("Mode: %s", gin.Mode())

	lgr.Ok("Running server at: %s", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+addr, prefix+addr))
	lgr.Ok("Access at: %s", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+"127.0.0.1:"+port, prefix+"127.0.0.1:"+port))
	lgr.Ok("Access at: %s", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+"localhost:"+port, prefix+"localhost:"+port))

	if vfs.IsLinux() {
		cmd := exec.Command("hostname", "-I")

		// Capture the output
		output, err := cmd.Output()
		if err != nil {
			lgr.Error(err.Error())

		} else {
			ips := strings.Fields(string(output))

			for _, ip := range ips {
				lgr.Ok("Access at: ", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+ip+":"+port, prefix+ip+":"+port))

			}
		}

	}

	if credentials.RunHttps() {
		lgr.Ok("Started HTTP(S) branch")
		error := srv.ListenAndServeTLS("cert.pem", "key.pem")

		if error != nil {
			lgr.ErrorJson(error)
		}

	} else {
		lgr.Ok("Started HTTP branch")
		error := srv.ListenAndServe()

		if error != nil {
			lgr.ErrorJson(error)
		}
	}

	lgr.Error("Execution ended")
}

func main() {

	//levenstein.TestLevenstein()

	models.RegisterClazzFactory()

	//db.InitGorm()

	lgr.Info("Starting infinity twin application")
	lgr.Info("DbName: ", credentials.GetDBName())

	dev_main()
	//TestMain()
}

var (
	// These variables will be set during the build process
	Version   string
	BuildTime string
)
