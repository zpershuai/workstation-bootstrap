# Manico Configuration Backup

This directory contains Manico app configuration backups.

## Files

- `com.lintie.manico.plist.xml` - Main app settings (shortcuts, preferences)
- `com.lintie.manico.helper.plist.xml` - Helper process settings
- `backup.sh` - Export current config to XML
- `restore.sh` - Import config from XML to system

## Usage

### Backup current configuration

```bash
bash misc/manico/backup.sh
```

### Restore configuration

```bash
bash misc/manico/restore.sh
```

## Notes

- Manico uses binary plist format, but we store as XML for version control
- XML format is human-readable and diff-friendly
- Restore script automatically restarts Manico to apply changes
