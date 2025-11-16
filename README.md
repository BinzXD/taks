# Golang + PostgreSQL + Docker

Project ini adalah aplikasi **backend** menggunakan **Golang**, **GORM**, dan **PostgreSQL**, containerized dengan **Docker** dan **Docker Compose**.

📄 **Dokumentasi API:**  
👉 [Lihat di Postman](https://documenter.getpostman.com/view/49571281/2sB3WwowQz)

## Clone Repository
git clone https://github.com/BinzXD/taks.git
cd task


## Build & Run with Docker
Jalankan:
docker-compose up --build
Lalu:
docker ps
Kemudian buka dibrowser anda:
http://localhost:8080/


## Default Account
Gunakan akun default untuk login:
Email:    superadmin@gmail.com
Password: secret123


## Notes
Pastikan port 8080 & 5432 belum digunakan di lokal.
Database akan otomatis dibuat dan di-seed saat pertama kali container dijalankan.
Untuk development, gunakan GO_ENV=development.