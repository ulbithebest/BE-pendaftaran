# BE-pendaftaran (Refactored)

Backend pendaftaran HIMATIF dengan struktur modular/clean ala pdfmbackend.

## Struktur Folder Baru

```
cmd/main.go
config/
  - db.go
  - token.go
controller/
  - auth.go
  - registration.go
  - admin.go
model/
  - user.go
  - registration.go
route/
  - route.go
helper/
  - jwt.go
  - response.go
middleware/
  - auth.go
uploads/
mod/ (untuk fitur/layanan lain di masa depan)
```

## Menjalankan

1. `go mod tidy`
2. Jalankan MongoDB lokal
3. Jalankan:
   ```sh
   go run ./cmd/main.go
   ```

## Endpoint Utama
- /api/register, /api/login
- /api/me, /api/registration, /api/upload/cv
- /api/admin/registrations, /api/admin/registration/:id, /api/admin/registration/:id/status

## Catatan
- Semua handler, model, helper, config, dan middleware sudah dipisah sesuai best practice.
- Siap untuk pengembangan modular di folder `mod/`.

---

Struktur dan kode sudah mengikuti inspirasi dari project pdfmbackend untuk maintainability dan scalability.


Backend web service untuk aplikasi pendaftaran organisasi HIMATIF berbasis Golang (Fiber v2) dan MongoDB.

## Struktur Folder

```
cmd/main.go
controllers/
  - auth_controller.go
  - registration_controller.go
  - admin_controller.go
models/
  - user.go
  - registration.go
routes/routes.go
utils/
  - db.go
  - jwt.go
  - response.go
middleware/auth.go
uploads/ (untuk file CV)
```

## Setup & Menjalankan

1. **Clone repo dan install dependencies**
   ```sh
   go mod tidy
   ```
2. **Set environment variable (opsional):**
   - `MONGODB_URI` (default: mongodb://localhost:27017)
   - `JWT_SECRET` (default: supersecretkey)
3. **Jalankan server**
   ```sh
   go run ./cmd/main.go
   ```

## Endpoint API

### AUTH
- `POST /api/register` — Register user
- `POST /api/login` — Login user, return JWT token

### USER (wajib JWT)
- `GET /api/me` — Info user
- `POST /api/registration` — Submit form pendaftaran
- `GET /api/registration` — Data pendaftaran user
- `POST /api/upload/cv` — Upload CV (PDF)

### ADMIN (wajib JWT + admin)
- `GET /api/admin/registrations` — Semua pendaftar
- `GET /api/admin/registration/:id` — Detail pendaftar
- `POST /api/admin/registration/:id/status` — Update status (`lulus`, `tidak_lulus`, `menunggu`)

## Format Collection MongoDB

### users
- `_id`, `name`, `nim`, `email`, `password` (hashed), `role` ("admin"/"user")

### registrations
- `_id`, `user_id`, `division`, `motivation`, `cv_path`, `status`, `note`, `updated_at`

## Catatan
- Semua response dalam format JSON rapi.
- Upload CV ke folder `uploads/`, path file disimpan di database.
- Untuk akses admin, ubah field `role` user di DB menjadi `admin`.

---

Jika ada error atau butuh bantuan lebih lanjut, silakan hubungi pengembang atau cek file source code untuk dokumentasi lebih lanjut.