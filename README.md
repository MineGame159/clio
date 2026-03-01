# Clio
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/MineGame159/clio?label=Code%20Size)
[![Lines of Code](https://img.shields.io/endpoint?url=https%3A%2F%2Ftokei.kojix2.net%2Fbadge%2Fgithub%2FMineGame159%2Fclio%2Flines)](https://tokei.kojix2.net/github/MineGame159/clio)

A CLI frontend for Streamio addons along with two custom, built-in addons for browsing your Real Debrid library and scraping torrents.

[![asciicast](https://asciinema.org/a/2BJ0OMX7rsLo2Qqw.svg)](https://asciinema.org/a/2BJ0OMX7rsLo2Qqw)

## Usage
To use Clio you need to create a config file in this location:
- Linux: `~/.config/clio.json`
- Windows: `%APPDATA%/clio.json`
- MacOS: `~/Library/Application Support/clio.json`

### Example
Replace `real-debrid-token` with your private API token for your Real Debrid account if you want to use the built-in addons:
```json
{
    "addons": [
        "<library:real-debrid-token>",
        "https://v3-cinemeta.strem.io/manifest.json",
        "<scraper:real-debrid-token>"
    ]
}
```