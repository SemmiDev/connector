#!/bin/bash

# Script untuk menginstal Go (Golang) versi terbaru di Ubuntu

set -e

echo "Memulai instalasi Golang terbaru..."

# 1. Tentukan URL untuk versi terbaru
#GO_LATEST=$(curl -s https://go.dev/VERSION?m=text)
GO_LATEST=go1.24.0
GO_URL="https://go.dev/dl/${GO_LATEST}.linux-amd64.tar.gz"

# 2. Unduh Golang terbaru
echo "Mengunduh Golang versi $GO_LATEST..."
curl -OL "$GO_URL"

# 3. Hapus instalasi Golang sebelumnya (jika ada)
echo "Menghapus instalasi Golang sebelumnya..."
rm -rf /usr/local/go

# 4. Ekstrak dan pasang Golang
echo "Ekstrak dan memasang Golang..."
tar -C /usr/local -xzf "${GO_LATEST}.linux-amd64.tar.gz"

# 5. Hapus file unduhan
echo "Membersihkan file unduhan..."
rm "${GO_LATEST}.linux-amd64.tar.gz"

# 6. Tambahkan Go ke PATH untuk semua user
if ! grep -q "export PATH=\$PATH:/usr/local/go/bin" /etc/profile; then
  echo "Menambahkan Go ke PATH di /etc/profile..."
  echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
fi

# Terapkan perubahan PATH
source /etc/profile

# 7. Verifikasi instalasi
echo "Verifikasi instalasi Golang..."
go version

echo "Instalasi Golang versi $GO_LATEST selesai!"

# 8. Hapus file unduhan
rm "${GO_LATEST}.linux-amd64.tar.gz"
