# Bear CLI - MURF TTS Setup Guide

## Overview

The `bear speak` command converts Bear notes to speech using MURF AI's text-to-speech API.

## Prerequisites

- Node.js installed (`node --version`)
- MURF API key from https://murf.ai/api/dashboard

## Installation & Setup

### 1. Copy MURF CLI Script (Global Install)

```bash
mkdir -p ~/.config/bear-cli
cp murf/murf-cli.js ~/.config/bear-cli/murf-cli.js
chmod +x ~/.config/bear-cli/murf-cli.js
```

### 2. Configure MURF Settings

**Option A: Using .env file (Recommended)**

```bash
cp murf/.env ~/.config/bear-cli/.env
```

Then edit `~/.config/bear-cli/.env` and set your API key:

```bash
MURF_API_KEY=your-api-key-here
MURF_VOICE_ID=en-UK-mason
MURF_FORMAT=MP3
MURF_SAMPLE_RATE=24000
MURF_OUTPUT_DIR=~/.config/bear-cli/audio
MURF_AUTO_PLAY=false
```

**Option B: Using CLI Commands**

```bash
bear config set-murf --api-key "your-api-key"
bear config set-murf --voice "en-US-sara"
bear config set-murf --format "MP3"
```

### 3. Verify Configuration

```bash
bear config show-murf
```

You should see your settings with the API key masked.

## Usage

### Basic Usage

Convert a note by ID:
```bash
bear speak --id "7E4B681B-..."
```

Convert a note by title:
```bash
bear speak --title "Meeting Notes"
```

### Options

```bash
# Override voice for specific call
bear speak --id "UUID" --voice "en-US-emma"

# Auto-play audio after generation
bear speak --id "UUID" --play

# Save to custom location
bear speak --id "UUID" --output ~/my-audio.mp3

# Extract specific section
bear speak --id "UUID" --header "Summary"

# Check current configuration
bear config show-murf
```

## Available Voices

Common voices available through MURF:
- `en-UK-mason` (default) - British English male
- `en-US-sara` - American English female
- `en-GB-william` - British English male
- `en-US-emma` - American English female
- `en-AU-karen` - Australian English female

See full list: https://murf.ai/api/docs/voices

## Audio Formats

Supported audio formats:
- `MP3` (default) - Good compression, widely compatible
- `WAV` - Uncompressed, better quality
- `FLAC` - Lossless compression
- `OGG` - Vorbis compression

## Output Location

Audio files are saved to `~/.config/bear-cli/audio/` by default with names like:
```
bear-tts-2026-01-29T14-30-00-abc123.mp3
```

## Text Processing

The text from your Bear note is automatically processed:
- Code blocks removed: ` ``` code ``` ` → "[code block omitted]"
- Inline code cleaned: `` `code` `` → "[code]"
- Markdown formatting stripped
- Bear wiki links converted: `[[Link]]` → "Link"
- Multiple newlines collapsed for readability
- Text truncated to max 5000 characters if needed

## Troubleshooting

### "ERROR: MURF_API_KEY not configured"

Make sure your API key is set in either:
1. `~/.config/bear-cli/.env` file, or
2. Environment variable: `export MURF_API_KEY=your-key`
3. Via CLI: `bear config set-murf --api-key "your-key"`

### "TTS_SCRIPT_NOT_FOUND"

Ensure `murf-cli.js` is installed:
```bash
cp murf/murf-cli.js ~/.config/bear-cli/murf-cli.js
```

### "Note not found"

Make sure:
- The note ID is correct (visible in Bear app)
- Or use `--title` if you know the note title
- The note isn't in trash

### Audio quality issues

Try adjusting in `~/.config/bear-cli/.env`:
- Increase `MURF_SAMPLE_RATE` to 44100 or 48000 (higher = better quality)
- Change `MURF_FORMAT` to WAV or FLAC for better quality (larger file size)

### Network/API errors

- Check your internet connection
- Verify MURF_API_KEY is valid at https://murf.ai/api/dashboard
- Check MURF API status at https://murf.ai/

## Configuration Priority

Settings are loaded in this order (first match wins):
1. Command-line flags (`--voice`, `--format`, etc.)
2. `~/.config/bear-cli/.env` file
3. Environment variables (`MURF_API_KEY`, `MURF_VOICE_ID`, etc.)
4. `~/.config/bear-cli/config.json` (saved via `bear config set-murf`)
5. Built-in defaults

## Advanced: Custom Script Path

If you want to use a custom MURF script:

```bash
export MURF_CLI_SCRIPT=/path/to/custom/murf-cli.js
bear speak --id "UUID"
```

## Sample Workflow

```bash
# 1. Setup (one time)
cp murf/murf-cli.js ~/.config/bear-cli/murf-cli.js
cp murf/.env ~/.config/bear-cli/.env
# Edit .env with your API key

# 2. Verify setup
bear config show-murf

# 3. Use
bear speak --title "Q1 Roadmap" --play

# 4. Locate audio
ls ~/.config/bear-cli/audio/
```

## Cost Estimation

MURF charges per character processed. Typical costs:
- Small note (100 chars): ~$0.001
- Medium article (2,000 chars): ~$0.02
- Long document (5,000 chars): ~$0.05

No fees for failed requests.

## Support

- MURF API Docs: https://murf.ai/api/docs
- MURF Dashboard: https://murf.ai/api/dashboard
- Report issues in your MURF account
