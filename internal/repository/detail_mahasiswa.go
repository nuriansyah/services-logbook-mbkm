package repository

import (
	"database/sql"
	"time"
)

type DetailMahasiswaRepository struct {
	db *sql.DB
}

func NewDetailMahasiswaRepository(db *sql.DB) *DetailMahasiswaRepository {
	return &DetailMahasiswaRepository{db: db}
}

type MahasiswaDetail struct {
	Id        int            `db:"id"`
	Name      string         `db:"name"`
	Nrp       string         `db:"nrp"`
	Password  string         `db:"password"`
	Company   sql.NullString `db:"company"`
	Program   sql.NullString `db:"program_km"`
	Prodi     sql.NullString `json:"prodi"`
	LearnPath sql.NullString `db:"learn_path"`
	Batch     sql.NullInt32  `db:"batch"`
	DosenID   int            `json:"dosen_id"`
	DosenName string         `json:"dosenName"`
}

type MahasiswaDetails struct {
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	Nrp       string         `db:"nrp"`
	Company   sql.NullString `db:"company"`
	Program   sql.NullString `db:"program_km"`
	Prodi     sql.NullString `json:"prodi"`
	LearnPath sql.NullString `db:"learn_path"`
	Batch     sql.NullInt32  `db:"batch"`
	ReportID  int            `db:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Status    string         `json:"status"`
	CreatedAT time.Time      `db:"createdAt"`
}

func (d *DetailMahasiswaRepository) InsertDetailMahasiswa(mahasiswaID int, company string, programKM string, learnPath string, batch int) error {
	query := `INSERT INTO mahasiswa_details (mahasiswa_id, company, program_km, learn_path, batch)
						VALUES ($1, $2, $3, $4, $5)`
	_, err := d.db.Exec(query, mahasiswaID, company, programKM, learnPath, batch)
	if err != nil {
		return err
	}
	return nil
}

func (d *DetailMahasiswaRepository) EditDetailMahasiswa(mahasiswaID int, company string, programKM string, learnPath string, batch int) error {
	query := `UPDATE mahasiswa_details
						SET company = $2, program_km = $3, learn_path = $4, batch = $5
						WHERE mahasiswa_id = $1`
	_, err := d.db.Exec(query, mahasiswaID, company, programKM, learnPath, batch)
	if err != nil {
		return err
	}
	return nil
}

func (d *DetailMahasiswaRepository) FetchMahasiswaByID(mahasiswaID int) ([]MahasiswaDetail, error) {
	sqlStatement := `SELECT m.id,m.nrp,m.name, md.company, md.learn_path,md.program_km,md.batch, d.name
					 FROM mahasiswa m
					 LEFT JOIN mahasiswa_details md on m.id = md.mahasiswa_id
					 RIGHT OUTER JOIN pembimbing p on m.id = p.mahasiswa_id
					 JOIN dosen d on d.id = p.dosen_pembimbing_id
					 WHERE m.id = $1`
	rows, err := d.db.Query(sqlStatement, mahasiswaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mahasiswa := []MahasiswaDetail{}
	for rows.Next() {
		var m MahasiswaDetail

		err = rows.Scan(&m.Id, &m.Nrp, &m.Name, &m.Company, &m.LearnPath, &m.Program, &m.Batch, &m.DosenName)
		if err != nil {
			return nil, err
		}

		mahasiswa = append(mahasiswa, m)
	}
	return mahasiswa, nil

}
