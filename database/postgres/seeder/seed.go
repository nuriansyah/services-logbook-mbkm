package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nuriansyah/services-logbook-mbkm/cmd/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func main() {
	//Konfigurasi koneksi ke database
	configuration := config.New(".env")
	db, err := config.NewInitializedDatabase(configuration)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Slice yang berisi nama-nama mahasiswa
	mhsMap := map[string]string{
		"Muhammad Fathurrohman Haris":     "203040015",
		"Eldi Anugrah Pratama":            "203040031",
		"Muhammad Anendha Zaska":          "203040090",
		"Salsabila Febrianty Nurul Putri": "203040104",
		"Devin Febrian":                   "203040141",
		"Chandra Arcychan Azfar":          "203040151",
		"Agung Alfatah":                   "203040171",
		"Fikri Ridwan Faturohman":         "203040001",
		"Aril Fajri Tolani":               "203040013",
		"Mohammad Nabil Fadilah":          "203040020",
		"Hilma Sulaeman":                  "203040029",
		"Caprio Parlindungan Sinaga":      "203040036",
		"Alfi Syahrin Rachmat":            "203040054",
		"R.S. Ahmad Agil Alatas":          "203040093",
		"Hervin Fakhrul Mahardika":        "203040097",
		"Yudha Permana":                   "203040101",
		"Mochamad Yazman Yazid":           "203040164",
		"Harits Naufal Ichsan":            "203040077",
		"Abizar Valentino Fachri":         "203040086",
		"Dhaifullah Hilmy":                "203040119",
		"Kalam Mahardhika":                "203040123",
		"Dian Nurcahya Ningrum":           "203040124",
		"Fresabayu Anggoro":               "203040137",
		"Renandra Rahadian Putri":         "203040153",
		"Fajar Nur Siqid":                 "203040074",
		"Dimas Putra Kusumah":             "203040088",
		"Mochamad Rafi Ramadhan":          "203040051",
	}

	dosens := []string{"Fajar Darmawan, ST., M.Kom", "Erik S.T., M.Kom", "R.Sandhika Galih Amalga S.T., M.T", "Miftahul Fadli Muttaqin S.T., M.T"}

	// Seed data pada tabel dosen

	for _, dosen := range dosens {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		// Generate a unique email for each dosen
		email := strings.ToLower(strings.ReplaceAll(strings.Split(dosen, " ")[0], " ", "")) + "@gmail.com"

		// Insert data into the dosen table
		_, err = db.Exec("INSERT INTO dosen (name, email, password) VALUES ($1, $2, $3)", dosen, email, string(hashedPassword))
		if err != nil {
			log.Fatal(err)
		}
	}

	//Seed data pada tabel mahasiswa

	for nama, nrp := range mhsMap {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		// Insert data ke tabel mahasiswa
		_, err = db.Exec("INSERT INTO mahasiswa (name, nrp, password) VALUES ($1, $2, $3)", nama, nrp, string(hashedPassword))
		if err != nil {
			log.Fatal(err)
		}
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
	fmt.Println("Data berhasil ditambahkan.")
}
