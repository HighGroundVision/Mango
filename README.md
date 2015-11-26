# Steam DOTA2 Source 2 Replay Parser Webserver

**Build Status Placeholder**

Uses [dotabuff/manta](https://github.com/dotabuff/manta) to prase replay files for Source 2 (Dota 2 Reborn) engine. (currently dose not support source 1 replays)

**Project Status:**

- Main web server is online. Yippy!
- Need to pipe in url parameter to get replay files from.
- Validate Json Structure for best load / prase conditions.

## Getting Started

Manta will return json array of events. Each event is timestamped with the game time and each event has a type. You will need to prase each event or event of speific types and decide how you're going to use it.
