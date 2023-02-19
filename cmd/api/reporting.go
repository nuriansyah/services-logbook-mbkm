package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RequestReporting struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ResponseReporting struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	DosenID   int    `json:"dosen_id"`
	StatusID  int    `json:"status_id"`
	CreatedAt string `json:"created_at"`
	SuccessReportingResponse
}

type ResponseInsertReporting struct {
	ID int `json:"id"`
	SuccessReportingResponse
}

type SuccessReportingResponse struct {
	Message string `json:"message"`
}

type ErrorReportingResponse struct {
	Message string `json:"error"`
}

func (api *API) insertRequest(c *gin.Context) {
	var req RequestReporting

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: "Invalid Request Body"})
		return
	}
	mhsID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	DosenID, err := api.reportRepo.FetchPembimbingByID(mhsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	reportID, err := api.reportRepo.InsertReporting(req.Title, req.Content, int(DosenID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ResponseInsertReporting{
		ID:                       reportID,
		SuccessReportingResponse: SuccessReportingResponse{Message: "Berhasil!"},
	})
}

func (api *API) readReporting(c *gin.Context) {

	var (
		reportID int
		err      error
	)
	userID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorReportingResponse{Message: err.Error()})
		return
	}

	if reportID, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, ErrorReportingResponse{Message: "Invalid Post ID"})
		return
	}
	reports, err := api.reportRepo.FetchAuthorIDbyReportID(reportID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
		return
	}
	if len(reports) == 0 {
		c.JSON(http.StatusNotFound, ErrorReportingResponse{Message: "Post Not Found"})
		return
	}
	c.JSON(http.StatusOK, ResponseReporting{
		ID:        int64(reportID),
		DosenID:   reports[0].PembimbingID,
		StatusID:  reports[0].StatusID,
		Title:     reports[0].Title,
		Content:   reports[0].Content,
		Type:      reports[0].Type,
		Status:    reports[0].Status,
		CreatedAt: reports[0].CreatedAT.Format("2006-01-02 15:04:05"),
		SuccessReportingResponse: SuccessReportingResponse{
			Message: "Success",
		},
	})
}

func (api *API) readsReportingByDosen(c *gin.Context) {

	var (
		err error
	)
	dosenID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorReportingResponse{Message: err.Error()})
		return
	}

	reportLists, err := api.reportRepo.FetchReportByDosenID(dosenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
		return
	}

	var reportListsResponse []ResponseReporting
	for _, reports := range reportLists {
		reportListsResponse = append(reportListsResponse, ResponseReporting{
			ID:        int64(reports.ID),
			Title:     reports.Title,
			Content:   reports.Content,
			DosenID:   reports.PembimbingID,
			StatusID:  reports.StatusID,
			CreatedAt: reports.CreatedAT.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, reportListsResponse)

}

func (api *API) readsReportingByMhs(c *gin.Context) {
	mhsID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorReportingResponse{Message: err.Error()})
		return
	}
	reportLists, err := api.reportRepo.FetchAuthorByMhsID(mhsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
		return
	}
	var reportListsResponse []ResponseReporting
	for _, reports := range reportLists {
		reportListsResponse = append(reportListsResponse, ResponseReporting{
			ID:        int64(reports.ID),
			DosenID:   reports.PembimbingID,
			StatusID:  reports.StatusID,
			Title:     reports.Title,
			Content:   reports.Content,
			Type:      reports.Type,
			Status:    reports.Status,
			CreatedAt: reports.CreatedAT.Format("2 Jan 2006 15:04"),
			SuccessReportingResponse: SuccessReportingResponse{
				Message: "Success",
			},
		})
	}
	c.JSON(http.StatusOK, reportListsResponse)
}

func (api *API) editRequest(c *gin.Context) {
	var req RequestReporting
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: "Invalid report ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: "Invalid Request Body"})
		return
	}
	mhsID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	DosenID, err := api.reportRepo.FetchPembimbingByID(mhsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	err = api.reportRepo.UpdateReporting(req.Title, req.Content, int(DosenID), reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessReportingResponse{
		Message: "Update Report Success",
	})
}

func (api *API) getUserIDAvoidPanic(ctx *gin.Context) (authorID int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover from panic")
		}
	}()

	authorID, _ = api.getUserIdFromToken(ctx)
	return
}

func (api *API) uploadPostDocs(ctx *gin.Context) {
	var (
		postID int
		err    error
	)

	if postID, err = strconv.Atoi(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorReportingResponse{Message: "Invalid Post ID"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorReportingResponse{Message: err.Error()})
		return
	}

	folderPath := "media/post"
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
		return
	}

	docs := form.File["docs"]
	var wg sync.WaitGroup
	var mu sync.Mutex
	allowedExtensions := []string{".pdf", ".docx"}

	for _, file := range docs {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()

			defer func() {
				if v := recover(); v != nil {
					log.Println(v)
					ctx.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: "Internal Server Error"})
					return
				}
			}()

			if !isWordOrPDFFile(file.Filename) {
				ctx.JSON(http.StatusBadRequest, ErrorReportingResponse{Message: "Invalid file type. Only Word or PDF files are allowed."})
				return
			}

			uploadedFile, err := file.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
				return
			}

			defer uploadedFile.Close()

			extension := filepath.Ext(file.Filename)
			extensionValid := false
			for _, ext := range allowedExtensions {
				if ext == extension {
					extensionValid = true
					break
				}
			}

			if !extensionValid {
				ctx.JSON(http.StatusBadRequest, ErrorReportingResponse{Message: "Invalid file extension. Allowed extensions are: .pdf, .docx"})
				return
			}

			unixTime := time.Now().UTC().UnixNano()
			fileName := fmt.Sprintf("%d-%d-%s", postID, unixTime, file.Filename)
			fileLocation := filepath.Join(folderPath, fileName+extension)
			targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
				return
			}

			defer targetFile.Close()

			if _, err := io.Copy(targetFile, uploadedFile); err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
				return
			}
			mu.Lock()
			//if err := api.postRepo.InsertPostDoc(postID, fileLocation); err != nil {
			//	ctx.JSON(http.StatusInternalServerError, ErrorReportingResponse{Message: err.Error()})
			//	return
			//}
			mu.Unlock()
		}(file)
	}

	wg.Wait()

	ctx.JSON(http.StatusOK, SuccessReportingResponse{Message: "Post Documents Uploaded"})
}
func isWordOrPDFFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".doc" || ext == ".docx" || ext == ".pdf"
}
