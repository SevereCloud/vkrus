// Package vkrus VK Hook for Logrus
package vkrus // import "github.com/SevereCloud/vkrus"

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/SevereCloud/vksdk/api"
	"github.com/sirupsen/logrus"
)

const maxLenMessage = 4096

// DefaultLevels to be fired when logging on
var DefaultLevels = []logrus.Level{ // nolint: gochecknoglobals
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

// VkHook is a logrus Hook for dispatching messages to the specified
// channel on Slack.
//
// PeerID - destination ID.
// For user: 'User ID', e.g. '12345'.
// For chat: '2000000000' + 'chat_id', e.g. '2000000001'
type VkHook struct {
	PeerID int // Destination ID
	VK     *api.VK

	AppName         string
	UseLevels       []logrus.Level
	Extra           map[string]interface{}
	Asynchronous    bool
	Disabled        bool
	DontParseLinks  bool // true - links will not attach snippet
	DisableMentions bool // true - mention of user will not generate notification for him
}

// NewHook return VK Hook for Logrus
func NewHook(peerID int, token string) *VkHook {
	hook := &VkHook{
		PeerID:    peerID,
		VK:        api.NewVK(token),
		UseLevels: DefaultLevels,
	}

	return hook
}

// Levels sets which levels to sent to VK
func (hook *VkHook) Levels() []logrus.Level {
	if hook.UseLevels == nil {
		return logrus.AllLevels
	}

	return hook.UseLevels
}

// Fire - Sent event to VK
func (hook *VkHook) Fire(entry *logrus.Entry) error {
	if hook.Disabled {
		return nil
	}

	newEntry := hook.newEntry(entry)
	msg := hook.createMessage(newEntry)

	return hook.sendMessage(msg)
}

func (hook *VkHook) newEntry(entry *logrus.Entry) *logrus.Entry {
	data := map[string]interface{}{}

	for k, v := range hook.Extra {
		data[k] = v
	}

	for k, v := range entry.Data {
		data[k] = v
	}

	newEntry := &logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
	}

	return newEntry
}

func (hook *VkHook) createMessage(entry *logrus.Entry) string {
	var msg string

	if hook.AppName != "" {
		msg += hook.AppName + "\n"
	}

	nameLevel := strconv.Itoa(int(entry.Level))

	nameLevelByte, err := entry.Level.MarshalText()
	if err == nil {
		nameLevel = strings.ToUpper(string(nameLevelByte))
	}

	msg += nameLevel + ": " + entry.Message

	if len(entry.Data) > 0 {
		msg += "\n\nMessage fields:\n"
		for k, v := range entry.Data {
			msg += fmt.Sprintf("%s=%+v\n", k, v)
		}
	}

	return msg
}

func (hook *VkHook) sendMessage(msg string) error {
	offset := 0
	now := 0

	count := len(msg)
	if len(msg) > maxLenMessage {
		count = maxLenMessage
	}

	for now < len(msg) {
		var text string

		for now < offset+count {
			runeValue, width := utf8.DecodeRuneInString(msg[now:])
			if now+width <= offset+count {
				text += string(runeValue)
				now += width
			}
		}

		params := api.Params{
			"peer_id":          hook.PeerID,
			"message":          text,
			"random_id":        0,
			"dont_parse_links": hook.DontParseLinks,
			"disable_mentions": hook.DisableMentions,
		}

		if hook.Asynchronous {
			go func() {
				_, _ = hook.VK.MessagesSend(params)
			}()
		} else {
			_, err := hook.VK.MessagesSend(params)
			if err != nil {
				return err
			}
		}

		offset = now
		if len(msg) >= offset+maxLenMessage {
			count = maxLenMessage
		} else {
			count = len(msg) - now
		}
	}

	return nil
}
