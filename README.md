# 📂 Archive Scraper

`archive_scraper` is a command-line tool to **download all files** from an [Archive.org](https://archive.org/) item URL, with optional checksum verification and smart file management.

GitHub Repository: [jbdemonte/archive_scraper](https://github.com/jbdemonte/archive_scraper)

---

## ✅ Features
- **Download all files from an Archive.org item** 💽
- **Preserves folder structure** 🗂️
- **Automatic checksum verification (SHA1, MD5, CRC32)** 🔒
- **Avoids re-downloading existing files** by verifying integrity 📂
- **Supports skipping checksum verification with `--no-checksum`** 🚀
- **Graceful interruption handling (Ctrl+C)** ⏸️

---

## 🛠️ **Installation**
### **1️⃣ Build from source (Go required)**
```sh
git clone https://github.com/jbdemonte/archive_scraper.git
cd archive_scraper
go build -o archive_scraper
```

#### Build using the Makefile

##### Build to your system 
```shell
make build
```

##### Build to Synology NAS

Identify your Synology architecture

```shell
uname -m
```

###### Build for `x86_64` → PC-based (Intel/AMD)
```shell
make build_x86_64
```

###### Build for `armv7l` → ARM 32-bit
```shell
mke build_armv7l
```

###### Build for `aarch64` → ARM 64-bit
```shell
make build_aarch64
```

###### Build for `mipsel` → MIPS Little Endian
```shell
make build_mipsel
```


### **2️⃣ Run the executable**
```sh
./archive_scraper <url> [<path target>] [--no-checksum]
```

---

## 🚀 **Usage**
### **Basic Download**
```sh
./archive_scraper itemID
./archive_scraper https://archive.org/details/itemID
./archive_scraper https://archive.org/download/itemID
./archive_scraper https://archive.org/metadata/itemID
```

💽 Downloads all files from **`itemID`** into the **default directory** (_downloads_).

### **Download to a specific directory**
```sh
./archive_scraper <itemID> </my/download/path>
```
🗂️ Saves the downloaded files into **`/my/download/path`**.

### **Skip checksum verification**
```sh
./archive_scraper "https://archive.org/download/itemID" "/my/download/path" --no-checksum
```
⚡ **Faster download** by skipping checksum verification.

---

## 🛠️ **Options**
| **Option**       | **Description**                                                                      |
|-----------------|---------------------------------------------------------------------------------------|
| `<url>`         | The **Archive.org item URL** or **TokenID**                                           |
| `<path target>` | *(Optional)* Target folder for downloaded files. Defaults to a predefined directory.  |
| `--no-checksum` | *(Optional)* **Disables checksum verification** for faster downloads.                 |

---

## ⚠️ **Notes**
- **Existing files are verified and preserved** to avoid unnecessary re-downloads.
- **Interrupted downloads can be resumed**, saving time and bandwidth.
- Ensure you **have enough disk space** for large downloads.
- If **using a VPN**, Archive.org may block some requests. **Try disabling VPN** if downloads fail.
-- checksums may sometime be wrong in the metadata, if, so, retry using the `--no-checksum` options (i.e. - md5 may sometime be wrong (i.e. `arcadeflyer_donkey-kong_files.xml` in [arcadeflyer_donkey-kong](https://archive.org/metadata/arcadeflyer_donkey-kong))

---

## 📝 **License**
💜 Open-source license (MIT). Feel free to modify and contribute!

---

## ✨ **Contributing**
🛠️ PRs welcome! Report bugs or suggest features in the **issues section**.

---

## 📩 **Contact**
For questions, reach out via **GitHub Issues**.

🚀 **Happy downloading!** 🎉
