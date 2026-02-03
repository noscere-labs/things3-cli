# bear - A CLI for Bear Notes

A powerful command-line interface for [Bear](https://bear.app) that enables programmatic interaction with your notes. Perfect for automation, scripts, and integration with tools like Claude Code.

## Features

- ✅ **Create notes** with title, content, tags, and file attachments
- ✅ **Read notes** by ID or title with optional header extraction
- ✅ **Update notes** with append, prepend, or replace modes
- ✅ **Archive notes** to Bear's trash
- ✅ **Manage tags** - list, rename, and delete tags
- ✅ **Search notes** by content (with API token)
- ✅ **JSON output** for easy parsing and integration
- ✅ **Token management** for persistent API credentials

## Installation

### Prerequisites

- macOS (Bear is macOS-only)
- Go 1.21 or later
- Bear app installed

### Build from Source

```bash
# Clone or download the repository
cd bear-cli

# Build the binary
make build

# Install to /usr/local/bin/bear
make install

# Or install to ~/.local/bin without sudo
make install-user
```

After installation, verify it works:
```bash
bear --help
```

## Quick Start

### Create a Note

```bash
bear create --title "Meeting Notes" --content "Discussed Q1 roadmap" --tags "work,important"
```

Output:
```json
{
  "success": true,
  "data": {
    "note_id": "7E4B681B-...",
    "title": "Meeting Notes",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Read a Note

```bash
# By ID
bear read --id "7E4B681B-..."

# By title
bear read --title "Meeting Notes"
```

### Update a Note

```bash
# Append content (default)
bear update --id "7E4B681B-..." --content "New item added"

# Replace entire note
bear update --id "7E4B681B-..." --content "Completely new content" --mode replace_all

# Append with timestamp
bear update --id "7E4B681B-..." --content "Final notes" --timestamp
```

### List Notes by Tag

```bash
bear list --tag "work"
```

### Archive a Note

```bash
bear archive --id "7E4B681B-..."
```

## API Token Setup

Some operations (search, tags list) require your Bear API token.

### Get Your Token

1. Open Bear
2. Go to **Settings** → **Account** → **API Token**
3. Generate or copy your token

### Store Token

```bash
bear config set-token --token "your-api-token-here"
```

This securely stores your token in `~/.bear-cli/config.json` (permissions: 0600).

### View Token Status

```bash
bear config show
```

## Command Reference

### Create Command

```bash
bear create [OPTIONS]

Options:
  -t, --title STRING        Note title
  -c, --content STRING      Note content
  -g, --tags STRING         Comma-separated tags (e.g., "work,urgent")
  -f, --file PATH          File to attach
  -p, --pin                Pin note to top
      --timestamp          Prepend current date/time
```

### Read Command

```bash
bear read [OPTIONS]

Options:
  -i, --id STRING          Note ID
  -t, --title STRING       Note title
  -h, --header STRING      Extract specific header section
      --exclude-trashed    Skip trashed notes
```

### Update Command

```bash
bear update [OPTIONS]

Options:
  -i, --id STRING          Note ID (required)
  -c, --content STRING     Content to add/update
  -f, --file PATH          File to attach
  -m, --mode STRING        Update mode (default: append)
                           Values: append, prepend, replace, replace_all
  -h, --header STRING      Target specific header section
  -g, --tags STRING        Comma-separated tags to add
      --new-line          Add content on new line (append mode)
      --timestamp         Prepend date/time
```

### List Command

```bash
bear list [OPTIONS]

Options:
  -t, --tag STRING         Filter by tag
  -s, --search STRING      Search by content (requires token)
  -f, --filter STRING      Filter type (all, untagged, todo, today, locked)
  -k, --token STRING       API token (uses config if not provided)
```

### Archive Command

```bash
bear archive [OPTIONS]

Options:
  -i, --id STRING          Note ID (required)
      --no-window         Don't show Bear window
```

### Tags Command

```bash
# List all tags
bear tags list --token "YOUR_TOKEN"

# Rename a tag
bear tags rename --name "old-tag" --new-name "new-tag"

# Delete a tag
bear tags delete --name "old-tag"
```

### Config Command

```bash
# Set API token
bear config set-token --token "YOUR_TOKEN"

# View stored token (masked)
bear config get-token

# Show full configuration
bear config show
```

## JSON Output Format

All commands return JSON in this format:

### Success Response

```json
{
  "success": true,
  "data": {
    // Command-specific data
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Human-readable error message",
  "error_code": "MACHINE_READABLE_CODE",
  "details": "Additional technical details"
}
```

## Integration with Claude Code

The `bear` CLI is designed for seamless integration with Claude Code:

```javascript
// In Claude Code
const { execSync } = require('child_process');

function createNote(title, content, tags) {
  const cmd = `bear create --title "${title}" --content "${content}" --tags "${tags}"`;
  const output = execSync(cmd, { encoding: 'utf-8' });
  return JSON.parse(output);
}

function readNote(noteId) {
  const cmd = `bear read --id "${noteId}"`;
  const output = execSync(cmd, { encoding: 'utf-8' });
  return JSON.parse(output);
}

// Usage
const note = createNote("Claude's Analysis", "Key findings...", "ai,analysis");
console.log(note.data.note_id);
```

## Configuration

Configuration is stored in `~/.bear-cli/config.json`:

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

### Environment Variables

- `BEAR_TOKEN` - Override token from config
- `BEAR_CLI_CONFIG` - Custom config file path
- `BEAR_SHOW_WINDOW` - Force show Bear window (true/false)

## Troubleshooting

### "Note not found" Error

- Verify the note ID or title is correct
- Use `--exclude-trashed` if the note might be in trash

### "API token required" Error

- Set token with: `bear config set-token --token "YOUR_TOKEN"`
- Or pass with: `--token "YOUR_TOKEN"`

### Callback timeout errors

- Ensure Bear is running
- Check that port 8765 is available (or set a custom port in config)

### File attachment issues

- Use absolute paths or paths with `~` for home directory
- Ensure file exists and is readable

## Performance Notes

- Each command creates a temporary local HTTP server to receive callbacks
- Network latency depends on Bear's response time
- Large files may take longer to attach

## Limitations

- **macOS Only** - Bear is only available on macOS
- **No Direct Delete** - Notes cannot be deleted via CLI (only archived)
- **Encrypted Notes** - Cannot access encrypted notes without unlocking Bear
- **Token in Plaintext** - Token stored in `~/.bear-cli/config.json` (file permission 0600)
- **Batch Operations** - Each command is separate (no bulk operations yet)

## Development

### Build Commands

```bash
# Development build with verbose output
make dev

# Run linter
make lint

# Format code
make fmt

# Run tests
make test

# Show all make targets
make help
```

### Project Structure

```
bear-cli/
├── main.go                 # CLI entry point
├── cmd/commands.go         # All command implementations
├── pkg/
│   ├── bear/
│   │   ├── client.go       # URL scheme client
│   │   ├── callback.go     # HTTP callback server
│   │   └── models.go       # Data structures
│   ├── formatter/
│   │   └── json.go         # JSON output formatting
│   └── util/
│       ├── encoding.go     # URL/file encoding
│       └── config.go       # Configuration management
├── go.mod                  # Go module definition
├── Makefile                # Build automation
└── README.md              # This file
```

## Error Codes

- `INVALID_ARGUMENTS` - Missing or invalid command-line arguments
- `NOTE_NOT_FOUND` - Note doesn't exist
- `BEAR_ERROR` - Bear app returned an error
- `CALLBACK_TIMEOUT` - No response from Bear within timeout
- `CONFIG_ERROR` - Configuration file issue
- `CLIENT_ERROR` - Failed to initialize Bear client
- `FORMAT_ERROR` - Failed to format response

## Security Considerations

- **Token Storage** - Tokens are stored in `~/.bear-cli/config.json` with restricted permissions (0600)
- **Token Display** - `bear config get-token` masks tokens for security (shows first 6 and last 6 chars)
- **URL Encoding** - All parameters are properly URL-encoded to prevent injection
- **Local Callback** - Responses are received via localhost only, not exposed to network

## License

This project is provided as-is. Bear is a product of Shiny Frog.

## Contributing

Contributions welcome! Areas for expansion:
- Batch operations
- Additional filters and search options
- Better error handling for edge cases
- Performance optimizations
- Test coverage

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review the Bear x-callback-url documentation: https://bear.app/faq/x-callback-url-scheme-documentation/
3. Check your Bear installation is up to date
4. Verify API token is correct (if needed)

---

Made with ❤️ for Bear note enthusiasts and automation fans.
