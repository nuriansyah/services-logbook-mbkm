package main

import (
	"fmt"
	"github.com/nuriansyah/services-logbook-mbkm/cmd/config"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Konfigurasi koneksi ke database
	configuration := config.New(".env")
	db, err := config.NewInitializedDatabase(configuration)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Slice yang berisi nama-nama mahasiswa
	mahasiswa := []string{"Ahmad Rizki", "Budi Santoso", "Cici Kurniawati", "Dedi Prasetyo", "Eka Fitriani", "Fajar Wibowo", "Gita Dewi", "Hendra Nugraha", "Irma Sari", "Joko Widodo", "Kiki Anggraini", "Laila Nurul", "Maman Syahid", "Nurul Hidayah", "Oktavia Sari"}

	// Seed data pada tabel dosen

	for i := 1; i <= 8; i++ {
		// Hash password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO dosen (name, email, password) VALUES ($1, $2, $3)", "Dosen "+strconv.Itoa(i), "dosen"+strconv.Itoa(i)+"@gmail.com", string(hashedPassword))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Seed data pada tabel mahasiswa
	nrp := "1830400"
	for i, nama := range mahasiswa {
		// Generate NRP
		var nrpValue string
		if i+1 < 10 {
			nrpValue = nrp + "0" + strconv.Itoa(i+1)
		} else {
			nrpValue = nrp + strconv.Itoa(i+1)
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		// Insert data ke tabel mahasiswa
		_, err = db.Exec("INSERT INTO mahasiswa (name, nrp, password) VALUES ($1, $2, $3)", nama, nrpValue, string(hashedPassword))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data berhasil ditambahkan.")
}
