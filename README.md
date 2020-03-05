**Disclaimer 1 - features, command line options and playlist format is likelly going to change in the near future once more testing and feedback is received. This also means that I will not provide binaries for a while, since everything is a subject to change. Performance is also questionable, but it will be improved over time.**

# funstream

This application is used to create your very own M3U (HLS) playlist. Application requires a special *funstream playlist*, which contains information about other sources, such as separate M3U8 channels, media files or M3U playlists. It also allows you to customize defined external sources, such as overriding logo URL, changing or exluding both channel names and groups.

Features:
* Supports HLS (M3U8), streams (`Content-Type: application/octet-stream`), audio and Video files.
* No `FFmpeg` dependency - very lightweight and written in pure Go.
* Flexible *funstream playlists*. See bottom of this README.md

# Roadmap

Features that *might* never be implemented:
1. Define and edit EPG guides (in the same *funstream playlist*).
2. Docker image
3. Proper versioning

# Installation (Linux OSes with SystemD)

Make sure you read (#Documentation)[#Documentation] first!

1. Clone this repo and build binary:
```
$ git clone git@github.com:erkexzcx/funstream.git
$ cd funstream
$ go build -o funstream ./cmd/funstream/funstream.go
```

2. Move files into required paths:
```
# cp funstream /usr/bin/
# mkdir /etc/funstream
# cp <funstream_playlist.yml> /etc/funstream
# cp funstream.service /etc/systemd/system/
# systemctl daemon-reload
```

If you want to change port or any other command line argument - edit `/etc/systemd/system/funstream.service` line `ExecStart=`. Do not forget to re-run `systemctl daemon-reload` afterwards.

3. Start and enable SystemD service:
```
# systemctl start funstream.service
# systemctl enable funstream.service
```

# Upgrade

1. Update local copy of this repo and re-build binary:
```
$ git pull
$ go build -o funstream ./cmd/funstream/funstream.go
```

2. Replace existing binary with new binary:
```
# cp funstream /usr/bin/
```

If playlist format or systemd script has changed, then you might need to update these files too!

# Documentation

## Command line arguments

Execute binary. These command-line options are optional and used if you are not happy with default values:
* `-port 8989` - set custom web server's port. By default it uses `8989`.
* `-useragent "VLC/3.0.2.LibVLC/3.0.2"` - set custom user agent. By default it uses what VLC use (`VLC/3.0.2.LibVLC/3.0.2`).
* `-playlist "funstream_playlist.yaml"` - set location of your very personal funstream playlist. By default it uses `funstream_playlist.yaml` in current working directory.

## Playlist customization

See [funstream_playlist.example.yaml](https://github.com/erkexzcx/funstream/blob/master/funstream_playlist.example.yaml).

You don't need to explicitly define all fields. For example, this simple one-channel `yaml` file would perfectly work:
```
channels:
  - title: ExampleTV
    url: http://example.com/path/to/stream.m3u8
    logo: http://example.com/logos/exampletv.png
    group: Example TVs
```