package vk

import (
	"github.com/SevereCloud/vksdk/object"
)

// ChatCreateEvent struct
type ChatCreateEvent struct {
	Channel string // The id of created channel.
	Text    string // Chat title

	// original response (objects.MessageNewObject)
	Data interface{}
}

// ChatTitleUpdateEvent struct
type ChatTitleUpdateEvent struct {
	Channel string // The id of created channel.
	NewText string // New chat title

	// original response (objects.MessageNewObject)
	Data interface{}
}

// ChatPhotoUpdateEvent struct
type ChatPhotoUpdateEvent struct {
	Channel  string                            // The id of created channel.
	NewPhoto object.MessagesMessageActionPhoto // The object with new cover photo urls
	Removed  bool                              // true if cover photo has changed, otherwise false

	// original response (objects.MessageNewObject)
	Data interface{}
}

// ChatPinUpdateEvent struct
type ChatPinUpdateEvent struct {
	Channel   string // The id of created channel.
	UserID    string // A string identifying the user who changed pin
	MessageID string // A string identifying the changed pin message
	Unpinned  bool   // true if message has unpinned, otherwise false

	// original response (objects.MessageNewObject)
	Data interface{}
}

// UserEnteredChatEvent struct
type UserEnteredChatEvent struct {
	Channel string // The channel over which the message was received.
	UserID  string // A string identifying the new user in chat
	ByLink  bool

	// original response (objects.MessageNewObject)
	Data interface{}
}

// UserLeavedChatEvent struct
type UserLeavedChatEvent struct {
	Channel string // The channel over which the message was received.
	UserID  string // A string identifying the leaved user in chat

	// original response (objects.MessageNewObject)
	Data interface{}
}
