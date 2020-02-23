<h1 align="center">Joe Bot - VK Adapter</h1>

<p align="center">Connecting joe with the VK chat application. https://github.com/go-joe/joe</p>
<p align="center">
	<a href="https://goreportcard.com/report/github.com/tdakkota/joe-vk-adapter"><img src="https://goreportcard.com/badge/github.com/tdakkota/joe-vk-adapter"></a>
	<a href="https://www.codefactor.io/repository/github/tdakkota/joe-vk-adapter"><img src="https://www.codefactor.io/repository/github/tdakkota/joe-vk-adapter/badge"></a>
	<a href="https://godoc.org/github.com/tdakkota/joe-vk-adapter"><img src="https://godoc.org/github.com/tdakkota/joe-vk-adapter?status.svg"></a>
	<a href="https://opensource.org/licenses/BSD-3-Clause"><img src="https://img.shields.io/badge/License-BSD%203--Clause-blue.svg"></a>
</p>

---

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
	"os"

	"github.com/go-joe/joe"
	"github.com/tdakkota/joe-vk-adapter"
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