# User Entity (Mini Jira)

## ERD (User)

```
users
-----
PK  id           bigint
    name         varchar(100)
    username     varchar(100) UNIQUE
    password     varchar(255)
    email        varchar(100) UNIQUE
    role         varchar(50) DEFAULT 'member'
    created_at   timestamp
    updated_at   timestamp
    deleted_at   timestamp (nullable)
```

## Role User

| Role | Deskripsi | Contoh Hak Akses |
| --- | --- | --- |
| `admin` | Pengelola sistem utama | Mengelola user, konfigurasi global, audit. |
| `project_manager` | Penanggung jawab proyek | Membuat proyek, mengatur sprint, assign task. |
| `member` | Anggota tim | Mengerjakan task, update status, komentar. |

## Catatan Implementasi

- Field `role` disimpan sebagai string untuk kemudahan query, tetapi direpresentasikan sebagai tipe `Role` di aplikasi agar nilai lebih terkontrol.
- Default role adalah `member`.
