# 🚀 Quick Start - IMAP Migrator v11

## 📦 Package Contents

This package contains everything you need to migrate IMAP email accounts.

### Code Files
- `main.go` - Main program
- `config.go` - Configuration system
- `duplicates.go` - Duplicate detection
- `report.go` - Report generation
- `go.mod` - Go dependencies

### Configuration Files
- `config.json` - Default configuration (use this)
- `config.json.exemplo` - Example with all options
- `accounts.csv` - **YOU NEED TO CREATE THIS FILE**
- `accounts.csv.exemplo` - Template for your accounts.csv

### Documentation
- `QUICKSTART.md` - This file (start here!)
- `FEATURES.md` - Complete features guide
- `README.md` - General documentation
- `GMAIL.md` - Gmail-specific guide
- `CHANGELOG.md` - Change history
- `REPORT_EXAMPLE.txt` - Sample generated report

---

## ⚡ Getting Started

### 1. **Extract Package**
```bash
unzip imap-migrator-v11-final.zip
cd imap-migrator-v11
```

### 2. **Create Accounts File**
Copy the example and edit with your accounts:
```bash
cp accounts.csv.exemplo accounts.csv
nano accounts.csv  # or use your preferred editor
```

**accounts.csv Format:**
```csv
source_email,source_user,source_pass,source_host,dest_email,dest_user,dest_pass,dest_host
user@source.com,user,password123,imap.source.com,user@dest.com,user,password456,imap.dest.com
```

### 3. **Build the Program**

**Using the build script (Recommended):**
```bash
chmod +x build.sh
./build.sh
```

The script will automatically:
- ✅ Check if Go is installed
- ✅ Install Go if needed (on Linux)
- ✅ Download all dependencies
- ✅ Compile the program
- ✅ Create `imap-migrator` executable

**Manual build (if you already have Go):**
```bash
go build -o imap-migrator *.go
```

### 4. **Configure Options (Optional)**
If you want to use filters, edit `config.json`:
```bash
nano config.json
```

**To start simple, leave it like this:**
```json
{
  "accounts_file": "accounts.csv",
  "skip_duplicates": false,
  "dry_run": false,
  "max_retries": 3,
  "max_message_size_mb": 0,
  "flatten_folders": false,
  "exclude_folders": [],
  "include_folders": [],
  "date_from": "",
  "date_to": "",
  "folder_mapping": {},
  "system_folders": {}
}
```

### 5. **Run the Program**

**If you built with build.sh:**
```bash
./imap-migrator
```

**Or run directly without building:**
```bash
go run *.go
```

### 6. **View Reports**
After migration, reports will be in:
```bash
ls reports/
cat reports/migration_*.txt
```

---

## 🧪 Test First (Recommended!)

Before migrating for real, do a test with **dry-run**:

1. Edit `config.json`:
```json
{
  "dry_run": true
}
```

2. Run:
```bash
go run *.go
```

3. See what would be done without copying anything!

---

## 📋 Quick Examples

### Simple Migration (Everything)
```json
{
  "skip_duplicates": false,
  "dry_run": false
}
```

### Migrate Only 2024
```json
{
  "skip_duplicates": true,
  "date_from": "2024-01-01",
  "date_to": "2024-12-31"
}
```

### Exclude Trash and Spam
```json
{
  "exclude_folders": [
    "INBOX.Trash",
    "INBOX.Junk",
    "INBOX.Drafts"
  ]
}
```

### Migrate Only Important Folders
```json
{
  "include_folders": [
    "INBOX",
    "INBOX.Important",
    "INBOX.Projects"
  ]
}
```

---

## ❓ Frequently Asked Questions

### How to run?
```bash
go run *.go
```

### Do I need to install anything?
Yes, only Go (version 1.19+). Dependencies are downloaded automatically.

### How to install Go?
- **Ubuntu/Debian**: `sudo apt install golang-go`
- **macOS**: `brew install go`
- **Windows**: Download from https://go.dev/dl/

### Is the program safe?
Yes! It:
- ✅ Tests all connections before starting
- ✅ Doesn't delete anything from source
- ✅ Generates detailed reports
- ✅ Supports dry-run for testing

### Can I stop in the middle?
Yes, use Ctrl+C. The program stops gracefully. You can run again and it will continue (use `skip_duplicates: true` to avoid recopying).

### How to migrate Gmail?
See the `GMAIL.md` file for specific instructions.

### Where are the reports?
In the `reports/` directory that is created automatically.

---

## 🆘 Common Problems

### "Quota exceeded"
The destination account is full. Increase storage limit.

### "Connection closed"
Normal! The program reconnects automatically.

### "Invalid credentials"
Check username and password in `accounts.csv`.

### For Gmail: "Authentication failed"
You need to use **App Password**, not your regular password. See `GMAIL.md`.

---

## 📚 Next Steps

1. ✅ Read `FEATURES.md` to see all advanced options
2. ✅ Test with `dry_run: true` first
3. ✅ Run the actual migration
4. ✅ Check reports in `reports/`

---

## 🎯 Complete Command

```bash
# 1. Extract
unzip imap-migrator-v11-final.zip
cd imap-migrator-v11

# 2. Build
chmod +x build.sh
./build.sh

# 3. Create accounts.csv
cp accounts.csv.exemplo accounts.csv
nano accounts.csv

# 4. Test (dry-run)
# Edit config.json and set "dry_run": true
./imap-migrator

# 5. Run for real
# Edit config.json and set "dry_run": false
./imap-migrator

# 6. View reports
ls reports/
```

---

## 💡 Final Tip

**Always start with dry-run!** This avoids surprises and allows you to validate your settings before copying thousands of messages.

Happy migrating! 🚀
