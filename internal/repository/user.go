package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) LoginDosen(email, password string) (*int, error) {
	sqlStatement := "SELECT id, password FROM dosen WHERE email = $1"
	res := u.db.QueryRow(sqlStatement, email)
	var hashedPassword string
	var id int
	err := res.Scan(&id, &hashedPassword)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	fmt.Print(&hashedPassword)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			fmt.Println("Wrong Password")
		}
		panic(err)
	}
	return &id, nil
}

func (u *UserRepository) LoginMahasiswa(nrp, password string) (*int, error) {
	sqlStatement := "SELECT id, password FROM mahasiswa WHERE nrp = $1"
	res := u.db.QueryRow(sqlStatement, nrp)
	var hashedPassword string
	var id int
	err := res.Scan(&id, &hashedPassword)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("wrong Password")
		}
		return nil, err
	}
	return &id, nil
}

func (u *UserRepository) CheckEmail(email string) (bool, error) {
	sqlStatement := "SELECT count(*) FROM dosen WHERE email = $1"
	res := u.db.QueryRow(sqlStatement, email)
	var count int
	err := res.Scan(&count)
	if count > 0 {
		return false, err
	}
	return true, err
}

func (u *UserRepository) CheckNrp(nrp string) (bool, error) {
	sqlStatement := "SELECT count(*) FROM mahasiswa WHERE nrp = $1"
	res := u.db.QueryRow(sqlStatement, nrp)
	var count int
	err := res.Scan(&count)
	if count > 0 {
		return false, err
	}
	return true, err
}

func (u *UserRepository) GetUserData(id int) (*Mahasiswa, error) {
	statement := `SELECT * FROM mahasiswa WHERE id = $1`
	var mhs Mahasiswa
	res := u.db.QueryRow(statement, id)
	err := res.Scan(&mhs.Id, &mhs.Name, &mhs.Nrp, &mhs.Prodi, &mhs.Company, &mhs.Program, &mhs.Batch)
	return &mhs, err
}

