## Instalasi Golang

Setelah berhasil clone, lakukan perintah :

1. cp .env.example .env

- APP_PORT=8080
- DB_URL=
- FIREBASE_CREDENTIALS_PATH=./serviceAccountKey.json
- OPENWEATHERMAP_API_KEY=
- WEATHER_CITY=

Isi data diatas, dengan db yang telah dibuat melalui laravel, openweather api key dan city untuk mengecek cuaca

Untuk file serviceAccoountKey.json, bisa didapatkan melalui akun firebase yang anda miliki

Step by step mendapatkan serviceAccountKey.json

- Masuk project pada firebase
- Masuk ke project setting
- Masuk ke tab service account
- Pilih admin sdk golang
- Tekan tombol generate private key

Setelah berhasil, akan terunduh file serviceAccountKey.json. Kemudian masukan file tersebut, ke folder root proyek

Setelah langkah diatas berhasil dijalankan, jalankan perintah "go run main.go"

