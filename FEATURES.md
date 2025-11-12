# Advanced Features - IMAP Migrator

## 📋 Complete Feature List

### ✅ Implemented Features

#### 1. **Duplicate Detection**
- Checks if messages already exist at destination before copying
- Uses Message-ID as unique identifier
- Fallback to MD5 hash (subject + sender + date + size) when Message-ID unavailable
- Configurable via `skip_duplicates` in config.json

#### 2. **Folder Filter - Exclusion**
- Exclude specific folders from migration
- Useful for skipping Drafts, Trash, Junk, etc.
- Configurable via `exclude_folders` in config.json

#### 3. **Folder Filter - Inclusion (Whitelist)**
- Migrate ONLY specific folders
- When configured, all other folders are ignored
- Configurable via `include_folders` in config.json

#### 4. **Date Filter - From**
- Migrate only messages >= specified date
- Format: YYYY-MM-DD (e.g., 2024-01-15)
- Configurable via `date_from` in config.json

#### 5. **Date Filter - To**
- Migrate only messages <= specified date
- Format: YYYY-MM-DD (e.g., 2015-10-10)
- Configurable via `date_to` in config.json

#### 6. **Folder Name Mapping**
- Rename folders during migration
- Useful for compatibility between different naming conventions
- Example: "INBOX.Sent Messages" → "INBOX.Sent"
- Configurable via `folder_mapping` in config.json

#### 7. **Message Size Limit**
- Skip messages larger than X MB
- Useful when destination server has size limitations
- 0 = no limit
- Configurable via `max_message_size_mb` in config.json

#### 8. **Dry-Run Mode (Simulation)**
- Runs without actually copying
- Shows exactly what would be done
- Useful for validating filters before execution
- Configurable via `dry_run` in config.json

#### 9. **Automatic Retry**
- Retries failed messages
- Configurable number of attempts
- Useful for handling temporary network errors
- Configurable via `max_retries` in config.json

#### 10. **Flatten Folder Hierarchy**
- Converts "INBOX.Sent.2024" to "INBOX_Sent_2024"
- Useful when destination server has hierarchy limitations
- Configurable via `flatten_folders` in config.json

#### 11. **Configurable System Folders**
- Define alternative names for Drafts, Sent, Junk, Trash, Archive
- Supports multiple names (Gmail, Outlook, etc.)
- Configurable via `system_folders` in config.json

---

## 🔧 Configuration File (config.json)

```json
{
  "accounts_file": "accounts.csv",
  "skip_duplicates": true,
  "dry_run": false,
  "max_retries": 3,
  "max_message_size_mb": 50,
  "flatten_folders": false,
  
  "exclude_folders": [
    "INBOX.Drafts",
    "INBOX.Trash",
    "INBOX.Junk"
  ],
  
  "include_folders": [],
  
  "date_from": "2024-01-01",
  "date_to": "2024-12-31",
  
  "folder_mapping": {
    "INBOX.Sent Messages": "INBOX.Sent",
    "INBOX.Deleted Items": "INBOX.Trash"
  },
  
  "system_folders": {
    "drafts": ["Drafts", "INBOX.Drafts", "[Gmail]/Drafts"],
    "sent": ["Sent", "Sent Messages", "INBOX.Sent", "[Gmail]/Sent Mail"],
    "junk": ["Junk", "Spam", "INBOX.Junk", "[Gmail]/Spam"],
    "trash": ["Trash", "Deleted Items", "INBOX.Trash", "[Gmail]/Trash"],
    "archive": ["Archive", "INBOX.Archive", "[Gmail]/All Mail"]
  }
}
```

---

## 📖 Usage Examples

### Example 1: Simple Migration (No Filters)
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

### Example 2: Migrate Only 2024, No Trash
```json
{
  "skip_duplicates": true,
  "exclude_folders": ["INBOX.Trash", "INBOX.Junk", "INBOX.Drafts"],
  "date_from": "2024-01-01",
  "date_to": "2024-12-31"
}
```

### Example 3: Migrate Only Specific Folders
```json
{
  "skip_duplicates": true,
  "include_folders": [
    "INBOX",
    "INBOX.Important",
    "INBOX.Projects"
  ]
}
```

### Example 4: Test (Dry-Run)
```json
{
  "dry_run": true,
  "skip_duplicates": true,
  "date_from": "2024-01-01"
}
```

### Example 5: Migration with Size Limit
```json
{
  "skip_duplicates": true,
  "max_message_size_mb": 25,
  "max_retries": 5
}
```

---

## ⚙️ Additive Rules

All rules are **additive** (AND logic). A message is only copied if it passes ALL filters:

1. ✅ Folder is in whitelist (if configured)
2. ✅ Folder is NOT in blacklist
3. ✅ Date >= date_from (if configured)
4. ✅ Date <= date_to (if configured)
5. ✅ Size <= max_message_size_mb (if configured)
6. ✅ Not a duplicate (if skip_duplicates = true)

**Example:**
```json
{
  "include_folders": ["INBOX", "INBOX.Important"],
  "date_from": "2015-10-11",
  "date_to": "2024-01-14",
  "max_message_size_mb": 50
}
```

Result: Migrates only messages from INBOX and INBOX.Important, between 10/11/2015 and 01/14/2024, smaller than 50MB.

---

## 📊 Reports

Reports now include:
- Messages skipped by date filter
- Messages skipped by size
- Messages skipped by duplication
- Specific reason for each skipped message

---

## 🚀 How to Use

1. Edit `config.json` with your preferences
2. Run: `go run *.go`
3. Check report in `reports/`

---

## ⚠️ Important Notes

- **Folder Validation**: If using `include_folders` or `exclude_folders`, the program checks if all folders exist before starting
- **Dry-Run**: Always test with `dry_run: true` first
- **Duplicates**: Detection uses Message-ID, which is reliable but not 100% guaranteed
- **Performance**: Duplicate detection adds overhead (searches existing messages)
