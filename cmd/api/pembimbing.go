package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nuriansyah/services-logbook-mbkm/internal/repository"
	"net/http"
	"strconv"
)

type RequestPembimbing struct {
	DosenID int `json:"dosen_id"`
}
type ResponseListDosen struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ResponsePembimbing struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	SuccessPembimbingResponse
}
type ResponsePembimbingRequest struct {
	ID int `json:"id"`
	SuccessPembimbingResponse
}
type UpdatedPembimbing struct {
	MahasiswaID int `json:"mahasiswa_id"`
}

type ResponseAccepted struct {
	ID int `json:"id"`
	SuccessPembimbingResponse
}
type SuccessPembimbingResponse struct {
	Message string `json:"message"`
}

type ErrorPembimbingResponse struct {
	Message string `json:"error"`
}

func (api *API) reqPembimbing(c *gin.Context) {
	var (
		req = RequestPembimbing{}
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: "Invalid Request Body"})
		return
	}

	userId, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorPembimbingResponse{Message: "Unauthorized"})
		return
	}
	mhsID, err := api.pembRepo.FetchMhsID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	pemId, err := api.pembRepo.RequestPembimbing(mhsID, req.DosenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ResponsePembimbingRequest{
		ID: pemId,
		SuccessPembimbingResponse: SuccessPembimbingResponse{
			Message: "Success Request Pembimbing",
		},
	})
}

func (api *API) updatePembimbing(c *gin.Context) {

	var req UpdatedPembimbing

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: err.Error()})
		return
	}

	if err := api.pembRepo.AcceptedPembimbing(req.MahasiswaID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, SuccessPembimbingResponse{
		Message: "Diterima Pembimbing",
	})
}

func (api *API) deletePembimbing(c *gin.Context) {
	mhsId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = api.pembRepo.RejectRequestPembimbing(mhsId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembimbing request deleted successfully"})
}

func (api *API) rejectBimbingan(c *gin.Context) {
	mhsId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	responseCode, err := api.pembRepo.RejectedBimbingan(mhsId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(responseCode, gin.H{"message": "Bimbingan rejected successfully"})
}

func (api *API) getAllRequest(c *gin.Context) {
	dosenId, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorPembimbingResponse{"Unauthorized"})
		return
	}
	mentoredList, err := api.pembRepo.FetchAllRequestByID(dosenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	var mentoredListResponse []ResponsePembimbing
	for _, mentorship := range mentoredList {
		mentoredListResponse = append(mentoredListResponse, ResponsePembimbing{
			ID:                        mentorship.ID,
			Name:                      mentorship.Name,
			Type:                      mentorship.Type,
			Status:                    mentorship.Status,
			SuccessPembimbingResponse: SuccessPembimbingResponse{Message: "Success"},
		})
	}
	c.JSON(http.StatusOK, mentoredListResponse)
}

func (api *API) getAllDosen(c *gin.Context) {
	dosenList, err := api.pembRepo.FetchDosenID()
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	var dosenLists []ResponseListDosen
	for _, dosen := range dosenList {
		dosenLists = append(dosenLists, ResponseListDosen{
			ID:   dosen.Id,
			Name: dosen.Name,
		})
	}
	c.JSON(http.StatusOK, dosenLists)
}

func (api *API) getAllRequestBimbingan(c *gin.Context) {
	dosenID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	mentoredList, err := api.pembRepo.FetchAllBimbibinganReqByID(dosenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: err.Error()})
		return
	}
	var mentoredListResponse []ResponsePembimbing
	for _, mentorship := range mentoredList {
		mentoredListResponse = append(mentoredListResponse, ResponsePembimbing{
			ID:   mentorship.ID,
			Name: mentorship.Name,
			Type: mentorship.Type,
		})
	}
	c.JSON(http.StatusOK, mentoredListResponse)
}

func (api *API) updatedBimbingan(c *gin.Context) {
	var req UpdatedPembimbing

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: "Invalid Request Body"})
		return
	}

	mahasiswaID, err := api.pembRepo.FetchMhsIdByRequestId(req.MahasiswaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorPembimbingResponse{Message: "Your ID cann't read"})
	}

	if mhsID, err := api.pembRepo.AcceptedBimbingan(req.MahasiswaID); err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, ErrorPembimbingResponse{Message: "Post Not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorPembimbingResponse{Message: "Internal Server Error"})
		return
	} else if mahasiswaID != mhsID {
		c.JSON(http.StatusForbidden, ErrorPembimbingResponse{Message: "Forbidden"})
		return
	}

	c.JSON(http.StatusAccepted, ResponseAccepted{
		SuccessPembimbingResponse: SuccessPembimbingResponse{
			Message: "Diterima Pembimbing",
		},
	})

}
