# Proyek GO

## Pendahuluan

Dokumen ini menjelaskan berbagai fitur dan metode yang diterapkan dalam proyek ini, termasuk layanan HTTP, manajemen layanan, health check, dan API untuk memperbarui data ID card.

---

## Metode

### Ringkasan HTTP Methods

| Metode   | Tujuan Utama                | Mengubah Resource? | Idempotent |
|----------|-----------------------------|--------------------|------------|
| PATCH    | Memperbarui sebagian data   | Ya                 | Tidak      |
| POST     | Membuat resource baru       | Ya                 | Tidak      |
| GET      | Membaca data                | Tidak              | Ya         |
| PUT      | Membuat/Mengganti resource  | Ya                 | Ya         |
| DELETE   | Menghapus resource          | Ya                 | Ya         |
| OPTIONS  | Memeriksa metode tersedia   | Tidak              | Ya         |
| HEAD     | Memeriksa metadata resource | Tidak              | Ya         |

Metode HTTP di atas digunakan sesuai dengan kebutuhan dan tindakan yang ingin dilakukan pada resource di server.

---

## Deskripsi Tugas

### 1. Layanan HTTP POST /refresh
Buat layanan HTTP di port **8080** yang menerima permintaan **POST /refresh** dan mengembalikan respons dengan status **"created"**.

#### Detail Fitur:
- Jika proses refresh sedang berjalan, permintaan baru harus **diantrikan**.
- Jika tidak ada proses refresh yang berjalan, proses refresh harus langsung dimulai.
- Proses refresh:
  - Memakan waktu sekitar **10 detik** (gunakan `sleep`).
  - Hindari mengantrikan permintaan dalam detik yang sama (berdasarkan timestamp Unix).
  - Catat timestamp kapan permintaan dilakukan.
- Proses refresh juga harus dijalankan sekali saat aplikasi pertama kali dijalankan.

---

### 2. Implementasi Service Manager

#### Fitur Utama
1. Cetak nilai integer (incremental) setiap **10 detik** saat server berjalan.
2. Endpoint **Start** untuk memulai pencetakan.
3. Endpoint **Stop** untuk menghentikan proses setelah **10 kali hit**.
4. Endpoint **Reload** untuk mengatur ulang nilai incremental menjadi **0**.

#### Cara Kerja
1. Nilai hit dimulai dari **0** saat server dijalankan.
2. Setiap **10 detik**, server akan mencetak **Hit ke-X**.
3. Server otomatis berhenti setelah mencapai **10 hit**.
4. Endpoint **Reload** akan mengatur ulang nilai hit menjadi **0**.

#### API Endpoint
| Method | Endpoint  | Deskripsi                    |
|--------|-----------|------------------------------|
| `POST` | `/start`  | Memulai proses pencetakan    |
| `POST` | `/stop`   | Menghentikan proses          |
| `POST` | `/reload` | Mereset nilai hit menjadi 0  |

#### Langkah-Langkah Menjalankan
1. Pastikan **Golang** sudah terinstal.
2. Jalankan server:
   ```sh
   go run main.go
   ```
3. Akses endpoint menggunakan browser atau alat seperti `curl`:
   - Mulai server: `http://localhost:8080/start`
   - Stop setelah 10 hit: `http://localhost:8080/stop`
   - Reset hit: `http://localhost:8080/reload`

---

### 3. Health Check

#### Fitur
1. **Health Check Server**:
   - Endpoint untuk memeriksa status server.
2. **Health Check Database**:
   - Melakukan query sederhana (`SELECT NOW()`) untuk memverifikasi koneksi ke database.
   - Mengembalikan pesan **"Database connected with time"** jika berhasil.

---

### 4. Disk Health Check

#### Fitur
- Endpoint untuk memeriksa status disk.
- Method: **GET**
- Endpoint: `/disk-health`

---

### 5. API Update ID Card

