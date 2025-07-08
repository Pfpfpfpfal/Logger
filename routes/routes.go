package routes

import (
	"fmt"
	"log"
	"logview_ui/config"
	"logview_ui/handlers"
	hLogfiles "logview_ui/handlers/logfiles"
	hPages "logview_ui/handlers/pages"
	hTabs "logview_ui/handlers/tabs"

	"github.com/gin-gonic/gin"

	"database/sql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
)

const baseUrl = "/logview"

func SetupRouter() *gin.Engine {

	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	dbConfig := config.LoadDBConfig()

	println(dbConfig.Host)
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", "postgres_user", "postgres_password", "postgres_container", "5432", "postgres_db")
	println(dbUrl)
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Не удалось установить соединение с базой данных")
	}
	store, err := postgres.NewStore(db, []byte("secret"))
	if err != nil {
		log.Fatalf("Не удалось создать хранилище сессий: %v", err)
	}

	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static/")

	{
		// Корневая группа роутов по base url
		base := r.Group(baseUrl)
		{
			// Группа роутов pages
			pages := base.Group("/pages")
			pages.GET("/logtable", hPages.GetLogViewPage)
		}
		{
			// Группа роутов data
			pages := base.Group("/data")
			pages.GET("/table_structure", handlers.GetTableStructure)

			pages.POST("/tabs", hTabs.SaveTab)
			pages.PATCH("/tabs/:tab_id", hTabs.ModifyTab)
			pages.DELETE("/tabs/:tab_id", hTabs.DeleteTab)
			pages.GET("/tabs/:tab_id", hTabs.GetTab)
			pages.GET("/tabs", hTabs.GetTabs)

			pages.GET("/tabs/:tab_id/history", handlers.GetFiltersFromHistory)
			pages.POST("/tabs/:tab_id/history", handlers.SaveFilterToHistory)
			pages.POST("/tabs/:tab_id/history/index", handlers.SaveHistoryIndex)

			pages.GET("/logfiles", hLogfiles.GetLogFiles)
			pages.GET("/logfiles/:file_id", hLogfiles.GetLogFile)
			pages.GET("/logfiles/:file_id/logs", hLogfiles.GetLogFileLogs)
			pages.GET("/logfiles/:file_id/logs/top/:field_name", hLogfiles.GetFieldTop10)
		}
	}

	return r
}
