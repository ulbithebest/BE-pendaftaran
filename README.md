# BE-pendaftaran

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