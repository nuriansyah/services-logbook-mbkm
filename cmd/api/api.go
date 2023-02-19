package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nuriansyah/services-logbook-mbkm/cmd/config"
	"github.com/nuriansyah/services-logbook-mbkm/internal/repository"
)

type API struct {
	userRepo        repository.UserRepository
	pembRepo        repository.PembimbingRepository
	reportRepo      repository.ReportingRepository
	detailmhsReport repository.DetailMahasiswaRepository
	router          *gin.Engine
}

func NewAPi(userRepo repository.UserRepository, pembRepo repository.PembimbingRepository, reportRepo repository.ReportingRepository, detailmhsReport repository.DetailMahasiswaRepository) API {
	router := gin.Default()
	api := &API{
		userRepo:        userRepo,
		pembRepo:        pembRepo,
		reportRepo:      reportRepo,
		detailmhsReport: detailmhsReport,
		router:          router,
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	router.POST("/loginDosen", api.loginDosen)
	router.POST("/loginMahasiswa", api.loginMahasiswa)
	router.POST("/registerMahasiswa", api.registerMahasiswa)
	router.POST("/registerDosen", api.registerDosen)

	router.GET("/dosen", api.getAllDosen)

	pembimbingRouter := router.Group("/api/pembimbing", AuthMiddleware())
	{
		pembimbingRouter.POST("/reqPembimbing", api.reqPembimbing)
		pembimbingRouter.PUT("/accepted", api.updatedBimbingan)
		pembimbingRouter.PUT("/", api.updatePembimbing)
		pembimbingRouter.GET("/getRequest", api.getAllRequest)
		pembimbingRouter.GET("/getBimbinganRequest", api.getAllRequestBimbingan)
		pembimbingRouter.DELETE("/:id", api.deletePembimbing)
		pembimbingRouter.PUT("/:id/reject", api.rejectBimbingan)
	}
	reportingRouter := router.Group("/api/reports", AuthMiddleware())
	{
		reportingRouter.POST("/", api.insertRequest)
		reportingRouter.GET("/post/:id", api.readReporting)
		reportingRouter.PUT("/:id", api.editRequest)
		reportingRouter.GET("/postsDosen", api.readsReportingByDosen)
		reportingRouter.GET("/postsMhs", api.readsReportingByMhs)

	}
	dosenRouter := router.Group("/api/dosen", AuthMiddleware())
	{
		dosenRouter.PUT("/changePassword", api.changePassword)
		dosenRouter.GET("/fetch", api.fetchMahasiswaByDsn)
		dosenRouter.GET("/fetchData", api.fetchDataDosen)
	}
	detailmhsRouter := router.Group("/api/mahasiswa", AuthMiddleware())
	{
		detailmhsRouter.POST("/detail", api.insertDetailMhs)
		detailmhsRouter.PUT("/", api.editDetailMhs)
		detailmhsRouter.GET("/fetch", api.fetchMahasiswaByMhs)
	}
	router.Use(gin.Recovery())

	return *api

}

func (api *API) Handler() *gin.Engine {
	return api.router
}

func (api *API) Start() {
	setPort := config.New(".env")
	api.Handler().Run(setPort.Get("APP_PORT"))
}
