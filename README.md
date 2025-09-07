# Synthera CLI

**Command line interface for interacting with Synthera backend**, build with Go for fast and efficient operations (even though its overkill).

## Purpose

Synthera CLI serves as the user interface to interact with the Synthera backend, handling requests and responses efficiently.

## Installation & Setup

### Prerequisites

- Go (version compatible with project)
- Git
- Termux (if you're using Android)

### Installation method

#### Build from Source

1. Close the repository:

```bash
git clone https://github.com/slrmyapi/synthera-cli.git
```

2. Navigate into the repo:

```bash
cd synthera-cli
```

3. Ensure dependencies are installed:

```bash
go mod tidy
```

4. Build the executable:

```bash
go build
```

5. Run the CLI:

```bash
./synthera
```

#### Using Precompiled binary

##### Linux / Mac

```bash
wget https://github.com/slrmyapi/synthera-cli/releases/latest/download/synthera-$(uname -s)-$(uname -m) -O synthera
chmod +x synthera
./synthera
```

##### Windows

```powershell
wget https://github.com/slrmyapi/synthera-cli/releases/latest/download/synthera-windows-x86_64.exe -OutFile synthera.exe
.\synthera.exe
```

##### Android (Termux)

```bash
apt update && apt upgrade -y
apt install wget -y
wget https://github.com/slrmyapi/synthera-cli/releases/latest/download/synthera-android-$(uname -m) -O synthera
chmod +x synthera
./synthera
```

> **Note:** Always run `./synthera` (or `.\synthera.exe` on Windows) for consistency

## Usage

1. Run the CLI:

```
./synthera
```

2. When prompted for token, enter your API token (can be found on website → Account → Profile → Token)
3. After token verification, you can start using the CLI to interact with the backend.

### Flow:

```
User enters command → CLI sends request → Backend processes → Response returned
```

## Features

- Basic CLI interface to interact with Synthera backend
- Token based authentication
- Sends requests to backend and receives responses
- Build with Go for fast and efficient execution
- Easy to extend with future commands / features

## Architecture & Structure

```bash
cli
├── api # Handles backend API requests
│   ├── client.go
│   └── types.go
├── build.sh
├── go.mod
├── main.go # Entry point
├── go.sum
├── .token.json # Your API token
├── README.md
├── ui # CLI interface components
│   ├── items.go
│   ├── models.go
│   ├── types.go
│   └── views.go
└── utils # Utility functions
    └── token.go
```

### Notes:

- Token is stored in `.token.json`; if invalid, user must delete it to re-enter new token. For example:

```bash
rm .token.json
./synthera
```

- Future improvements: double-check API response, add more features.

## Contributing

- Contributions are welcome.
- Please keep the code clean, modular and maintainable.
- Submit PRs or issues for review.

## License

No license assigned. Use at your own risk. Do not redistribute or use commercially without permission.

## Known Issues

- Token validation: if the token invalid, user must manually delete `.token.json` to re-enter a new token.
- Future improvements:
    - Add token double-check / validation to prevent manual deletion
    - Expand CLI commands to match website functionality
    - Improve error handling and user feedback
