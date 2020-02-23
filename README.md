<h1 align="center">Joe Bot - Vk Adapter</h1>
<p align="center">
	<a href="https://godoc.org/github.com/tdakkota/vk"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?color=blue"></a>
</p>

This repository contains a module for the [Joe Bot library][joe]. Built using 
[vksdk][vksdk].

## Getting Started

This library is packaged using [Go modules][go-modules]. You can get it via:

```
go get github.com/tdakkota/joe-vk-adapter
```

### Example usage

In order to connect your bot to VK you can simply pass it as module when
creating a new bot:

```go
package main

import (
	"github.com/go-joe/joe"
	"github.com/tdakkota/joe-vk-adapter"
	"os"
)

func main() {
	b := joe.New("example-bot",
		vk.Adapter(os.Getenv("BOT_TOKEN")),
	…
	)

	b.Respond("ping", func(msg joe.Message) error {
		msg.Respond("pong")
		return nil
	})

	err := b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}
```

This adapter will emit the following events to the robot brain:

- `joe.ReceiveMessageEvent`
- `ChatCreateEvent`
- `ChatTitleUpdateEvent`
- `ChatPhotoUpdateEvent`
- `ChatPinUpdateEvent`
- `UserEnteredChatEvent`
- `UserLeavedChatEvent`

## License

[BSD-3-Clause](LICENSE)

[joe]: https://github.com/go-joe/joe
[vksdk]: https://github.com/SevereCloud/vksdk
[go-modules]: https://github.com/golang/go/wiki/Modules