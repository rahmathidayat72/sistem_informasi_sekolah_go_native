# 🔧 Troubleshooting - Mengatasi Masalah Git Push

## ❌ Masalah: Permission Denied saat Git Push

Jika Anda mengalami error seperti ini:

```
git@github.com: Permission denied (publickey).
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
and the repository exists.
```

### 🔍 Penyebab Masalah

Error ini terjadi karena Anda mencoba menggunakan **SSH authentication** untuk mengakses GitHub, tetapi:
1. SSH key belum dikonfigurasi di komputer Anda, ATAU
2. SSH key belum ditambahkan ke akun GitHub Anda

### ✅ Solusi 1: Menggunakan HTTPS dengan Personal Access Token (Recommended)

Ini adalah cara termudah dan paling cepat:

#### Langkah 1: Ubah Remote URL ke HTTPS

```bash
git remote set-url origin https://github.com/rahmathidayat72/sistem_informasi_sekolah_go_native.git
```

Verifikasi perubahan:
```bash
git remote -v
```

Output seharusnya:
```
origin  https://github.com/rahmathidayat72/sistem_informasi_sekolah_go_native.git (fetch)
origin  https://github.com/rahmathidayat72/sistem_informasi_sekolah_go_native.git (push)
```

#### Langkah 2: Buat Personal Access Token (PAT)

1. Login ke GitHub → Klik foto profil Anda → **Settings**
2. Scroll ke bawah → Klik **Developer settings** (di bagian paling bawah)
3. Klik **Personal access tokens** → **Tokens (classic)**
4. Klik **Generate new token** → **Generate new token (classic)**
5. Beri nama token (contoh: "Sistem Sekolah Laptop")
6. Pilih expiration (masa berlaku token)
7. Centang scope minimal:
   - ✅ **repo** (Full control of private repositories)
8. Scroll ke bawah → Klik **Generate token**
9. **PENTING**: Salin token yang muncul (Anda tidak akan bisa melihatnya lagi!)

#### Langkah 3: Push dengan Token

Ketika melakukan push pertama kali:

```bash
git push -u origin main
```

GitHub akan meminta username dan password:
- **Username**: username GitHub Anda
- **Password**: **PASTE TOKEN** yang sudah Anda salin (bukan password akun GitHub!)

Setelah sekali berhasil, Git akan menyimpan kredensial Anda.

---

### ✅ Solusi 2: Setup SSH Key (Untuk Pengguna Advanced)

Jika Anda ingin tetap menggunakan SSH:

#### Langkah 1: Cek SSH Key yang Ada

```bash
ls -al ~/.ssh
```

Cari file seperti `id_rsa.pub`, `id_ed25519.pub`, atau `id_ecdsa.pub`.

#### Langkah 2: Generate SSH Key Baru (jika belum ada)

```bash
ssh-keygen -t ed25519 -C "email@anda.com"
```

Atau jika sistem tidak support ed25519:
```bash
ssh-keygen -t rsa -b 4096 -C "email@anda.com"
```

Tekan **Enter** untuk semua pertanyaan (menggunakan default).

#### Langkah 3: Start SSH Agent

```bash
eval "$(ssh-agent -s)"
```

#### Langkah 4: Tambahkan SSH Key ke Agent

```bash
ssh-add ~/.ssh/id_ed25519
```

Atau jika menggunakan RSA:
```bash
ssh-add ~/.ssh/id_rsa
```

#### Langkah 5: Salin Public Key

**Linux/Mac:**
```bash
cat ~/.ssh/id_ed25519.pub
```

**Windows (Git Bash):**
```bash
clip < ~/.ssh/id_ed25519.pub
```

Salin output yang muncul (dimulai dengan `ssh-ed25519` atau `ssh-rsa`).

#### Langkah 6: Tambahkan SSH Key ke GitHub

1. Login ke GitHub → Klik foto profil → **Settings**
2. Klik **SSH and GPG keys** (di sidebar kiri)
3. Klik **New SSH key**
4. Beri judul (contoh: "Laptop Pribadi")
5. Paste public key yang sudah disalin
6. Klik **Add SSH key**

#### Langkah 7: Test Koneksi SSH

```bash
ssh -T git@github.com
```

Jika berhasil, akan muncul:
```
Hi username! You've successfully authenticated, but GitHub does not provide shell access.
```

#### Langkah 8: Ubah Remote URL ke SSH

```bash
git remote set-url origin git@github.com:rahmathidayat72/sistem_informasi_sekolah_go_native.git
```

#### Langkah 9: Push

```bash
git push -u origin main
```

---

### ✅ Solusi 3: Clone Ulang Repository dengan HTTPS

Jika solusi di atas tidak berhasil, clone ulang repository:

```bash
# Backup folder saat ini
cd ..
mv sistem_informasi_sekolah_go_native sistem_informasi_sekolah_go_native_backup

# Clone ulang dengan HTTPS
git clone https://github.com/rahmathidayat72/sistem_informasi_sekolah_go_native.git

# Pindahkan file yang sudah dimodifikasi (jika ada)
# cp sistem_informasi_sekolah_go_native_backup/file_anda.go sistem_informasi_sekolah_go_native/
```

---

## 🆘 Masalah Lain yang Sering Terjadi

### Problem: "fatal: remote origin already exists"

**Solusi:**
```bash
git remote remove origin
git remote add origin https://github.com/rahmathidayat72/sistem_informasi_sekolah_go_native.git
```

### Problem: "Updates were rejected because the tip of your current branch is behind"

**Solusi:**
```bash
git pull origin main --rebase
# Jika ada konflik, selesaikan secara manual:
# 1. Edit file yang konflik
# 2. git add .
# 3. git rebase --continue
git push origin main
```

### Problem: "fatal: refusing to merge unrelated histories"

**Solusi:**
```bash
git pull origin main --allow-unrelated-histories
```

---

## 📞 Butuh Bantuan?

Jika masih mengalami masalah:
1. Buat issue baru di repository ini
2. Sertakan:
   - Output lengkap dari error yang muncul
   - Hasil dari `git remote -v`
   - Solusi mana yang sudah dicoba

---

## 📚 Referensi

- [GitHub Docs - Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [GitHub Docs - Connecting to GitHub with SSH](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)
- [Git Documentation - git remote](https://git-scm.com/docs/git-remote)
