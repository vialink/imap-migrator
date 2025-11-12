# Gmail IMAP Migration Guide

## 🔐 Gmail Requirements

To migrate Gmail accounts, you need:

1. **IMAP enabled** in Gmail settings
2. **2-Step Verification** enabled
3. **App Password** generated (16 characters)
4. Use **full email** as username (e.g., user@gmail.com)

## 📋 Step-by-Step Setup

### 1. Enable IMAP

1. Go to Gmail Settings (gear icon → See all settings)
2. Click **Forwarding and POP/IMAP** tab
3. In **IMAP Access** section, select **Enable IMAP**
4. Click **Save Changes**

### 2. Enable 2-Step Verification

1. Go to https://myaccount.google.com/security
2. Click **2-Step Verification**
3. Follow the setup wizard

### 3. Generate App Password

1. Go to https://myaccount.google.com/apppasswords
2. Select **Mail** and **Other (Custom name)**
3. Enter a name (e.g., "IMAP Migrator")
4. Click **Generate**
5. **Copy the 16-character password** (spaces don't matter)

## 📝 accounts.csv Format for Gmail

```csv
source_email,source_user,source_pass,source_host,dest_email,dest_user,dest_pass,dest_host
user@gmail.com,user@gmail.com,abcdabcdabcdabcd,imap.gmail.com,user@dest.com,user,password456,imap.dest.com
```

**Important:**
- Use **full email** as username: `user@gmail.com`
- Use **App Password** (16 chars), not regular password
- Server: `imap.gmail.com`
- Port: 993 (TLS) - automatic

## 🗂️ Gmail Folder Structure

Gmail uses special folder names:

- **[Gmail]/All Mail** - All messages archive
- **[Gmail]/Sent Mail** - Sent messages
- **[Gmail]/Drafts** - Draft messages
- **[Gmail]/Spam** - Spam folder
- **[Gmail]/Trash** - Trash folder
- **[Gmail]/Important** - Important messages
- **[Gmail]/Starred** - Starred messages

### Folder Mapping Example

If migrating to a standard IMAP server:

```json
{
  "folder_mapping": {
    "[Gmail]/Sent Mail": "INBOX.Sent",
    "[Gmail]/Drafts": "INBOX.Drafts",
    "[Gmail]/Spam": "INBOX.Junk",
    "[Gmail]/Trash": "INBOX.Trash",
    "[Gmail]/All Mail": "INBOX.Archive"
  }
}
```

## ⚠️ Important Notes

### All Mail Folder
- **[Gmail]/All Mail** contains ALL messages (including those in other folders)
- If you migrate this folder, you'll get duplicates
- **Recommendation**: Exclude it or use `skip_duplicates: true`

```json
{
  "skip_duplicates": true,
  "exclude_folders": ["[Gmail]/All Mail"]
}
```

### Labels vs Folders
- Gmail uses **labels**, not folders
- A message can have multiple labels
- In IMAP, it appears in multiple folders
- Use `skip_duplicates: true` to avoid copying the same message multiple times

### Rate Limits
- Gmail has IMAP rate limits
- The program handles this with automatic reconnection
- Large migrations may take time

## 🎯 Complete Example

### Gmail to Gmail
```csv
user1@gmail.com,user1@gmail.com,apppassword1234,imap.gmail.com,user2@gmail.com,user2@gmail.com,apppassword5678,imap.gmail.com
```

```json
{
  "skip_duplicates": true,
  "exclude_folders": ["[Gmail]/All Mail", "[Gmail]/Spam", "[Gmail]/Trash"]
}
```

### Gmail to Standard IMAP
```csv
user@gmail.com,user@gmail.com,apppassword1234,imap.gmail.com,user@dest.com,user,password456,imap.dest.com
```

```json
{
  "skip_duplicates": true,
  "exclude_folders": ["[Gmail]/All Mail"],
  "folder_mapping": {
    "[Gmail]/Sent Mail": "INBOX.Sent",
    "[Gmail]/Drafts": "INBOX.Drafts",
    "[Gmail]/Spam": "INBOX.Junk",
    "[Gmail]/Trash": "INBOX.Trash"
  }
}
```

## 🆘 Troubleshooting

### "Authentication failed"
- Check if you're using **App Password**, not regular password
- Verify 2-Step Verification is enabled
- Generate a new App Password

### "Too many simultaneous connections"
- Gmail limits concurrent connections
- Reduce parallel migrations in code (default is 5)

### "Quota exceeded"
- Gmail has storage limits
- Check available space in destination account

### Slow migration
- Normal for Gmail due to rate limits
- Be patient, the program will complete

## 📚 Additional Resources

- Gmail IMAP settings: https://support.google.com/mail/answer/7126229
- App Passwords: https://support.google.com/accounts/answer/185833
- 2-Step Verification: https://support.google.com/accounts/answer/185839

---

**Tip:** Always test with `dry_run: true` and a small folder first!
