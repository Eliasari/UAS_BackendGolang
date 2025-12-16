# Sistem Pelaporan Prestasi Mahasiswa – Backend API  
**UAS Pemrograman Backend Lanjut**

**Nama:** Elia Sari  
**NIM:** 434231023  
**Kelas:** C8  

Backend API ini dikembangkan sebagai bagian dari **Ujian Akhir Semester (UAS)** mata kuliah *Pemrograman Backend Lanjut*. Sistem ini mengimplementasikan **Sistem Pelaporan Prestasi Mahasiswa** sesuai dengan *Software Requirement Specification (SRS)* yang telah ditetapkan.

Aplikasi ini berfungsi sebagai **RESTful API** untuk mendukung proses pelaporan, verifikasi, dan pengelolaan prestasi mahasiswa dengan mekanisme autentikasi dan otorisasi berbasis peran.

---

## Deskripsi Sistem
Sistem Pelaporan Prestasi Mahasiswa merupakan aplikasi backend yang dirancang untuk mengelola data prestasi mahasiswa secara terstruktur dan fleksibel. Sistem ini mendukung beberapa peran pengguna, yaitu **Admin**, **Mahasiswa**, dan **Dosen Wali**, dengan hak akses yang berbeda sesuai dengan role dan permission masing-masing.

Data disimpan menggunakan pendekatan hybrid:
- **PostgreSQL** digunakan untuk data relasional seperti user, role, permission, dan relasi prestasi
- **MongoDB** digunakan untuk menyimpan detail prestasi yang bersifat dinamis sesuai jenis prestasi

Pendekatan ini diterapkan untuk memenuhi kebutuhan fleksibilitas data sebagaimana dijelaskan dalam dokumen SRS.

---

## Fitur Utama

### 1. Autentikasi dan Otorisasi (JWT & RBAC)
Sistem menerapkan autentikasi menggunakan **JSON Web Token (JWT)** dan otorisasi berbasis **Role-Based Access Control (RBAC)**.  
Setiap endpoint dilindungi oleh permission tertentu sesuai dengan peran pengguna.

Fitur ini mencakup:
- Login pengguna
- Validasi token
- Pembatasan akses endpoint berdasarkan role dan permission

---

### 2. Manajemen User dan Role
Admin memiliki hak akses untuk mengelola data pengguna dalam sistem, termasuk:
- Menambahkan dan mengubah data user
- Menetapkan role ke user
- Mengelola relasi mahasiswa dan dosen wali

Fitur ini mendukung kebutuhan pengelolaan sistem sesuai dengan SRS.

---

### 3. Manajemen Prestasi Mahasiswa
Mahasiswa dapat membuat dan mengelola data prestasi akademik maupun non-akademik.  
Sistem mendukung beberapa jenis prestasi, antara lain:
- Competition
- Academic
- Organization
- Publication
- Certification
- Other

Fitur yang tersedia:
- Membuat prestasi dalam status draft
- Mengubah data prestasi draft
- Menghapus prestasi draft
- Mengajukan prestasi untuk diverifikasi

---

### 4. Verifikasi Prestasi oleh Dosen Wali
Dosen wali dapat melakukan proses verifikasi terhadap prestasi mahasiswa bimbingannya, meliputi:
- Melihat daftar prestasi mahasiswa
- Menyetujui atau menolak prestasi
- Memberikan catatan pada prestasi yang ditolak

Alur verifikasi mengikuti workflow prestasi yang ditentukan dalam SRS.

---

### 5. List Prestasi dengan Filter, Sorting, dan Pagination
Sistem menyediakan endpoint untuk menampilkan daftar prestasi dengan fitur:
- Filter berdasarkan status prestasi
- Sorting data
- Pagination untuk efisiensi pengambilan data

Hak akses list prestasi disesuaikan dengan role pengguna:
- Admin dapat melihat seluruh data prestasi
- Dosen wali hanya melihat prestasi mahasiswa bimbingannya
- Mahasiswa hanya melihat prestasi miliknya sendiri

---

### 6. Dokumentasi API (Swagger)
Seluruh endpoint REST API didokumentasikan menggunakan **Swagger (swaggo)** untuk memudahkan pengujian API dan integrasi dengan frontend.

---

### 7. Unit Testing
Sistem dilengkapi dengan unit testing untuk menguji fungsi dan service utama guna memastikan logika bisnis berjalan sesuai dengan requirement.

---

## Tech Stack
- **Go (Golang)**
- **Fiber Framework**
- **PostgreSQL**
- **MongoDB**
- **Swagger (swaggo)**

---

## Cara Menjalankan Aplikasi

1. Clone repository
```bash
git clone https://github.com/Eliasari/UAS_BackendGolang.git
cd uas-prestasi
```
3. Konfigurasi environment
```bash
cp .env.example .env
```
5. Jalankan aplikasi
```bash
go run main.go
```
7. Akses dokumentasi Swagger
```bash
[go run main.go](http://localhost:3000/swagger/index.html
```
