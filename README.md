# MVC Go Template

Template ini adalah versi boilerplate dari struktur `mini_jira` agar bisa dipakai ulang untuk project baru.

## Struktur

- `handler/` controller layer
- `service/` business logic
- `repository/` data access
- `models/` entity database
- `dto/` request/response contract
- `routes/` route registration
- `internal/` bootstrap server + database
- `migrations/` migrasi schema

## Cara Pakai

1. Clone/copy folder template ini menjadi nama project baru.
2. Masuk ke folder project baru.
3. Rename module Go dan semua import internal:

```bash
./init_template.sh github.com/username/nama-project
```

4. Buat file environment:

```bash
cp .env.example .env
```

5. Sesuaikan isi `.env` (database, JWT key, SMTP, dll).
6. Jalankan:

```bash
go mod tidy
go run main.go
```

## Catatan

- Script `init_template.sh` wajib dijalankan sekali setelah clone agar import path tidak lagi memakai `mini_jira`.
- Folder ini tidak menyertakan file sensitif seperti `.env` dan artefak build.
