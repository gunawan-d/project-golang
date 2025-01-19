# Project GO


## Description Task


## Deskripsi:
### 1. Buatlah sebuah layanan HTTP yang menerima permintaan POST /refresh di port 8080. Layanan ini harus segera mengembalikan respons HTTP dengan status "created".

```bash
 - Jika proses refresh sedang berjalan, permintaan baru harus diantrikan.
 - Jika tidak ada proses refresh yang berjalan, maka harus langsung memulai proses refresh.
 - Detail Proses Refresh:

 - Proses refresh memakan waktu sekitar 10 detik untuk selesai (bisa gunakan sleep).
 - Hindari mengantrikan permintaan dalam detik yang sama (berdasarkan timestamp Unix).
 - Setiap refresh harus dijalankan secara berurutan dan mencatat timestamp kapan permintaan dilakukan.
 - Tidak perlu menyimpan antrean setelah restart aplikasi.
 - Proses refresh juga harus dijalankan sekali saat aplikasi pertama kali dijalankan.
```


### 2. Service Manager Implementation


Task :
1. Print nilai int (incremental) Setiap 10 seconds, ketika server jalan dia HIT nya ke 
2. Start Untuk menjelankan
3. Stop untuk 10X hit
4. Reload ketika di reload (nge clear / Reset incremental menjadi 0)


### Deskripsi
Tugas ini adalah membuat server sederhana menggunakan **Golang** yang akan mencetak nilai **incremental** setiap **10 detik** saat server berjalan. Server ini memiliki beberapa fitur utama:

1. **Start** → Memulai proses pencetakan nilai incremental setiap 10 detik.
2. **Stop** → Menghentikan proses setelah mencapai **10 kali hit**.
3. **Reload** → Mereset nilai incremental menjadi **0** saat direload.

### Cara Kerja
1. Saat server dijalankan, nilai **hit** akan dimulai dari **0**.
2. Setiap **10 detik**, server akan mencetak **Hit ke-X**.
3. Jika sudah mencapai **10 hit**, server akan otomatis berhenti.
4. Jika endpoint **reload** dipanggil, nilai **hit** akan kembali menjadi **0**.

### API Endpoint
| Method | Endpoint  | Deskripsi |
|--------|----------|-----------|
| `POST`  | `/start` | Memulai proses |
| `POST`  | `/stop`  | Menghentikan proses |
| `POST`  | `/reload` | Mereset nilai hit menjadi 0 |

### Prasyarat
- Pastikan **Golang** sudah terinstall di komputer Anda.
- Gunakan **Go modules** untuk mengatur dependency (opsional).

### Cara Menjalankan
1. Clone repositori ini atau buat file baru **main.go**.
2. Jalankan perintah berikut untuk menjalankan server:
   ```sh
   go run main.go
   ```
3. Akses endpoint melalui browser atau menggunakan `curl`:
   - Mulai server: `http://localhost:8080/start`
   - Stop setelah 10 hit: `http://localhost:8080/stop`
   - Reset hit: `http://localhost:8080/reload`

Setelah membaca ini, Anda bisa mulai mengembangkan server menggunakan Golang. Selamat mencoba!


## 3. Task Healthcheck

1. Healtcheck
2. Healtcheck to DB (Select date) message "Database connected with time"

## 4. Task Disk Healtcheck 

1. Create disk healtcheck
2. ednpoint method Get hit /disk-health