#### Deskripsi
API untuk memperbarui kolom `idcard` di tabel `profile` berdasarkan data ID KTP yang dikirim melalui permintaan **POST**.

#### Fitur Utama
- Validasi panjang ID KTP (harus **16 karakter**).
- Mengembalikan pesan sukses atau error sesuai kondisi.

#### Langkah-Langkah Instalasi
1. Tambahkan package berikut:
   ```bash
   go get -u github.com/go-sql-driver/mysql
   go get -u github.com/labstack/echo/v4
   ```

#### Contoh Output
```json
{
    "status": "success",
    "message": "ID card updated successfully"
}
```

---

### 6. Create Token with JWT on API

#### Pendahuluan
Anda diminta untuk membuat sebuah endpoint API yang berfungsi untuk menghasilkan token JWT. Token ini harus dienkripsi menggunakan kunci rahasia dan berisi payload tertentu yang dikirimkan melalui request body.

---

#### Tugas

#### 6.a. Endpoint API Pembuatan Token
Buatlah endpoint HTTP dengan spesifikasi berikut:
- **Method:** `POST`
- **Endpoint:** `/api/create-token`
- **Request Body:** Mengandung payload data (contoh: `userId` dan `role`).
- **Respons:** Mengembalikan token JWT yang valid.

#### 6.b. Spesifikasi Token
- Token harus berisi payload yang dikirimkan melalui request body.
- Token harus memiliki masa berlaku selama 1 jam.
- Token harus dienkripsi menggunakan algoritma HMAC-SHA256 dengan kunci rahasia.

---

#### Input & Output yang Diharapkan

#### Contoh Request
**HTTP Request:**

```http
POST /api/create-token HTTP/1.1
Content-Type: application/json

{
  "userId": 123,
  "role": "admin"
}
```
Contoh Response
HTTP Response (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Jika terdapat error: HTTP Response (400 Bad Request):


```json
{
  "error": "Failed to create token"
}
```


### 7. API Get Data via Parameter

#### Deskripsi Endpoint
API ini mengambil informasi ID card untuk pengguna tertentu berdasarkan parameter `idcard`.

#### Detail Endpoint
- **Query Param**: `idcard`
- **Method**: **GET**
- **Route**: `/get-users/:idcard`

#### Contoh Request
```http
GET http://localhost:8080/get-users?idcard=*** HTTP/1.1
Host: localhost:8080
```
---

### 8. API Authentication
#### Deskripsi
API ini mengelola autentikasi pengguna, termasuk Register, Login, dan Get Profile, menggunakan tabel authentication di basis data.

## Endpoints

#### 8.a. Register
- **Method**: `POST`  
- **Route**: `/auth/register`  
- **Body**:  
  ```json
  {
      "name": "John",
      "email": "john@example.com",
      "password": "123456",
      "photo_url": "http://example.com/photo.jpg"
  }
  ```
- **Response**:  
  - `201`: User registered.  
  - `400`: Email already exists.  

---

#### 8.b. Login
- **Method**: `POST`  
- **Route**: `/auth/login`  
- **Body**:  
  ```json
  {
      "email": "john@example.com",
      "password": "123456"
  }
  ```
- **Response**:  
  - `200`: Token returned.  
  - `401`: Invalid credentials.  

---

#### 8.c. Get Profile
- **Method**: `GET`  
- **Route**: `/auth/profile`  
- **Header**:  
  `Authorization: Bearer <token>`  
- **Response**:  
  - `200`: Profile data.  
  - `401`: Unauthorized.  

---

#### Instruksi
a. **Database**: Buat tabel `authentication`:
   ```sql
   CREATE TABLE authentication (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100) NOT NULL,
       email VARCHAR(50) NOT NULL UNIQUE,
       password VARCHAR(100) NOT NULL,
       photo_url TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   );
   ```
b. **Implementasi**: Gunakan hashing untuk menyimpan password dan JWT untuk autentikasi.
c. **Pengujian**: Gunakan Postman atau alat lainnya untuk menguji endpoint.


