# PaceRank Client

The PaceRank client observes processes running on the operating system and collect anonymous telemetry data and sends it
to the PaceRank digest service. The client can be compiled into an CLI command, or compile into a GUI app.

When the client runs it will recognize if you are typing into an editor, and it will listen to file changes in specified
folders. On file change it will analyze which syntax was used in the file, and save the metric.

The client **won't** save or send any data that is personal, or belong to the source code that has been analyzed.

## Run CLI command
```
go run ./cmd/cli -f /path/to/watch -f /path/to/watch2
```

First time you run this command it will open a browser to authorize the client to send statistics to your account in
PaceRank.

## Compile Windows
First you need to install josephspurrier/goversioninfo and make the command available in path:
```
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
```

Then we can compile

```
go generate
go build -ldflags -H=windowsgui -o pacerank.exe
```

The GUI client has a dependency on sciter.dll and the DLL file need to be present in the same folder as the executable.
