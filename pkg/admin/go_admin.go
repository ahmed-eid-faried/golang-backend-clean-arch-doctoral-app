package Admin

// import (
// 	"net/http"
//

// 	_ "github.com/GoAdminGroup/go-admin/adapter/gin" // Import the adapter, it must be imported. If it is not imported, you need to define it yourself.
// 	"github.com/GoAdminGroup/go-admin/engine"
// 	config "github.com/GoAdminGroup/go-admin/modules/config"
// 	// _ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // Import the sql driver
// 	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
// 	"github.com/GoAdminGroup/go-admin/modules/language"
// 	"github.com/GoAdminGroup/go-admin/plugins/admin"
// 	_ "github.com/GoAdminGroup/themes/adminlte" // Import the theme
// 	"github.com/gin-gonic/gin"
// 	_ "github.com/lib/pq"
// )

// type AdminPanel struct {
// 	eng      *engine.Engine
// 	cfgAdmin config.Config
// }

// // router interface{} //ex ::===>>> // r *gin.Engine

// func (s *AdminPanel) Run(router *gin.Engine) {
// 	//*********************************************************

// 	// Instantiate a GoAdmin engine object.
// 	s.eng = engine.Default()

// 	// GoAdmin global configuration, can also be imported as a json file.
// 	s.cfgAdmin = config.Config{
// 		Databases: config.DatabaseList{
// 			"default": {
// 				Host:         "127.0.0.1",
// 				Port:         "5432",
// 				User:         "postgres",
// 				Pwd:          "postgres",
// 				Name:         "postgres",
// 				Driver:       "postgresql",
// 				MaxIdleConns: 50,
// 				MaxOpenConns: 150,
// 				// ConnMaxLifetime: 0,
// 				// ConnMaxIdleTime: 0,
// 				// Params:          nil,
// 			},
// 		},
// 		UrlPrefix: "admin", // The url prefix of the website.
// 		// Store must be set and guaranteed to have write access, otherwise new administrator users cannot be added.
// 		Store: config.Store{
// 			Path:   "./uploads",
// 			Prefix: "uploads",
// 		},
// 		Language: language.EN,
// 	}

// 	// Add configuration and plugins, use the Use method to mount to the web framework.
// 	_ = s.eng.AddConfig(&s.cfgAdmin).
// 		Use(router)

// 		//*********************************************************
// 	// Define a simple home route
// 	router.GET("/", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Welcome to the E-commerce Admin Panel")
// 	})

// 	// Add custom tables and charts
// 	addTables(s.eng)
// }

// func addTables(eng *engine.Engine) {
// 	adminPlugin := admin.NewAdmin(nil)

// 	// Add the Profile table
// 	adminPlugin.AddGenerator("products", GetProfileTable)
// 	// Add the products table
// 	adminPlugin.AddGenerator("posts", GetPostsTable)
// 	// Add the users table
// 	adminPlugin.AddGenerator("posts", GetUserTable)

// 	eng.AddPlugins(adminPlugin)
// }
