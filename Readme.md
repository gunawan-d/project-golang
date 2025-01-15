# Project GO


## Description Task


### Deskripsi:
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
