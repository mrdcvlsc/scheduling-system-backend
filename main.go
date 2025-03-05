package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Routes/RoutesV1"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

var SessionStore = cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

func main() {
	fmt.Println("Starting backend service")

	//////////////////////////////////////////////////////////////////////////
	//                    Initialize Persistence To Use
	//////////////////////////////////////////////////////////////////////////

	RouteGlobals.ResourcesPersistence = &StorageResources.Persistence{
		ReaderService: &StorageResources.JsonReader{},
		WriterService: &StorageResources.JsonWriter{},
	}

	RouteGlobals.SchedulePersistence = &StorageSchedule.Persistence{
		LoadService: &StorageSchedule.JsonReader{},
		SaveService: &StorageSchedule.JsonWriter{},
	}

	RouteGlobals.InitializeCachedUniversitySchedule()

	//////////////////////////////////////////////////////////////////////////
	// MongoDB Setup
	//////////////////////////////////////////////////////////////////////////

	// fmt.Println("Connecting to MongoDB...")

	// fmt.Printf("MONGO_DB_USER     = %s\n", os.Getenv("MONGO_DB_USER"))
	// fmt.Printf("MONGO_DB_PASSWORD = %s\n", os.Getenv("MONGO_DB_PASSWORD"))
	// fmt.Printf("PORT              = %s\n", os.Getenv("PORT"))

	// // Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// opts := options.Client().ApplyURI(fmt.Sprintf(
	// 	"mongodb+srv://%s:%s@testcluster.sz6qg.mongodb.net/?retryWrites=true&w=majority&appName=TestCluster",
	// 	os.Getenv("MONGO_DB_USER"),
	// 	os.Getenv("MONGO_DB_PASSWORD"),
	// )).SetServerAPIOptions(serverAPI)

	// // Create a new client and connect to the server
	// client, err := mongo.Connect(context.TODO(), opts)
	// if err != nil {
	// 	panic(err)
	// }

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// // Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	//////////////////////////////////////////////////////////////////////////

	use_secure_cookie := false
	same_site := http.SameSiteDefaultMode

	if os.Getenv("GIN_MODE") == "release" {
		use_secure_cookie = true
		same_site = http.SameSiteNoneMode
	}

	SessionStore.Options(sessions.Options{
		// MaxAge:   259200, // 3 days
		// MaxAge:   60, // 1 minute
		MaxAge:   60 * 15, // 15 minute
		Secure:   use_secure_cookie,
		HttpOnly: true,
		SameSite: same_site,
	})

	//////////////////////////////////////////////////////////////////////////

	router := gin.Default()

	// maximum memory limit for multipart form file uploads
	router.MaxMultipartMemory = 5 << 20 // 5 MiB

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))
	router.Use(sessions.Sessions("session_id", SessionStore))

	if gin.Mode() != gin.ReleaseMode {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://192.168.1.*:5173", "http://192.168.0.*:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowWildcard:    true, // Enable wildcard support for 192.168.1.*
		}))
	}

	//////////////////////////////////////////////////////////////////////////
	//                              API-v1
	//////////////////////////////////////////////////////////////////////////

	v1 := router.Group("/v1")

	v1.GET("/const", RoutesV1.GetConst)

	v1.GET("/all_departments", RoutesV1.GetAllDepartments)
	v1.GET("/department_data", RoutesV1.GetDepartmentData)

	v1.GET("/instructors/d", RoutesV1.GetDepartmentInstructorsDefaults)
	v1.GET("/instructors/a", RoutesV1.GetDepartmentInstructorsAllocated)
	v1.POST("/instructor_add", RoutesV1.PostInstructor)
	v1.PATCH("/instructor_update", RoutesV1.PatchInstructor)
	v1.DELETE("/instructor_remove", RoutesV1.DeleteInstructor)

	v1.DELETE("/room_remove", RoutesV1.DeleteRoom)

	v1.GET("/university_schedule", RoutesV1.GetUniversitySchedule)
	v1.POST("/university_schedule", RoutesV1.PostUniversitySchedule)

	v1.GET("/class_schedule", RoutesV1.GetClassSchedule)
	v1.GET("/class_json_schedule", RoutesV1.GetJsonClassSchedule)

	v1.POST("/generate_schedule", RoutesV1.GenerateSchedule)

	if os.Getenv("GIN_MODE") != "release" {
		v1.GET("/generate_schedule", RoutesV1.GenerateSchedule) // for dev only
	}

	//////////////////////////////////////////////////////////////////////////

	if os.Getenv("GIN_MODE") != "release" {
		router.GET("/cpu", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"cpu": runtime.NumCPU()})
		})

		router.GET("/die", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "bye-bye")
			os.Exit(0)
		})
	}

	Utils.DisplayOutboundIP(os.Getenv("PORT"))
	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

}
