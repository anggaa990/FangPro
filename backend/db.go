package main

import (
    "database/sql"
    "log"
    "os"

    _ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
    dbPath := "tobacco.db"

    // Cek apakah file DB sudah ada
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        log.Println("Database belum ada. Membuat baru...")

        // Buat file kosong
        file, err := os.Create(dbPath)
        if err != nil {
            log.Fatal("Gagal membuat file database:", err)
        }
        file.Close()
    }

    // Koneksi ke SQLite dengan parameter anti-lock
    // PENTING: tambahkan query parameters untuk mengatasi database locking
    database, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
    if err != nil {
        log.Fatal("Gagal membuka database:", err)
    }

    // Set connection pool - KRUSIAL untuk SQLite!
    database.SetMaxOpenConns(1)  // SQLite hanya support 1 writer
    database.SetMaxIdleConns(1)  // Keep 1 connection alive
    
    if err := database.Ping(); err != nil {
        log.Fatal("Tidak bisa terhubung ke database:", err)
    }

    log.Println("Database terhubung:", dbPath)

    // Jalankan schema.sql
    schema, err := os.ReadFile("../sql/schema.sql")
    if err != nil {
        log.Fatal("Gagal membaca schema.sql:", err)
    }

    if _, err := database.Exec(string(schema)); err != nil {
        log.Fatal("Gagal menjalankan schema:", err)
    }

    log.Println("Schema database OK")
    DB = database
}