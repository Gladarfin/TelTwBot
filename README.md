# TelTwBot

/twitch-bot
├── Cmd/
│   └── Bot/
│       └── main.go      # Minimal main, just wires things together
├── Internal/
│   ├── Bot/             # Core bot logic
│   ├── Commands/        # Command handlers
│   ├── Config/          # Configuration loading
│   └── Twitch/          # Twitch-specific code
├── go.mod
├── go.sum
└── Data/                # Data files         