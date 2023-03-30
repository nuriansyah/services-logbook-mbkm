package repository

import (
	"database/sql"
	"time"
)

type Mahasiswa struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Nrp       string `json:"nrp"`
	Password  string `json:"password"`
	Company   string `json:"company"`
	Program   string `json:"program"`
	Prodi     string `json:"prodi"`
	LearnPath string `json:"learnPath"`
	Batch     int    `json:"batch"`
	DosenID   int    `json:"dosen_id"`
	DosenName string `json:"dosenName"`
}

type Dosen struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	passwrod string `json:"passwrod"`
}

type Pembimbing struct {
	ID                int    `json:"id"`
	MahasiswaID       int    `json:"mahasiswa_id"`
	DosenPembimbingID int    `json:"dosen_pembimbing_id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	StatusID          int    `json:"status_id"`
}

type Reporting struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Type         string         `json:"type"`
	Status       string         `json:"status"`
	CreatedAT    time.Time      `json:"created_at"`
	PembimbingID int            `json:"pembimbing_id"`
	StatusID     int            `json:"statusID"`
	ImageID      sql.NullInt32  `json:"image_id"`
	ImagePath    sql.NullString `json:"image_path"`
}
type Comment struct {
	Id        int       `json:"id"`
	Comment   string    `json:"comments"`
	DosenName string    `json:"dosen_name"`
	MhsName   string    `json:"mhs_name"`
	CreatedAT time.Time `json:"created_at"`
	PostID    int       `json:"post_id"`
	DosenID   int       `json:"dosenID"`
	MhsID     int       `json:"mhsID"`
}
