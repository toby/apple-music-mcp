# 🎵 Apple Music MCP Server

A [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol/docs) server that enables Large Language Models to seamlessly interact with your Apple Music library. Built with Go and designed for the modern AI-powered music discovery experience.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![MCP Compatible](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://github.com/modelcontextprotocol/docs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 🎯 What is this?

Imagine asking an AI assistant to "play my recently added jazz albums" or "create a playlist of upbeat songs from the 2000s" and having it actually work with your Apple Music account. That's what this MCP server makes possible.

By implementing the Model Context Protocol, this server acts as a bridge between AI language models (like Claude, GPT-4, or local models) and the Apple Music API, enabling natural language interactions with your music library.

## ✨ Features

### 🎶 Library Management
- **Search & Discovery**: Find tracks, albums, artists, and playlists
- **Library Access**: Browse your personal music collection
- **Playlist Management**: Create, modify, and organize playlists
- **Recently Played**: Access your listening history
- **Recommendations**: Get personalized music suggestions

### 🤖 AI Integration
- **Natural Language**: Interact with your music using conversational commands
- **Context Awareness**: AI understands your music preferences over time
- **Smart Queries**: Complex searches like "indie rock from artists I've liked recently"
- **Batch Operations**: Perform multiple actions in a single conversation

### 🔒 Security & Privacy
- **OAuth 2.0**: Secure authentication with Apple Music
- **Token Management**: Automatic refresh and secure storage
- **Permission Scoping**: Only access what's needed
- **Local Storage**: Your tokens stay on your machine

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   LLM Client    │◄──►│  MCP Server      │◄──►│  Apple Music    │
│                 │    │                  │    │     API         │
│ • Claude        │    │ • Tool Registry  │    │                 │
│ • GPT-4         │    │ • Auth Handler   │    │ • User Library  │
│ • Local Models  │    │ • Request Router │    │ • Catalog       │
│ • Custom Apps   │    │ • Response Cache │    │ • Playlists     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Technical Stack
- **Language**: Go 1.21+
- **Apple Music API**: [go-apple-music](https://github.com/minchao/go-apple-music) library
- **Protocol**: Model Context Protocol (MCP)
- **Authentication**: Apple Music OAuth 2.0
- **Storage**: Local SQLite for caching and token storage
- **Configuration**: YAML-based configuration

## 🚀 Quick Start

### Prerequisites

1. **Apple Developer Account**: You'll need developer credentials to access the Apple Music API
2. **Go 1.21+**: [Install Go](https://golang.org/doc/install)
3. **Apple Music Subscription**: Required for full functionality

### Installation

```bash
# Clone the repository
git clone https://github.com/toby/apple-music-mcp.git
cd apple-music-mcp

# Build the server
go build -o apple-music-mcp ./cmd/server

# Or install directly
go install github.com/toby/apple-music-mcp/cmd/server@latest
```

### Apple Music API Setup

1. **Create an Apple Developer Account** at [developer.apple.com](https://developer.apple.com)

2. **Generate a MusicKit Private Key**:
   - Go to "Certificates, Identifiers & Profiles"
   - Under "Keys", create a new key with MusicKit enabled
   - Download the `.p8` private key file

3. **Get your credentials**:
   - **Team ID**: Found in your Apple Developer account
   - **Key ID**: The identifier for your MusicKit key
   - **Private Key**: The downloaded `.p8` file

### Configuration

Create a configuration file at `~/.config/apple-music-mcp/config.yaml`:

```yaml
apple_music:
  team_id: "YOUR_TEAM_ID"
  key_id: "YOUR_KEY_ID"
  private_key_path: "/path/to/your/AuthKey_KEYID.p8"
  
server:
  port: 8080
  log_level: "info"
  
storage:
  database_path: "~/.local/share/apple-music-mcp/data.db"
  
cache:
  ttl: "1h"
  max_size: 1000
```

### Running the Server

```bash
# Start the MCP server
./apple-music-mcp

# Or with custom config
./apple-music-mcp --config /path/to/config.yaml

# Enable debug logging
./apple-music-mcp --log-level debug
```

The server will start and listen for MCP connections, typically on a Unix socket or stdio.

## 🛠️ MCP Tools

The server exposes the following tools to AI assistants:

### Library Tools
- `search_music` - Search for tracks, albums, artists, playlists
- `get_library_albums` - Retrieve user's album collection
- `get_library_playlists` - Get user's playlists
- `get_recently_played` - Access listening history
- `get_recommendations` - Get personalized suggestions

### Playlist Tools
- `create_playlist` - Create a new playlist
- `add_to_playlist` - Add tracks to existing playlist
- `remove_from_playlist` - Remove tracks from playlist
- `update_playlist` - Modify playlist metadata

### Playback Tools
- `get_current_track` - Get currently playing song
- `get_queue` - View upcoming tracks
- `add_to_queue` - Add songs to play queue

### Discovery Tools
- `get_top_charts` - Browse popular music
- `get_new_releases` - Find latest albums
- `search_catalog` - Explore Apple Music catalog

## 💬 Usage Examples

### With Claude Desktop

Add this to your Claude Desktop MCP configuration:

```json
{
  "mcpServers": {
    "apple-music": {
      "command": "/path/to/apple-music-mcp",
      "args": ["--config", "/path/to/config.yaml"]
    }
  }
}
```

### Example Conversations

**User**: "Show me my recently added albums"
**AI**: Using the `get_library_albums` tool with recent filter...

**User**: "Create a workout playlist with high-energy pop songs from my library"
**AI**: I'll search your library for upbeat pop tracks and create a new playlist...

**User**: "What's similar to the last song I played?"
**AI**: Let me check your recently played tracks and find recommendations...

## 🔧 Development

### Project Structure

```
apple-music-mcp/
├── cmd/server/           # Main application entry point
├── internal/
│   ├── auth/            # Apple Music authentication
│   ├── api/             # Apple Music API client wrapper
│   ├── mcp/             # MCP protocol implementation
│   ├── tools/           # MCP tool definitions
│   ├── storage/         # Database and caching
│   └── config/          # Configuration management
├── pkg/                 # Public API packages
├── docs/                # Additional documentation
├── examples/            # Usage examples
└── tests/               # Test files
```

### Building from Source

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build for development
go build -o bin/apple-music-mcp ./cmd/server

# Build for production
go build -ldflags="-s -w" -o bin/apple-music-mcp ./cmd/server
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires API credentials)
go test -tags=integration ./...

# Test MCP protocol compliance
go test ./internal/mcp/...
```

### Adding New Tools

1. Define the tool in `internal/tools/`
2. Implement the handler logic
3. Register the tool in the MCP server
4. Add tests and documentation

Example tool structure:

```go
type SearchMusicTool struct {
    client *api.Client
}

func (t *SearchMusicTool) Execute(ctx context.Context, req mcp.ToolRequest) (*mcp.ToolResponse, error) {
    // Implementation here
}
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and add tests
4. **Run tests**: `go test ./...`
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Development Guidelines

- Follow Go best practices and idioms
- Write tests for new functionality
- Update documentation as needed
- Use conventional commits for commit messages
- Ensure backwards compatibility when possible

## 🐛 Troubleshooting

### Common Issues

**Authentication Errors**
```
Error: invalid_token
```
- Verify your Apple Developer credentials
- Check that your MusicKit key hasn't expired
- Ensure the private key file is readable

**Connection Issues**
```
Error: connection refused
```
- Verify the server is running
- Check the MCP client configuration
- Review firewall and network settings

**Permission Errors**
```
Error: insufficient_scope
```
- Your Apple Developer account may not have MusicKit enabled
- Check your Apple Music subscription status

### Debug Mode

Enable detailed logging:

```bash
./apple-music-mcp --log-level debug
```

Or set the environment variable:

```bash
export APPLE_MUSIC_MCP_LOG_LEVEL=debug
```

## 📚 Additional Resources

- [Model Context Protocol Documentation](https://github.com/modelcontextprotocol/docs)
- [Apple Music API Documentation](https://developer.apple.com/documentation/applemusicapi)
- [go-apple-music Library](https://github.com/minchao/go-apple-music)
- [MCP Server Examples](https://github.com/modelcontextprotocol/servers)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎵 Acknowledgments

- [MCP Team](https://github.com/modelcontextprotocol) for the excellent protocol
- [minchao](https://github.com/minchao) for the go-apple-music library
- Apple for the Music API
- The Go community for amazing tools and libraries

---

*Made with 🎶 for the AI-powered music experience*
