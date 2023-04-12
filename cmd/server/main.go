package main

import (
	"github.com/nuriansyah/services-logbook-mbkm/cmd/api"
	"github.com/nuriansyah/services-logbook-mbkm/cmd/config"
	"github.com/nuriansyah/services-logbook-mbkm/internal/repository"
)

func main() {

	configuration := config.New(".env")
	db, err := config.NewInitializedDatabase(configuration)
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	pembRepo := repository.NewPembimbingRepository(db)
	reportRepo := repository.NewReportingRepository(db)
	detailMhsRepo := repository.NewDetailMahasiswaRepository(db)
	commentsReport := repository.NewCommnetsRepository(db)

	mainAPI := api.NewAPi(*userRepo, *pembRepo, *reportRepo, *detailMhsRepo, *commentsReport)
	mainAPI.Start()
}
