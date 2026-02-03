# Bear CLI - Installation & Setup Guide

## Overview

You have a complete, production-ready Go CLI tool for Bear notes. The executable will be named `bear` (not `bear-cli`).

## System Requirements

- **OS**: macOS (10.12 or later)
- **Go**: Version 1.21 or later
- **Bear**: Latest version of Bear app installed

Check Go installation:
```bash
go version
```

## Installation Steps

### Step 1: Verify the Project Structure

After extracting the archive, you should have:

```
bear-cli/
├── main.go                 # Entry point (uses cobra for CLI)
├── cmd/commands.go         # All command implementations with full docs
├── pkg/
│   ├── bear/
│   │   ├── client.go       # Bear x-callback-url client (150+ lines)
│   │   ├── callback.go     # HTTP callback server for responses
│   │   └── models.go       # Data structures for notes/tags
│   ├── formatter/
│   │   └── json.go         # JSON output formatting
│   └── util/
│       ├── encoding.go     # URL encoding and file utilities
│       └── config.go       # Config file management
├── go.mod                  # Go module (uses Cobra framework)
├── Makefile                # Build automation
├── README.md               # Full documentation
└── .gitignore              # Git ignore rules
```

### Step 2: Navigate to Project

```bash
cd bear-cli
```

### Step 3: Build the Executable

The simplest approach:

```bash
make build
```

This creates `./bin/bear` binary.

Or manually:

```bash
mkdir -p bin
go build -o bin/bear main.go
```

### Step 4: Install to PATH

#### Option A: Install System-Wide (Recommended)

```bash
make install
```

This installs to `/usr/local/bin/bear` and creates the executable.

If prompted for sudo password, enter it. This allows you to use `bear` command globally.

#### Option B: Install to User Home (No Sudo)

```bash
make install-user
```

This installs to `~/.local/bin/bear`. 

**You must add `~/.local/bin` to your PATH if not already present:**

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Or for zsh:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Step 5: Verify Installation

```bash
bear --help
```

You should see help text for the bear command.

## First Run: Set API Token

Many operations require a Bear API token. Get yours:

1. Open **Bear** app
2. Go to **Settings** → **Account** → **API Token**
3. Copy your token

Store it:

```bash
bear config set-token --token "YOUR_API_TOKEN_HERE"
```

Verify:

```bash
bear config show
```

## Quick Test

### Test 1: Create a Note

```bash
bear create --title "Test Note" --content "Hello from CLI!" --tags "test"
```

Expected output:
```json
{
  "success": true,
  "data": {
    "note_id": "...",
    "title": "Test Note",
    "created_at": "..."
  }
}
```

### Test 2: Read the Note (use the note_id from above)

```bash
bear read --id "..." 
```

Or by title:

```bash
bear read --title "Test Note"
```

### Test 3: Update the Note

```bash
bear update --id "..." --content "Updated content" --timestamp
```

### Test 4: List Tags

```bash
bear list --tag "test"
```

## Common Operations

### Create a note with timestamp

```bash
bear create --title "Daily Log" --content "Today's summary" --timestamp
```

### Append to a note

```bash
bear update --id "NOTE_ID" --content "New entry" --timestamp --new-line
```

### Replace entire note

```bash
bear update --id "NOTE_ID" --content "Completely new content" --mode replace_all
```

### Archive (trash) a note

```bash
bear archive --id "NOTE_ID"
```

### Get all tags (requires API token)

```bash
bear tags list --token "YOUR_TOKEN"
```

Or if token is stored:

```bash
bear tags list
```

## Integration with Claude Code

The CLI outputs JSON, perfect for parsing in scripts:

```javascript
// In Claude Code
const { execSync } = require('child_process');

// Create a note
const result = execSync('bear create --title "Test" --content "Content"', { encoding: 'utf-8' });
const parsed = JSON.parse(result);
const noteId = parsed.data.note_id;

// Read it back
const readResult = execSync(`bear read --id "${noteId}"`, { encoding: 'utf-8' });
const note = JSON.parse(readResult);
console.log(note.data.content);
```

## Make Targets

Useful make commands:

```bash
make help          # Show all targets
make build         # Build the binary
make install       # Install to /usr/local/bin
make install-user  # Install to ~/.local/bin
make clean         # Remove build artifacts
make fmt           # Format code
make dev           # Development build with verbose output
```

## Configuration

Config file location: `~/.bear-cli/config.json`

View current config:

```bash
bear config show
```

Typical config after setup:

```json
{
  "token": "your-api-token",
  "callback_port": 8765,
  "callback_timeout_seconds": 10,
  "show_window": false,
  "output_format": "json",
  "last_updated": "2024-01-15T10:30:00Z"
}
```

### Important Notes

- Token is stored with restricted permissions (0600, user-readable only)
- Callback port defaults to 8765 (used for receiving responses from Bear)
- Timeout is 10 seconds (configurable, increase if Bear is slow)
- `show_window: false` means Bear won't appear on screen during operations

## Troubleshooting

### "bear: command not found"

The binary isn't in your PATH. Try:

```bash
# Test if it works with full path
./bin/bear --help

# Then reinstall to PATH
make install
# or
make install-user
```

### "Failed to read file" error

File doesn't exist or isn't readable:

```bash
# Use absolute path
bear create --title "Test" --file ~/Downloads/document.pdf

# Or check file exists
ls -la ~/Downloads/document.pdf
```

### "API token required" error

For search and tags operations:

```bash
# Set token once
bear config set-token --token "YOUR_TOKEN"

# Or pass each time
bear tags list --token "YOUR_TOKEN"
```

### Callback timeout errors

1. Ensure Bear app is running
2. Try again (sometimes takes a moment)
3. Check port 8765 isn't in use: `lsof -i :8765`

### "Note not found"

- Verify note ID is correct
- Try reading by title instead: `bear read --title "Title"`
- Check note isn't in trash: `bear read --title "Title" --exclude-trashed`

## Uninstall

To remove the tool:

```bash
# Remove installation
rm /usr/local/bin/bear
# or
rm ~/.local/bin/bear

# Remove configuration
rm -rf ~/.bear-cli

# Clean up source
rm -rf bear-cli/
```

## Next Steps

1. **Read the main README.md** for complete command reference
2. **Review the code** - heavily commented for learning:
   - `pkg/bear/client.go` - Bear URL scheme handling
   - `cmd/commands.go` - All CLI commands
   - `pkg/util/config.go` - Configuration management
3. **Create an alias** for convenience:
   ```bash
   alias b='bear'
   ```
4. **Integrate with scripts** - the JSON output makes it perfect for automation
5. **Check Bear API docs** - https://bear.app/faq/x-callback-url-scheme-documentation/

## Support & Help

View help for any command:

```bash
bear --help
bear create --help
bear read --help
bear update --help
bear list --help
bear tags --help
bear config --help
bear archive --help
```

Each command has detailed help with examples.

## Code Quality

The code includes:

- ✅ Extensive inline comments explaining implementation
- ✅ Clear function documentation for every exported function
- ✅ Proper error handling with meaningful error codes
- ✅ Configuration management with safe defaults
- ✅ JSON output formatting for easy integration
- ✅ Modular design with clear separation of concerns

All code is production-ready and fully documented!

---

**Enjoy using the bear CLI!** 🐻
