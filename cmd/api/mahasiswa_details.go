package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var req struct {
	MahasiswaID int    `json:"mahasiswa_id"`
	Company     string `json:"company"`
	ProgramKM   string `json:"program_km"`
	LearnPath   string `json:"learn_path"`
	Batch       int    `json:"batch"`
}

type ResponseMahasiswa struct {
	ID        int    `json:"id"`
	DosenName string `json:"dosen_name"`
	Nama      string `json:"nama"`
	Nrp       string `json:"nrp"`
	Company   string `json:"company"`
	ProgramKM string `json:"program_km"`
	LearnPath string `json:"learn_path"`
	Batch     int    `json:"batch"`
}

type ErrorMahasiswaResponse struct {
	Message string `json:"message"`
}
type SuccessMahasiswaResponse struct {
	Message string `json:"message"`
}

func (api *API) insertDetailMhs(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mhsID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	// Call the repository to insert the detail mahasiswa
	if err := api.detailmhsReport.InsertDetailMahasiswa(mhsID, req.Company, req.ProgramKM, req.LearnPath, req.Batch); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (api *API) editDetailMhs(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the repository to edit the detail mahasiswa
	if err := api.detailmhsReport.EditDetailMahasiswa(req.MahasiswaID, req.Company, req.ProgramKM, req.LearnPath, req.Batch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessMahasiswaResponse{
		Message: "Success",
	})
}

func (api *API) fetchMahasiswaByMhs(c *gin.Context) {
	mhsID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	mhs, err := api.detailmhsReport.FetchMahasiswaByID(mhsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	if len(mhs) == 0 {
		c.JSON(http.StatusNotFound, ErrorMahasiswaResponse{Message: "Student not found"})
		return
	}
	var (
		company, programKM, learnPath string
		batch                         int
	)

	if mhs[0].Company.Valid {
		company = mhs[0].Company.String
	}
	if mhs[0].Program.Valid {
		programKM = mhs[0].Program.String
	}
	if mhs[0].LearnPath.Valid {
		learnPath = mhs[0].LearnPath.String
	}
	if mhs[0].Batch.Valid {
		batch = int(mhs[0].Batch.Int32)
	}
	c.JSON(http.StatusOK, ResponseMahasiswa{
		ID:        mhsID,
		DosenName: mhs[0].DosenName,
		Nama:      mhs[0].Name,
		Nrp:       mhs[0].Nrp,
		Company:   company,
		ProgramKM: programKM,
		LearnPath: learnPath,
		Batch:     batch,
	})
}
func (api *API) fetchMahasiswaByDsn(c *gin.Context) {
	dosenID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	mhs, err := api.userRepo.FetchMahasiswaByDosenID(dosenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	if len(mhs) == 0 {
		c.JSON(http.StatusNotFound, ErrorMahasiswaResponse{Message: "Student not found"})
		return
	}

	var reponseListMahasiswa []ResponseMahasiswa
	for _, mhsList := range mhs {
		var (
			company, programKM, learnPath string
			batch                         int
		)
		if mhsList.Company.Valid {
			company = mhsList.Company.String
		}
		if mhsList.Program.Valid {
			programKM = mhsList.Program.String
		}
		if mhsList.LearnPath.Valid {
			learnPath = mhsList.LearnPath.String
		}
		if mhsList.Batch.Valid {
			batch = int(mhsList.Batch.Int32)
		}
		reponseListMahasiswa = append(reponseListMahasiswa, ResponseMahasiswa{
			ID:        mhsList.Id,
			DosenName: mhsList.DosenName,
			Nama:      mhsList.Name,
			Nrp:       mhsList.Nrp,
			Company:   company,
			ProgramKM: programKM,
			LearnPath: learnPath,
			Batch:     batch,
		})
	}
	c.JSON(http.StatusOK, reponseListMahasiswa)
}
