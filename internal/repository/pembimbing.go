package repository

import (
	"database/sql"
	"errors"
	"net/http"
)

type PembimbingRepository struct {
	db *sql.DB
}

func NewPembimbingRepository(db *sql.DB) *PembimbingRepository {
	return &PembimbingRepository{db: db}
}

type PembimbingDetail struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func (p *PembimbingRepository) RequestPembimbing(mahasiswa_id, dosen_pembimbing_id int) (responseCode int, err error) {
	// Check if the mahasiswa exists in the table
	var countMahasiswa int
	err = p.db.QueryRow("SELECT COUNT(*) FROM mahasiswa WHERE id = $1", mahasiswa_id).Scan(&countMahasiswa)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// If the mahasiswa does not exist, return an error message
	if countMahasiswa == 0 {
		return http.StatusBadRequest, errors.New("Mahasiswa not found")
	}

	// Check if the dosen_pembimbing exists in the table
	var countDosenPembimbing int
	err = p.db.QueryRow("SELECT COUNT(*) FROM dosen WHERE id = $1", dosen_pembimbing_id).Scan(&countDosenPembimbing)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// If the dosen_pembimbing does not exist, return an error message
	if countDosenPembimbing == 0 {
		return http.StatusBadRequest, errors.New("Dosen Pembimbing not found")
	}

	// Check if the combination of mahasiswa_id and dosen_pembimbing_id already exists in the table
	var count int
	err = p.db.QueryRow("SELECT COUNT(*) FROM pembimbing WHERE mahasiswa_id = $1 AND dosen_pembimbing_id = $2", mahasiswa_id, dosen_pembimbing_id).Scan(&count)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// If the combination already exists, return an error message
	if count > 0 {
		return http.StatusConflict, errors.New("Pembimbing already exists")
	}
	// If the combination does not exist, insert the new row
	sqlStatement := "INSERT INTO pembimbing (mahasiswa_id, dosen_pembimbing_id, type, status_id) VALUES ($1,$2,'request',1) RETURNING id"
	var id int
	err = p.db.QueryRow(sqlStatement, mahasiswa_id, dosen_pembimbing_id).Scan(&id)

	return http.StatusCreated, err
}

func (p *PembimbingRepository) AcceptedPembimbing(mahasiswa_id int) error {
	sqlStatement := "UPDATE pembimbing SET type = 'bimbingan', status_id = 1 WHERE mahasiswa_id = $1;"

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(sqlStatement, mahasiswa_id)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PembimbingRepository) RejectRequestPembimbing(mhsId int) error {
	sqlStatement := "DELETE FROM pembimbing WHERE mahasiswa_id = $1"
	_, err := p.db.Exec(sqlStatement, mhsId)
	return err
}

func (p *PembimbingRepository) FetchAllRequestByID(dosenID int) ([]PembimbingDetail, error) {
	rows, err := p.db.Query("SELECT p.id,m.name,p.type,s.status  FROM pembimbing p left join mahasiswa m on m.id = p.mahasiswa_id LEFT JOIN status s on s.id = p.status_id WHERE type = 'request' AND dosen_pembimbing_id = $1", dosenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pembimbingList := []PembimbingDetail{}
	for rows.Next() {
		var pembimbing PembimbingDetail
		err := rows.Scan(&pembimbing.ID, &pembimbing.Name, &pembimbing.Type, &pembimbing.Status)
		if err != nil {
			return nil, err
		}
		pembimbingList = append(pembimbingList, pembimbing)
	}
	return pembimbingList, nil
}

func (p *PembimbingRepository) FetchAllBimbinganByID(dosenID int) ([]Pembimbing, error) {
	rows, err := p.db.Query("SELECT id, mahasiswa_id, dosen_pembimbing_id, type, status_id FROM pembimbing WHERE type = 'bimbingan' AND dosen_pembimbing_id = $1", dosenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pembimbingList := []Pembimbing{}
	for rows.Next() {
		var pembimbing Pembimbing
		err := rows.Scan(&pembimbing.ID, &pembimbing.MahasiswaID, &pembimbing.DosenPembimbingID, &pembimbing.Type, &pembimbing.StatusID)
		if err != nil {
			return nil, err
		}
		pembimbingList = append(pembimbingList, pembimbing)
	}
	return pembimbingList, nil
}

func (p *PembimbingRepository) FetchMhsIdByRequestId(reqId int) (int, error) {
	sqlStatement := "SELECT mahasiswa_id FROM pembimbing WHERE mahasiswa_id = $1"

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var mhsID int
	err = tx.QueryRow(sqlStatement, reqId).Scan(&mhsID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrPostNotFound
		}

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return mhsID, nil
}

func (p *PembimbingRepository) FetchAllBimbibinganReqByID(dosenID int) ([]Pembimbing, error) {
	rows, err := p.db.Query("SELECT id,mahasiswa_id,dosen_pembimbing_id,type, status_id FROM pembimbing WHERE type = 'bimbingan' AND status_id = 1 AND dosen_pembimbing_id = $1", dosenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pembimbingList := []Pembimbing{}
	for rows.Next() {
		var pembimbing Pembimbing
		err := rows.Scan(&pembimbing.ID, &pembimbing.MahasiswaID, &pembimbing.DosenPembimbingID, &pembimbing.Type, &pembimbing.StatusID)
		if err != nil {
			return nil, err
		}
		pembimbingList = append(pembimbingList, pembimbing)
	}
	return pembimbingList, err
}

func (p *PembimbingRepository) AcceptedBimbingan(mahasiswa_id int) (resposeCode int, err error) {
	sqlStatement := "UPDATE pembimbing SET type = 'bimbingan', status_id = 2 WHERE mahasiswa_id = $1;"
	var id int
	err = p.db.QueryRow(sqlStatement, mahasiswa_id).Scan(&id)

	return http.StatusAccepted, err
}

func (p *PembimbingRepository) RejectedBimbingan(mahasiswa_id int) (responseCode int, err error) {
	sqlStatement := "UPDATE pembimbing SET type = 'bimbingan', status_id = 3 WHERE mahasiswa_id = $1;"
	var id int
	err = p.db.QueryRow(sqlStatement, mahasiswa_id).Scan(&id)

	return http.StatusAccepted, err
}

func (p *PembimbingRepository) FetchMhsID(mhsID int) (int, error) {
	sqlStatement := `SELECT mahasiswa_id from pembimbing
					left join mahasiswa m on m.id = pembimbing.mahasiswa_id
					where m.id = $1`
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	err = tx.QueryRow(sqlStatement, mhsID).Scan(&mhsID)

	return mhsID, nil
}

func (p *PembimbingRepository) FetchDosenID() ([]Dosen, error) {
	rows, err := p.db.Query("SELECT id, name FROM dosen")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dosens []Dosen
	for rows.Next() {
		var dosen Dosen
		if err := rows.Scan(&dosen.Id, &dosen.Name); err != nil {
			return nil, err
		}
		dosens = append(dosens, dosen)
	}
	return dosens, nil
}
