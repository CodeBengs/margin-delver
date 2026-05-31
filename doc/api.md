# API Documentation

Base path untuk internal API versi pertama:

```text
/internal/v1
```

## Auth

### POST /internal/v1/auth/login

Login user aktif dari table `users`.

Request body:

```json
{
  "username": "DELVERADMIN1",
  "password": "delverAdmin1"
}
```

Contoh cURL:

```bash
curl -X POST "http://localhost:3030/internal/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"DELVERADMIN1\",\"password\":\"delverAdmin1\"}"
```

Response sukses:

```json
{
  "status": "success",
  "message": "Login success",
  "result": {
    "token": "generated-token",
    "token_type": "Bearer",
    "expires_at": "2026-05-25T10:00:00+07:00",
    "user": {
      "id": 1,
      "username": "DELVERADMIN1",
      "name": "Delver Administrator"
    }
  }
}
```

Catatan:

- Password disimpan sebagai `password_hash` memakai bcrypt.
- Default user hanya dibuat ketika `DB_SEED_DEFAULT_USER=true` dan table `users` masih kosong.
- Nilai default user diambil dari `AUTH_DEFAULT_USERNAME`, `AUTH_DEFAULT_PASSWORD`, dan `AUTH_DEFAULT_NAME`.
