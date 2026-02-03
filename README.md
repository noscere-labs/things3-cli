# things - A CLI for Things 3

A command-line interface for [Things 3](https://culturedcode.com/things/) that uses the Things URL scheme to automate to-dos, projects, and lists.

## Features

- ✅ **Add to-dos** with notes, tags, checklists, and scheduling
- ✅ **Add projects** with optional areas and initial to-dos
- ✅ **Update to-dos/projects** (requires auth token)
- ✅ **Show lists or items** by query or ID
- ✅ **Search** Things from the command line
- ✅ **JSON payloads** for batch creation/update
- ✅ **JSON output** for easy scripting

## Installation

### Prerequisites

- macOS (Things is macOS-only)
- Go 1.21 or later
- Things 3 installed

### Build from Source

```bash
cd things3-cli

make build
make install
# or without sudo
make install-user
```

Verify:
```bash
things --help
```

## Quick Start

### Add a To-Do

```bash
things add --title "Buy milk" --when today --tags "errands"
```

### Add a Project

```bash
things add-project --title "Website" --area "Work" --to-dos "Design" --to-dos "Build"
```

### Update a To-Do (requires auth token)

```bash
things update --id "THINGS-ID" --title "Updated title" --reveal
```

### Show a List

```bash
things show --query Today
```

### Search

```bash
things search --query "invoice"
```

## Auth Token Setup

Updating items in Things requires an auth token.

1. Open Things
2. Go to **Settings → General**
3. Enable **Things URLs** and choose **Manage** to create a token

Store the token:
```bash
things config set-token --auth-token "YOUR_TOKEN"
```

You can also use the `THINGS_AUTH_TOKEN` environment variable.

## Command Reference (Highlights)

### Add To-Do

```bash
things add [OPTIONS]

Options:
  --title STRING
  --titles STRING (repeat flag)
  --notes STRING
  --when STRING
  --deadline STRING
  --tags STRING
  --list STRING
  --list-id STRING
  --heading STRING
  --heading-id STRING
  --checklist-items STRING (repeat flag)
  --completed
  --canceled
  --show-quick-entry
  --reveal
  --creation-date STRING
  --completion-date STRING
  --use-clipboard STRING
```

### Update To-Do

```bash
things update [OPTIONS]

Options:
  --id STRING (required)
  --title STRING
  --notes STRING
  --prepend-notes STRING
  --append-notes STRING
  --when STRING
  --deadline STRING
  --tags STRING
  --add-tags STRING
  --checklist-items STRING (repeat flag)
  --prepend-checklist-items STRING (repeat flag)
  --append-checklist-items STRING (repeat flag)
  --list STRING
  --list-id STRING
  --heading STRING
  --heading-id STRING
  --completed
  --canceled
  --reveal
  --duplicate
  --creation-date STRING
  --completion-date STRING
  --use-clipboard STRING
  --auth-token STRING
```

### JSON Payloads

```bash
things json --file payload.json
```

## Configuration

Config file location:

```
~/.config/things3-cli/config.json
```

Show config:
```bash
things config show
```
