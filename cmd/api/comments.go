package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nuriansyah/services-logbook-mbkm/internal/repository"
	"github.com/nuriansyah/services-logbook-mbkm/utils"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type CreateCommentRequest struct {
	PostID  int    `json:"post_id" binding:"required,number"`
	Comment string `json:"comment" binding:"required"`
}
type CommentDosenResponse struct {
	Id        int    `json:"id"`
	Comment   string `json:"comments"`
	DosenName string `json:"dosen_name"`
	PostID    int    `json:"post_id"`
	CreatedAT string `json:"created_at"`
}
type CommentMahasiswaResponse struct {
	ID        int    `json:"id"`
	MhsName   string `json:"mhs_name"`
	Comment   string `json:"comments"`
	PostID    int    `json:"post_id"`
	CreatedAT string `json:"created_at"`
}

type CommentResponse struct {
	ID        int    `json:"id"`
	MhsName   string `json:"mhs_name,omitempty"`
	DosenName string `json:"dosen_name,omitempty"`
	Comment   string `json:"comments"`
	PostID    int    `json:"post_id"`
	CreatedAT string `json:"created_at"`
	Type      string `json:"types"`
}

func (api *API) ReadAllCommentsByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	mahasiswaComments, err := api.commentsRepo.SelectAllMahasiswaCommentsByPostID(postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dosenComments, err := api.commentsRepo.SelectAllDosenCommentsByPostID(postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var commentLists []CommentResponse

	for _, cld := range mahasiswaComments {
		id := uuid.New().ID()
		commentLists = append(commentLists, CommentResponse{
			ID:        int(id),
			MhsName:   cld.MhsName,
			Comment:   cld.Comment,
			PostID:    postID,
			CreatedAT: cld.CreatedAT.Format("2 Jan 2006 15:04 PM"),
			Type:      "mahasiswa",
		})
	}

	for _, cld := range dosenComments {
		id := uuid.New().ID()
		commentLists = append(commentLists, CommentResponse{
			ID:        int(id),
			DosenName: cld.DosenName,
			Comment:   cld.Comment,
			PostID:    postID,
			CreatedAT: cld.CreatedAT.Format("2 Jan 2006 15:04 PM"),
			Type:      "dosen",
		})
		// Sort the comments by CreatedAT in descending order
		sort.Slice(commentLists, func(i, j int) bool {
			timei, _ := time.Parse("2 Jan 2006 15:04 PM", commentLists[i].CreatedAT)
			timej, _ := time.Parse("2 Jan 2006 15:04 PM", commentLists[j].CreatedAT)
			return timei.Before(timej)
		})
	}

	c.JSON(http.StatusOK, commentLists)
}

func (api *API) ReadCommentMahasiswaByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if c.GetHeader("Authorization") != "" {
		_, err := api.getUserIdFromToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	comments, err := api.commentsRepo.SelectAllMahasiswaCommentsByPostID(postID)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	var commentListsMhs []CommentMahasiswaResponse
	for _, cld := range comments {
		commentListsMhs = append(commentListsMhs, CommentMahasiswaResponse{
			ID:        cld.Id,
			MhsName:   cld.MhsName,
			Comment:   cld.Comment,
			PostID:    postID,
			CreatedAT: cld.CreatedAT.Format("2 Jan 2006 15:04 PM"),
		})
	}
	c.JSON(http.StatusOK, commentListsMhs)
}

func (api *API) ReadCommentDosenByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	comments, err := api.commentsRepo.SelectAllDosenCommentsByPostID(postID)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	var commentListsDosen []CommentDosenResponse
	for _, cld := range comments {
		commentListsDosen = append(commentListsDosen, CommentDosenResponse{
			Id:        cld.Id,
			DosenName: cld.DosenName,
			Comment:   cld.Comment,
			PostID:    postID,
			CreatedAT: cld.CreatedAT.Format("2 Jan 2006 15:04 PM"),
		})
	}

	c.JSON(http.StatusOK, commentListsDosen)
}

func (api *API) CreateCommentDosen(c *gin.Context) {
	var createCommentRequest CreateCommentRequest
	err := c.ShouldBindJSON(&createCommentRequest)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"errors": utils.GetErrorMessage(ve)},
			)
		} else {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		}
		return
	}
	userId, err := api.getUserIdFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = api.commentsRepo.InsertCommentDosen(repository.Comment{
		PostID:  createCommentRequest.PostID,
		Comment: createCommentRequest.Comment,
		DosenID: userId,
	})
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"message": "Add Comment Successful"},
	)
}

func (api *API) CreateCommentMahasiswa(c *gin.Context) {
	var createCommentRequest CreateCommentRequest
	err := c.ShouldBindJSON(&createCommentRequest)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"errors": utils.GetErrorMessage(ve)},
			)
		} else {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		}
		return
	}
	userId, err := api.getUserIdFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = api.commentsRepo.InsertCommentMahasiswa(repository.Comment{
		PostID:  createCommentRequest.PostID,
		Comment: createCommentRequest.Comment,
		MhsID:   userId,
	})
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"message": "Add Comment Successful"},
	)
}
