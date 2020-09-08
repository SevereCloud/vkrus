package vkrus_test

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"

	"github.com/SevereCloud/vkrus/v2"
	"github.com/sirupsen/logrus"
)

func TestVkHook_Fire(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"k0": 123,
		"k1": "abc",
		"k2": fmt.Errorf("%s", "error"),
	}

	f := func(hook *vkrus.VkHook, entry *logrus.Entry, wantErr bool) {
		t.Helper()

		err := hook.Fire(entry)
		if (err != nil) != wantErr {
			t.Errorf("VkHook.Fire() error = %v, wantErr %v", err, wantErr)
		}
	}

	peerID, err := strconv.Atoi(os.Getenv("PEER_ID"))
	if err != nil {
		t.Fatal(err)
	}

	token := os.Getenv("TOKEN")

	f(
		&vkrus.VkHook{
			PeerID:          peerID,
			VK:              api.NewVK(token),
			AppName:         "App name",
			Extra:           map[string]interface{}{"k3": "aoa"},
			Asynchronous:    true,
			DontParseLinks:  true,
			DisableMentions: true,
		},
		&logrus.Entry{
			Data:    data,
			Level:   logrus.ErrorLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		&vkrus.VkHook{Disabled: true},
		&logrus.Entry{},
		false,
	)
	f(
		vkrus.NewHook(peerID, "badtoken"),
		&logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "Test message",
		},
		true,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.FatalLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.ErrorLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.WarnLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.InfoLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.DebugLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.TraceLevel,
			Message: "Test message",
		},
		false,
	)
	f(
		vkrus.NewHook(peerID, token),
		&logrus.Entry{
			Level:   logrus.InfoLevel,
			Message: strings.Repeat("ж", 4097),
		},
		false,
	)
}

func TestVkHook_Levels(t *testing.T) {
	t.Parallel()

	f := func(hook *vkrus.VkHook, wantLevel []logrus.Level) {
		t.Helper()

		level := hook.Levels()
		if !reflect.DeepEqual(level, wantLevel) {
			t.Errorf("VkHook.Fire() level = %v, wantLevel %v", level, wantLevel)
		}
	}

	f(vkrus.NewHook(0, ""), vkrus.DefaultLevels)
	f(&vkrus.VkHook{}, logrus.AllLevels)
	f(&vkrus.VkHook{UseLevels: []logrus.Level{}}, []logrus.Level{})
}
