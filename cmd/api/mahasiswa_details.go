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

type DetailsResponseReporting struct {
	ResponseMahasiswaDetails
	Details []ResponseReportingDetails `json:"details"`
}

type ResponseMahasiswaDetails struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Nrp       string `json:"nrp"`
	Company   string `json:"company"`
	ProgramKM string `json:"program_km"`
	LearnPath string `json:"learn_path"`
	Batch     int    `json:"batch"`
}

type ResponseReportingDetails struct {
	ID        int    `json:"ReportID"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
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

	c.JSON(http.StatusCreated, SuccessMahasiswaResponse{Message: "Success Insert Data!"})
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
	mhs, err := api.userRepo.FetchMahasiwaDetailsByDosenID(dosenID)
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

	details := make([]ResponseReportingDetails, 0)

	for _, detil := range mhs {

		details = append(details, ResponseReportingDetails{
			ID:      detil.ReportID,
			Title:   detil.Title,
			Content: detil.Content,
			Status:  detil.Status,
		})
	}

	postIDqueue := make([]int, 0)
	postsDetail := make(map[int]ResponseMahasiswaDetails)

	for _, post := range mhs {
		if _, ok := postsDetail[post.ID]; !ok {
			if len(postIDqueue) == 0 || postIDqueue[len(postIDqueue)-1] != post.ID {
				postIDqueue = append(postIDqueue, post.ID)
			}
			if post.Batch.Valid {
				batch = int(post.Batch.Int32)
			}
			if post.Company.Valid {
				company = post.Company.String
			}
			if post.Program.Valid {
				programKM = post.Program.String
			}
			if post.LearnPath.Valid {
				learnPath = post.LearnPath.String
			}
			postsDetail[post.ID] = ResponseMahasiswaDetails{
				ID:        post.ID,
				Nama:      post.Name,
				Nrp:       post.Nrp,
				Company:   company,
				ProgramKM: programKM,
				LearnPath: learnPath,
				Batch:     batch,
			}
		}
	}

	detailed := make(map[int][]ResponseReportingDetails)

	for _, post := range mhs {
		if _, ok := detailed[post.ID]; !ok {
			detailed[post.ID] = make([]ResponseReportingDetails, 0)
		}

		if post.ReportID != 0 {
			detailed[post.ID] = append(detailed[post.ID], ResponseReportingDetails{
				ID:        post.ReportID,
				Title:     post.Title,
				Content:   post.Content,
				Status:    post.Status,
				CreatedAt: post.CreatedAT.Format("2006-01-02 15:04:05"),
			})
		}
	}

	postsReponse := make([]DetailsResponseReporting, 0)

	for _, postID := range postIDqueue {
		postsReponse = append(postsReponse, DetailsResponseReporting{
			ResponseMahasiswaDetails: postsDetail[postID],
			Details:                  detailed[postID],
		})
	}

	c.JSON(http.StatusOK, postsReponse)
}
