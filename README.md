# IMAP Email Migration Tool

A robust and feature-rich tool for migrating email accounts between IMAP servers, preserving all folders, messages, and their status.

## ✨ Key Features

- **Parallel Processing**: Migrate up to 5 accounts simultaneously
- **Connection Pre-Check**: Tests all connections before starting migration
- **Automatic Reconnection**: Handles connection drops gracefully
- **Duplicate Detection**: Skip already migrated messages (by Message-ID)
- **Advanced Filtering**: Filter by date, size, folders
- **Dry-Run Mode**: Test migration without copying anything
- **Detailed Reports**: Comprehensive audit reports for each account
- **Gmail Support**: Full support for Gmail IMAP
- **Retry Logic**: Configurable automatic retry for failed messages
- **Folder Mapping**: Rename folders during migration
- **Quota Handling**: Graceful handling of quota exceeded errors

## 🚀 Quick Start

### 1. Extract the Package
```bash
unzip imap-migrator-v11-final.zip
cd imap-migrator-v11
```

### 2. Create Accounts File
```bash
cp accounts.csv.exemplo accounts.csv
nano accounts.csv
```

**Format:**
```csv
source_email,source_user,source_pass,source_host,dest_email,dest_user,dest_pass,dest_host
user@source.com,user,pass123,imap.source.com,user@dest.com,user,pass456,imap.dest.com
```

### 3. Build the Program

**Option A: Using the build script (Recommended)**
```bash
chmod +x build.sh
./build.sh
```
The script will:
- Check if Go is installed (install if needed)
- Download all dependencies
- Compile the program
- Create `imap-migrator` executable

**Option B: Manual build**
```bash
go build -o imap-migrator *.go
```

### 4. Configure Options (Optional)
Edit `config.json` to set filters and options.

### 5. Run
```bash
./imap-migrator
```

Or without building:
```bash
go run *.go
```

## 📋 Configuration Options

All options are configured in `config.json`:

- **accounts_file**: CSV file with accounts (default: "accounts.csv")
- **skip_duplicates**: Skip already migrated messages
- **dry_run**: Simulate migration without copying
- **max_retries**: Number of retry attempts for failed messages
- **max_message_size_mb**: Skip messages larger than X MB
- **flatten_folders**: Convert folder hierarchy to flat names
- **exclude_folders**: Blacklist of folders to skip
- **include_folders**: Whitelist of folders to migrate (if set, only these are migrated)
- **date_from**: Migrate only messages >= this date (YYYY-MM-DD)
- **date_to**: Migrate only messages <= this date (YYYY-MM-DD)
- **folder_mapping**: Rename folders during migration
- **system_folders**: Alternative names for system folders

## 📚 Documentation

- **QUICKSTART.md** - Quick start guide
- **FEATURES.md** - Complete features documentation
- **GMAIL.md** - Gmail-specific guide
- **CHANGELOG.md** - Version history
- **REPORT_EXAMPLE.txt** - Sample migration report

**Portuguese versions available with `-ptbr` suffix.**

## 🔧 Requirements

- Go 1.19 or higher
- Internet connection
- IMAP access to source and destination servers

## 📊 Reports

After migration, detailed reports are generated in the `reports/` directory with:
- Migration summary (duration, totals)
- Per-folder statistics
- List of errors (if any)
- Skipped messages with reasons

## ⚠️ Important Notes

- Always test with `dry_run: true` first
- The tool does NOT delete anything from the source
- For Gmail, use App Passwords (see GMAIL.md)
- All filters are additive (AND logic)

## 🆘 Support

For issues, questions, or feature requests, please visit:
https://help.manus.im

## 📄 License

This tool is provided as-is for email migration purposes.

---

**Version 11** - Full-featured IMAP migration tool with advanced filtering and reporting.
