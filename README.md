This is a simple initial Go project for a Twitch and Telegram bot that is not yet finished. The purpose is to explore some functions of APIs.

# TelTwBot
```
/twitch-bot
├── Cmd/
│   └── Bot/
│       └── main.go      # Minimal main, just wires things together
├── Internal/
│   ├── Config/          # Configuration loading
│   ├── Database/        # Core database logic and scripts
│   ├── Interfaces/      # Core interfaces
│   ├── Telegram/        # Core telegram bot logic
│   ├── TwitchBot/       # Core twitch bot logic
│   ├── Commands/        # Command handlers
│   └── Twitch/          # Twitch-specific code
├── go.mod
├── go.sum
└── Data/                # Data files
```       

#### Twitch Commands
```
!help - displays a list of available commands;
!hello - displays a random greeting to user;
!title - displays the current stream title;
!game - shows what game is currently being played;
!who - shows participating streamers;		
!role - shows the user role on current channel;
!stats - shows user stats;
!duel - starts the duel with other user;
!up - increase selected stat if there is enough free points.
!hl - shows game completion times from HowLongToBeat.com.
```

#### Telegram Commands
```
uptime - get stream uptime;
test - just for test;
math - do simple math (e.g. a + b, a*b, a/b, a-b);
help - show help;
```


