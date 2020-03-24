# vkrus

[![Build Status](https://travis-ci.com/SevereCloud/vkrus.svg?branch=master)](https://travis-ci.com/SevereCloud/vkrus)
[![Documentation](https://godoc.org/github.com/SevereCloud/vkrus?status.svg)](https://pkg.go.dev/github.com/SevereCloud/vkrus)
[![codecov](https://codecov.io/gh/SevereCloud/vkrus/branch/master/graph/badge.svg)](https://codecov.io/gh/SevereCloud/vkrus)
[![VK chat](https://img.shields.io/badge/VK%20chat-%234a76a8.svg?logo=VK&logoColor=white)](https://vk.me/join/AJQ1d6Or8Q00Y_CSOESfbqGt)
[![release](https://img.shields.io/github/v/tag/SevereCloud/vkrus?label=release)](https://github.com/SevereCloud/vkrus/releases)
[![license](https://img.shields.io/github/license/SevereCloud/vkrus.svg?maxAge=2592000)](https://github.com/SevereCloud/vkrus/blob/master/LICENSE)

[Logrus](https://github.com/sirupsen/logrus) hook for [VK](https://vk.com) using [VKSDK](https://github.com/SevereCloud/vksdk).

### Usage

This library is packaged using [Go modules](https://github.com/golang/go/wiki/Modules). You can get it via:

```sh
# go mod init mymodulename
go get github.com/SevereCloud/vkrus@latest
```


The hook must be configured with:
- A peer ID (
  - For user: 'User ID', e.g. `12345`
  - For chat: '2000000000' + 'chat_id', e.g. `2000000001` (**chat_id for community**)
- [Access Token](https://vk.com/dev/access_token) with **messages** rights (found in your community settings)

```go
package main

import (
	vkrus "github.com/SevereCloud/vkrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	peerID := 117253521   // USE strconv.Atoi(os.Getenv("PEER_ID"))
	groupToken := "token" // USE os.Getenv("TOKEN")
	hook := vkrus.NewHook(peerID, groupToken)
	hook.UseLevels = log.AllLevels

	log.AddHook(hook)
}

func main() {
	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
}

```

DefaultLevels to be fired when logging on
```go
var DefaultLevels = []logrus.Level{
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}
```

Optional params:

- AppName
- UseLevels
- Extra
- Asynchronous
- Disabled
- DontParseLinks
- DisableMentions

### License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FSevereCloud%2Fvkrus.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FSevereCloud%2Fvkrus?ref=badge_large)