func (u *UserRepository) UpdateDetailDataUser(userID, batch int, nrp, prodi, program, company string) error {
	sqlStmt := `UPDATE user_details SET nrp = $1,prodi = $2,program = $3,company = $4,batch = $5 WHERE user_id = $6`
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(sqlStmt, nrp, prodi, program, company, batch, userID)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) InsertUserMahasiswa(name, nrp, password string) (userId, responCode int, err error) {

	isAvailable, err := u.CheckNrp(nrp)
	if err != nil {
		return -1, http.StatusBadRequest, err
	}

	if !isAvailable {
		return -1, http.StatusBadRequest, errors.New("email has been used")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	sqlStatement := `INSERT INTO mahasiswa (name,nrp,password) VALUES ($1,$2,$3) RETURNING id`

	var id int
	err = u.db.QueryRow(sqlStatement, name, nrp, hashedPassword).Scan(&id)

	return userId, http.StatusOK, err
}

func (u *UserRepository) InsertUserDosen(name, email, password string) (userId, responCode int, err error) {

	isAvailable, err := u.CheckEmail(email)
	if err != nil {
		return -1, http.StatusBadRequest, err
	}

	if !isAvailable {
		return -1, http.StatusBadRequest, errors.New("email has been used")
	}
	regex, err := regexp.Compile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}

	isValid := regex.Match([]byte(email))
	if !isValid {
		return -1, http.StatusBadRequest, errors.New("invalid email")
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	sqlStatement := `INSERT INTO dosen (name,email,password) VALUES ($1,$2,$3) RETURNING id`

	var id int
	err = u.db.QueryRow(sqlStatement, name, email, hashedPassword).Scan(&id)

	return userId, http.StatusOK, err
}

func (u *UserRepository) InsertDetailMahasiswa(mhs_id, batch int, company, program_km, learn_path string) (userId, responseCode int, err error) {
	sqlStatement := "INSERT INTO mahasiswa_details (mahasiswa_id, company, program_km, learn_path, batch) VALUES ($1, $2, $3, $4, $5)"
	var id int
	err = u.db.QueryRow(sqlStatement, mhs_id, batch, company, program_km, learn_path).Scan(&id)

	return userId, responseCode, err
}

func (u *UserRepository) ChangePasswordMahasiswa(id int, password string) error {
	sqlStatement := "UPDATE mahasiswa SET password = $1 WHERE id = $2"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = u.db.Exec(sqlStatement, hashedPassword, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) ChangePasswordDosen(id int, password string) error {
	sqlStatement := "UPDATE dosen SET password = $1 WHERE id = $2"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = u.db.Exec(sqlStatement, string(hashedPassword), id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) FetchMahasiswaByDosenID(dosenID int) ([]MahasiswaDetail, error) {
	sqlStatement := `SELECT m.id,d.name, m.name,m.nrp,md.company,md.program_km,md.learn_path,md.batch
					 FROM dosen d
					 LEFT JOIN pembimbing p on d.id = p.dosen_pembimbing_id
					 LEFT JOIN mahasiswa m on p.mahasiswa_id = m.id
					 LEFT JOIN mahasiswa_details md on m.id = md.mahasiswa_id
					 WHERE p.dosen_pembimbing_id = $1`
	rows, err := u.db.Query(sqlStatement, dosenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mahasiswa := []MahasiswaDetail{}
	for rows.Next() {
		var m MahasiswaDetail
		err = rows.Scan(&m.Id, &m.DosenName, &m.Name, &m.Nrp, &m.Company, &m.LearnPath, &m.Program, &m.Batch)
		if err != nil {
			return nil, err
		}
		mahasiswa = append(mahasiswa, m)
	}
	return mahasiswa, nil
}

func (u *UserRepository) FetchMahasiwaDetailsByDosenID(dosenID int) ([]MahasiswaDetails, error) {
	sqlStatement := `SELECT  m.id,m.name,m.nrp,
        md.company, md.program_km,md.learn_path,md.batch,
        r.id,r.title,r.content,
        s.status,
        r.created_at
					 FROM mahasiswa m
					 LEFT JOIN mahasiswa_details md on m.id = md.mahasiswa_id
					 RIGHT OUTER JOIN pembimbing p on m.id = p.mahasiswa_id
				     JOIN status s on s.id = p.status_id
					 JOIN dosen d on d.id = p.dosen_pembimbing_id
                     JOIN reporting r on p.id = r.pembimbing_id
                     	WHERE p.dosen_pembimbing_id = $1`
	rows, err := u.db.Query(sqlStatement, dosenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mhsDetails := []MahasiswaDetails{}
	for rows.Next() {
		var mds MahasiswaDetails
		err = rows.Scan(&mds.ID, &mds.Name, &mds.Nrp,
			&mds.Company, &mds.Program, &mds.LearnPath, &mds.Batch,
			&mds.ReportID, &mds.Title, &mds.Content, &mds.Status, &mds.CreatedAT)
		if err != nil {
			return nil, err
		}
		mhsDetails = append(mhsDetails, mds)
	}
	return mhsDetails, err
}

func (u *UserRepository) FetchDataDosen(id int) ([]Dosen, error) {
	sqlStatement := `SELECT id,name,email FROM dosen WHERE id =$1`
	rows, err := u.db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dosen := []Dosen{}
	for rows.Next() {
		var d Dosen
		err = rows.Scan(&d.Id, &d.Name, &d.Email)
		if err != nil {
			return nil, err
		}
		dosen = append(dosen, d)
	}
	return dosen, nil
}

func (u *UserRepository) CheckPasswordDosen(userId int, password string) (bool, error) {
	var hashedPassword string
	err := u.db.QueryRow("SELECT password FROM dosen WHERE id = $1", userId).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
func (u *UserRepository) CheckPasswordMahasiswa(userId int, password string) (bool, error) {
	var hashedPassword string
	err := u.db.QueryRow("SELECT password FROM mahasiswa WHERE id = $1", userId).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil // Password does not match
		}
		return false, err
	}

	return true, nil // Password matches
}
