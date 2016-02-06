# Steam DOTA2 Source 2 Replay Parser Webserver

[![Build Status](https://travis-ci.org/HighGroundVision/Mango.svg?branch=master)](https://travis-ci.org/HighGroundVision/Mango)

Uses [dotabuff/manta](https://github.com/dotabuff/manta) to prase replay files for Source 2 (Dota 2 Reborn) engine. (currently dose not support source 1 replays)

**Project Status:**

- Main web server is online. Yippy!
- Need to pipe in URL parameter to get replay files from.
- Validate current JSON structure for best load / parse performance.

## Getting Started

Mango will return JSON array of events. Each event is timestamped with the game time and each event has a type. You will need to parse each event or event of specific types and decide how you're going to use it.

## License

See the LICENSE file.

## Help

If you have any questions, you can tweet us at [@DotaHGV](https://twitter.com/DotaHGV)

## Authors and Acknowledgements

Mango is maintained and development by [HGV](http://www.highgroundvision.com), a leading Dota 2 data visualization and analysis web site. HGV's Team:

* [Jamie Webster](https://github.com/RGBKnights) 
* [Graham Clifford](https://github.com/gclifford)

Special thanks to the following people:

* [Jason Coene](https://github.com/jcoene) built Dotabuff's Source 2 parser [manta](https://github.com/dotabuff/manta).
* [Robin Dietrich](https://github.com/invokr) built the C++ parser [Alice](https://github.com/AliceStats/Alice).
* [Howard Chung](https://github.com/howardchung) built yasp parser [yasp](https://github.com/yasp-dota/yasp).
