package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nuriansyah/services-logbook-mbkm/internal/repository"
	"os"
)

type API struct {
	userRepo        repository.UserRepository
	pembRepo        repository.PembimbingRepository
	reportRepo      repository.ReportingRepository
	detailmhsReport repository.DetailMahasiswaRepository
	commentsRepo    repository.CommentsRepository
	router          *gin.Engine
}

func NewAPi(
	userRepo repository.UserRepository,
	pembRepo repository.PembimbingRepository,
	reportRepo repository.ReportingRepository,
	detailmhsReport repository.DetailMahasiswaRepository,
	commentsRepo repository.CommentsRepository,
) API {
	router := gin.Default()
	api := &API{
		userRepo:        userRepo,
		pembRepo:        pembRepo,
		reportRepo:      reportRepo,
		detailmhsReport: detailmhsReport,
		commentsRepo:    commentsRepo,
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
		pembimbingRouter.PUT("/:id/accepted", api.approvedBimbingan)
		pembimbingRouter.PUT("/:id/rejected", api.rejectedBimbingan)
	}
	reportingRouter := router.Group("/api/reports", AuthMiddleware())
	{
		reportingRouter.POST("/", api.insertRequest)
		reportingRouter.GET("/post/:id", api.readReporting)
		reportingRouter.PUT("/:id", api.editRequest)
		reportingRouter.GET("/postsDosen", api.readsReportingByDosen)
		reportingRouter.GET("/postsMhs", api.readsReportingByMhs)
		reportingRouter.POST("/upload", api.uploadPostDocs)
		reportingRouter.POST("/upload/img/:id", api.uploadPostDocs)
		reportingRouter.POST("/download/:post_id", api.downloadPostDoc)
		reportingRouter.GET("/approvedList", api.getApprovedReports)
		reportingRouter.GET("/pendingList", api.getPendingReports)
		reportingRouter.GET("/rejectList", api.getRejectedReports)

	}
	dosenRouter := router.Group("/api/dosen", AuthMiddleware())
	{
		dosenRouter.PUT("/changePasswordDosen", api.changePasswordDosen)
		dosenRouter.GET("/fetch", api.fetchMahasiswaByDsn)
		dosenRouter.GET("/fetchData", api.fetchDataDosen)
	}
	detailmhsRouter := router.Group("/api/mahasiswa", AuthMiddleware())
	{
		detailmhsRouter.PUT("/changePasswordMhs", api.changePasswordMahasiswa)
		detailmhsRouter.POST("/detail", api.insertDetailMhs)
		detailmhsRouter.PUT("/", api.editDetailMhs)
		detailmhsRouter.GET("/fetch", api.fetchMahasiswaByMhs)

	}
	commentsRouter := router.Group("/api/comments", AuthMiddleware())
	{
		commentsRouter.POST("/dosen", api.CreateCommentDosen)
		commentsRouter.POST("/mhs", api.CreateCommentMahasiswa)

		//commentsRouter.GET("/dosen/:postID", api.ReadCommentDosenByPostID)
	}
	router.GET("api/comments/mhs/:post_id", api.ReadCommentMahasiswaByPostID)
	router.GET("/api/comments/dosen/:post_id", api.ReadCommentDosenByPostID)
	router.GET("/api/comments/:post_id", api.ReadAllCommentsByPostID)

	router.Use(gin.Recovery())

	return *api

}

func (api *API) Handler() *gin.Engine {
	return api.router
}

func (api *API) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	api.Handler().Run(port)
}
