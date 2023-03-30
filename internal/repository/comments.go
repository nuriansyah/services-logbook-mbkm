package repository

import (
	"database/sql"
	"time"
)

type CommentsRepository struct {
	db *sql.DB
}

func NewCommnetsRepository(db *sql.DB) *CommentsRepository {
	return &CommentsRepository{db: db}
}

func (c *CommentsRepository) SelectAllMahasiswaCommentsByPostID(PostID int) ([]Comment, error) {
	sqlStmt := `SELECT
    mc.id,r.id ,m.name,mc.comment,mc.created_at FROM mahasiswa_comment mc
    LEFT JOIN mahasiswa m on m.id = mc.mahasiswa_id
    LEFT JOIN reporting r on mc.reporting_id = r.id
    WHERE mc.reporting_id = $1
    ORDER BY created_at`

	rows, err := c.db.Query(sqlStmt, PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commnets []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.Id,
			&comment.PostID,
			&comment.MhsName,
			&comment.Comment,
			&comment.CreatedAT,
		)
		if err != nil {
			return nil, err
		}
		commnets = append(commnets, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return commnets, err
}

func (c *CommentsRepository) SelectAllDosenCommentsByPostID(PostID int) ([]Comment, error) {
	sqlStmt := `SELECT
    dc.id,r.id,d.name,dc.comment,dc.created_at FROM dosen_comment dc
    LEFT JOIN reporting r on dc.reporting_id = r.id
    LEFT JOIN dosen d on dc.dosen_id = d.id
    WHERE dc.reporting_id = $1
    ORDER BY created_at`
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := c.db.Query(sqlStmt, PostID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound
		}

		return nil, err
	}
	defer rows.Close()

	var commnets []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.Id,
			&comment.PostID,
			&comment.DosenName,
			&comment.Comment,
			&comment.CreatedAT,
		)
		if err != nil {
			return nil, err
		}
		commnets = append(commnets, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return commnets, err
}

//func (c *CommentsRepository) FetchCommentsAuthorsID(id int) (int, error) {
//
//}

func (c *CommentsRepository) InsertCommentDosen(comment Comment) (int, error) {
	sqlStmt := `INSERT INTO dosen_comment (comment, created_at, reporting_id, dosen_id)
                VALUES ($1, $2, $3, $4) 
                RETURNING id`
	var id int
	err := c.db.QueryRow(sqlStmt, comment.Comment, time.Now(), comment.PostID, comment.DosenID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CommentsRepository) InsertCommentMahasiswa(comment Comment) (int64, error) {
	sqlStmt := `INSERT INTO mahasiswa_comment (comment, created_at, reporting_id, mahasiswa_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var commentID int64
	err := c.db.QueryRow(sqlStmt, comment.Comment, time.Now(), comment.PostID, comment.MhsID).Scan(&commentID)
	if err != nil {
		return 0, err
	}
	return commentID, nil
}
