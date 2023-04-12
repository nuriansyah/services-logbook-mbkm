package main

import (
	"fmt"
	"github.com/nuriansyah/services-logbook-mbkm/cmd/config"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Mahasiswa struct {
	ID       uint
	Name     string
	NRP      string
	Password string
	Details  MahasiswaDetails
}

type MahasiswaDetails struct {
	ID          uint
	MahasiswaID uint
	Company     string
	ProgramKM   string
	LearnPath   string
	Batch       int
}
type Dosen struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Pembimbing struct {
	ID                int
	MahasiswaID       int
	DosenPembimbingID int
	Type              string
	StatusID          int
}

func main() {
	//Konfigurasi koneksi ke database
	configuration := config.New(".env")
	db, err := config.NewInitializedDatabase(configuration)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Slice yang berisi nama nama dosens
	dosenList := []Dosen{
		{
			Name:     "Fajar Darmawan, ST., M.Kom",
			Email:    "fajar@unpas.ac.id",
			Password: "password",
		},
		{
			Name:     "Erik, ST., M.Kom",
			Email:    "erik@unpas.ac.id",
			Password: "password",
		},
		{
			Name:     "R. Sandhika Galih Amalga, ST., MT",
			Email:    "r.sandhika@unpas.ac.id",
			Password: "password",
		},
		{
			Name:     "Miftahul Fadli Muttaqin, ST., MT",
			Email:    "miftahul@unpas.ac.id",
			Password: "password",
		},
	}
	// Slice yang berisi nama-nama mahasiswa
	mahasiswaList := []Mahasiswa{
		{
			Name:     "Muhammad Fathurrohman Haris",
			NRP:      "203040015",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Machine Learning",
				Batch:     4,
			},
		},
		{
			Name:     "Eldi Anugrah Pratama",
			NRP:      "203040031",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Machine Learning",
				Batch:     4,
			},
		},
		{
			Name:     "Mochamad Rafi Ramadhan",
			NRP:      "203040051",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)",
				Batch:     4,
			},
		},

		{
			Name:     "Muhammad Anendha Zaska",
			NRP:      "203040090",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)",
				Batch:     4,
			},
		},
		{
			Name:     "Salsabila Febrianty Nurul Putri",
			NRP:      "203040104",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Machine Learning",
				Batch:     4,
			},
		},
		{
			Name:     "Devin Febrian",
			NRP:      "203040141",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "BISA AI Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Blokchain",
				Batch:     4,
			},
		},
		{
			Name:     "Chandra Arcychan Azfar",
			NRP:      "203040151",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)\n",
				Batch:     4,
			},
		},
		{
			Name:     "Agung Alfatah",
			NRP:      "203040171",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Machine Learning",
				Batch:     4,
			},
		},
		{
			Name:     "Fikri Ridwan Faturohman",
			NRP:      "203040001",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Android Mobile Development",
				Batch:     4,
			},
		},
		{
			Name:     "Aril Fajri Tolani",
			NRP:      "203040013",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Alterra Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Engineer Career (ReactJS)\n",
				Batch:     4,
			},
		},
		{
			Name:     "Mohammad Nabil Fadilah",
			NRP:      "203040020",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Alterra Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Engineer Career (ReactJS)",
				Batch:     4,
			},
		},

		{
			Name:     "Hilman Sulaeman",
			NRP:      "203040029",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "BISA AI Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "UI/UX Design",
				Batch:     4,
			},
		},
		{
			Name:     "Caprio Parlindungan Sinaga",
			NRP:      "203040036",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "BISA AI Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "UI/UX Design",
				Batch:     4,
			},
		},
		{
			Name:     "Alfi Syahrin Rachmat",
			NRP:      "203040054",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Data Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Data Science",
				Batch:     4,
			},
		},
		{
			Name:     "R.S. Ahmad Agil Alatas",
			NRP:      "203040093",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Alterra Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Engineer Career (ReactJS)",
				Batch:     4,
			},
		},
		{
			Name:     "Hervin Fakhrul Mahardika",
			NRP:      "203040097",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)",
				Batch:     4,
			},
		},
		{
			Name:     "Yudha Permana",
			NRP:      "203040101",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Alterra Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Engineer Career (ReactJS)",
				Batch:     4,
			},
		},
		{
			Name:     "Mochamad Yazman Yazid",
			NRP:      "203040164",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Zenius",
				ProgramKM: "Studi Independen",
				LearnPath: "UI/UX Design",
				Batch:     4,
			},
		},
		{
			Name:     "Harits Naufal Ichsan",
			NRP:      "203040077",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Dicoding Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Web & Back-End Dev",
				Batch:     4,
			},
		},
		{
			Name:     "Abizar Valentino Fachri",
			NRP:      "23040086",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)",
				Batch:     4,
			},
		},
		{
			Name:     "Dhaifullah Hilmy",
			NRP:      "203040119",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Dicoding Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Web & Back-End Dev",
				Batch:     4,
			},
		},
		{
			Name:     "Kalam Mahardhika",
			NRP:      "203040123",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)",
				Batch:     4,
			},
		},
		{
			Name:     "Dian Nurcahya Ningrum",
			NRP:      "203040124",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Dicoding Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Front-End Web & Back-End Dev",
				Batch:     4,
			},
		},
		{
			Name:     "Fresabayu Anggoro",
			NRP:      "203040137",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Machine Learning",
				Batch:     4,
			},
		},
		{
			Name:     "Renandra Rahadian Putri",
			NRP:      "203040153",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Bangkit Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Google Cloud Platform)",
				Batch:     4,
			},
		},
		{
			Name:     "Fajar Nur Sidiq",
			NRP:      "203040074",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "Metrodata Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Network Security Center Operation",
				Batch:     4,
			},
		},
		{
			Name:     "Dimas Putra Kusumah",
			NRP:      "203040088",
			Password: "password",
			Details: MahasiswaDetails{
				Company:   "BISA AI Academy",
				ProgramKM: "Studi Independen",
				LearnPath: "Cloud Computing (Cyber Security)",
				Batch:     4,
			},
		},
	}
	//Slice data pembimbing
	bimbinganList := []Pembimbing{
		{
			MahasiswaID:       1,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       2,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       3,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       4,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       5,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       6,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       7,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       8,
			DosenPembimbingID: 1,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       9,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       10,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       11,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       12,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       13,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       14,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       15,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       16,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       17,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       18,
			DosenPembimbingID: 2,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       19,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       20,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       21,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       22,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       23,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       24,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       25,
			DosenPembimbingID: 3,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       26,
			DosenPembimbingID: 4,
			Type:              "Bimbingan",
			StatusID:          1,
		},
		{
			MahasiswaID:       27,
			DosenPembimbingID: 4,
			Type:              "Bimbingan",
			StatusID:          1,
		},
	}
	for _, dosen := range dosenList {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		// Insert data into the dosen table
		_, err = db.Exec("INSERT INTO dosen (name, email, password) VALUES ($1, $2, $3)", dosen.Name, dosen.Email, string(hashedPassword))
		if err != nil {
			log.Fatal(err)
		}
	}

	//Seed data pada tabel mahasiswa

	for _, mahasiswa := range mahasiswaList {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		// Insert data into mahasiswa table
		stmt, err := db.Prepare(`
		INSERT INTO mahasiswa (name, nrp, password)
		VALUES ($1, $2, $3)
		RETURNING id;
		`)
		if err != nil {
			log.Fatalf("Failed to prepare statement: %v", err)
		}
		defer stmt.Close()

		var mahasiswaID uint
		err = stmt.QueryRow(mahasiswa.Name, mahasiswa.NRP, hashedPassword).Scan(&mahasiswaID)
		if err != nil {
			log.Fatalf("Failed to insert data into mahasiswa table: %v", err)
		}

		// Insert data into mahasiswa_details table
		stmt, err = db.Prepare(`
		INSERT INTO mahasiswa_details (mahasiswa_id, company, program_km, learn_path, batch)
		VALUES ($1, $2, $3, $4, $5);
		`)
		if err != nil {
			log.Fatalf("Failed to prepare statement: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(mahasiswaID, mahasiswa.Details.Company, mahasiswa.Details.ProgramKM, mahasiswa.Details.LearnPath, mahasiswa.Details.Batch)
		if err != nil {
			log.Fatalf("Failed to insert data into mahasiswa_details table: %v", err)
		}
		fmt.Println("Data seeded successfully!")
	}

	_, err = db.Exec("INSERT INTO status (status) VALUES ('Pending')")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO status (status) VALUES ('Accepted')")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO status (status) VALUES ('Rejected')")
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range bimbinganList {
		_, err := db.Exec("INSERT INTO pembimbing (mahasiswa_id, dosen_pembimbing_id, type, status_id) VALUES ($1, $2, $3, $4)",
			p.MahasiswaID, p.DosenPembimbingID, p.Type, p.StatusID)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Data berhasil ditambahkan.")
}
