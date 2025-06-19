package main

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"turtle/api"
	"turtle/apiApp"
	"turtle/auth"
	"turtle/credentials"
	"turtle/db"
	"turtle/documents"
	"turtle/llm/llmApi"
	"turtle/llm/llmCtrl"

	"turtle/lg"
	"turtle/models"
	"turtle/server"
	"turtle/tools"
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

	lg.LogOk("Working directory: ", vfs.GetWorkingDirectory())

	r := setupRouter()
	use_cors(r)

	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".xyz", "text/plain")
	mime.AddExtensionType(".gzip", "application/x-gzip-compressed")
	mime.AddExtensionType(".gz", "application/x-gzip-compressed")

	llmCtrl.InitOllama()

	// r.Use()

	r.GET("/api/main", _MainRoute)

	documents.InitApi(r)
	api.InitApi(r)
	apiApp.InitApiApp(r)
	auth.Init_api_auth0(r)

	llmApi.InitLLMApi(r)

	if tools.IsInDevelopment() {
		lg.LogI("Going to take files from: ", "../static")
		r.Use(static.Serve("/", static.LocalFile("../static", true)))

	} else {
		lg.LogI("Going to take files from: ", "./static")
		r.Use(static.Serve("/", static.LocalFile("./static", true)))
	}

	//r.NoRoute(tools.ProxyMiddleware2())

	prefix := "http://"
	port := "8080"

	if credentials.RunHttps() {
		prefix = "https://"
	}

	addr := "0.0.0.0:" + port

	server.RunMyioServer(r)

	lg.LogI("Running in: ", vfs.GetExeFile())

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

	lg.LogOk("------------------------------------")
	lg.LogOk("Mode: ", gin.Mode())

	lg.LogOk("Running server at: ", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+addr, prefix+addr))
	lg.LogOk("Access at: ", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+"127.0.0.1:"+port, prefix+"127.0.0.1:"+port))
	lg.LogOk("Access at: ", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+"localhost:"+port, prefix+"localhost:"+port))

	if vfs.IsLinux() {
		cmd := exec.Command("hostname", "-I")

		// Capture the output
		output, err := cmd.Output()
		if err != nil {
			lg.LogE(err)

		} else {
			ips := strings.Fields(string(output))

			for _, ip := range ips {
				lg.LogOk("Access at: ", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+ip+":"+port, prefix+ip+":"+port))

			}
		}

	}

	if credentials.RunHttps() {
		lg.LogOk("Started HTTP(S) branch")
		error := srv.ListenAndServeTLS("cert.pem", "key.pem")

		if error != nil {
			lg.LogE(error)
		}

	} else {
		lg.LogOk("Started HTTP branch")
		error := srv.ListenAndServe()

		if error != nil {
			lg.LogE(error)
		}
	}

	lg.LogE("Execution ended")
}

func main() {

	models.RegisterClazzFactory()

	//db.InitGorm()

	lg.LogI("Starting infinity twin application")
	lg.LogI("DbName: ", credentials.GetDBName())

	dev_main()
	//TestMain()
}

var (
	// These variables will be set during the build process
	Version   string
	BuildTime string
)
