package repository

import (
	"database/sql"
	"errors"
	"time"
)

type ReportingRepository struct {
	db *sql.DB
}

func NewReportingRepository(db *sql.DB) *ReportingRepository {
	return &ReportingRepository{db: db}
}

var (
	ErrPostNotFound = errors.New("Post not found")
)

func (r *ReportingRepository) InsertReporting(title, content string, dosenID int) (userId int, err error) {

	sqlStatement := "INSERT INTO reporting (title,content,pembimbing_id,created_at,status_id) VALUES ($1,$2,$3,$4,1) RETURNING id"
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var id int
	err = tx.QueryRow(sqlStatement, title, content, dosenID, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, err
}

func (r *ReportingRepository) UpdateReporting(title, content string, dosenID, reportID int) (err error) {
	sqlStatement := "UPDATE reporting SET title = $1, content = $2, pembimbing_id = $3, updated_at = $4 WHERE id = $5"
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sqlStatement, title, content, dosenID, time.Now(), reportID)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ReportingRepository) FetchAuthorIDbyReportID(postID, AuthorID int) ([]Reporting, error) {

	var (
		reports []Reporting
	)

	sqlStatement := `SELECT r.id, title, content, created_at, p.type,s.status,r.pembimbing_id,r.status_id
    				 FROM reporting r
					 LEFT JOIN pembimbing p on p.id = r.pembimbing_id
					 LEFT JOIN status s on s.id = r.status_id
					 WHERE r.id = $1 AND p.mahasiswa_id = $2`
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := r.db.Query(sqlStatement, postID, AuthorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report Reporting
		err := rows.Scan(
			&report.ID,
			&report.Title,
			&report.Content,
			&report.CreatedAT,
			&report.Type,
			&report.Status,
			&report.PembimbingID,
			&report.StatusID,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportingRepository) FetchReportByDosenID(pembID int) ([]Reporting, error) {
	sqlStatement := "SELECT id,title,content,pembimbing_id,status_id,created_at FROM reporting WHERE pembimbing_id = $1"
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(sqlStatement, pembID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Reporting
	for rows.Next() {
		var report Reporting
		err := rows.Scan(
			&report.ID, &report.Title, &report.Content, &report.PembimbingID, &report.StatusID, &report.CreatedAT,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return reports, err
}

func (r *ReportingRepository) FetchAuthorByMhsID(AuthorID int) ([]Reporting, error) {

	var (
		reports []Reporting
	)

	sqlStatement := `SELECT r.id, title, content, created_at, p.type,s.status,r.pembimbing_id,r.status_id
    				 FROM reporting r
					 LEFT JOIN pembimbing p on p.id = r.pembimbing_id
					 LEFT JOIN status s on s.id = r.status_id
					 WHERE  p.mahasiswa_id = $1`
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := r.db.Query(sqlStatement, AuthorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report Reporting
		err := rows.Scan(
			&report.ID,
			&report.Title,
			&report.Content,
			&report.CreatedAT,
			&report.Type,
			&report.Status,
			&report.PembimbingID,
			&report.StatusID,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportingRepository) FetchPembimbingByID(mhsID int) (int, error) {
	sqlStatement := "SELECT id FROM pembimbing WHERE mahasiswa_id = $1"
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var id int
	err = tx.QueryRow(sqlStatement, mhsID).Scan(&id)

	return id, err
}

func (r *ReportingRepository) InsertFileReporting(path string, dosenID int) (err error) {
	sqlStatement := `INSERT INTO reporting_files (file_path,created_at,pembimbing_id) VALUES $1,$2,$3`
	tx, err := r.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	var id int
	err = tx.QueryRow(sqlStatement, path, time.Now(), dosenID).Scan(&id)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ReportingRepository) CountReportingApproved(mhsID int) (int, error) {
	sqlStatement := `SELECT count(r.status_id)
						FROM reporting r 
						LEFT JOIN pembimbing p ON r.pembimbing_id = p.id
						WHERE r.status_id = 2 AND p.mahasiswa_id = $1`
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow(sqlStatement, mhsID).Scan(&count)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportingRepository) CountReportingPending(mhsID int) (int, error) {
	sqlStatement := `SELECT count(r.status_id)
						FROM reporting r 
						LEFT JOIN pembimbing p ON r.pembimbing_id = p.id
						WHERE r.status_id = 1 AND p.mahasiswa_id = $1`
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow(sqlStatement, mhsID).Scan(&count)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportingRepository) CountReportingReject(mhsID int) (int, error) {
	sqlStatement := `SELECT count(r.status_id ) 
							FROM reporting r 
							LEFT JOIN pembimbing p on r.pembimbing_id = p.id
							WHERE r.status_id = 3 AND r.pembimbing_id = $1`
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(sqlStatement, mhsID).Scan(&mhsID)
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return mhsID, err
}

func (r *ReportingRepository) FetchMhsID(mhsID int) (int, error) {
	sqlStatement := `SELECT id from pembimbing
					left join mahasiswa m on m.id = pembimbing.mahasiswa_id
					where m.id = $1`
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	err = tx.QueryRow(sqlStatement, mhsID).Scan(&mhsID)

	return mhsID, nil
}
