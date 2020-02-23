// Package vk implements a VK adapter for the joe bot library.
package vk

import (
	"context"
	"strconv"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/longpoll-bot"
	"github.com/SevereCloud/vksdk/object"
	"github.com/go-joe/joe"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// BotAdapter implements a joe.Adapter that reads and writes messages to and
// from VK.
type BotAdapter struct {
	vk      *api.VK
	lp      longpoll.Longpoll
	context context.Context
	logger  *zap.Logger
}

// Config contains the configuration of a BotAdapter.
type Config struct {
	Token  string
	Logger *zap.Logger
}

var ErrGetBotInfo = errors.New("failed to get bot info")

// NewAdapter creates a new *BotAdapter that connects to VK. Note that you
// will usually configure the VK adapter as joe.Module (i.e. using the
// Adapter function of this package).
func NewAdapter(ctx context.Context, conf Config) (*BotAdapter, error) {
	vk := api.Init(conf.Token)

	logger := conf.Logger
	if logger == nil {
		logger = zap.NewNop()
	}

	r, err := vk.GroupsGetByID(api.Params{})
	if err != nil {
		return nil, errors.Wrap(err, ErrGetBotInfo.Error())
	}

	if len(r) < 1 {
		return nil, ErrGetBotInfo
	}

	group := r[0]
	logger.Info("got group info",
		zap.Int("group_id", group.ID),
		zap.String("name", group.Name),
	)

	lp, err := longpoll.Init(vk, group.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init longpolling")
	}

	b := &BotAdapter{
		vk:      vk,
		lp:      lp,
		context: ctx,
		logger:  logger,
	}

	return b, nil
}

// Send implements joe.Adapter by sending all received text messages to the
// given chat.
func (b *BotAdapter) Send(text, chat string) error {
	peerID, err := strconv.Atoi(chat)
	if err != nil {
		return b.sendDomain(text, chat)
	}

	return b.sendPeerID(text, peerID)
}

func (b *BotAdapter) sendDomain(text, domain string) error {
	b.logger.Info("Sending message to chat",
		zap.String("domain", domain),
	)

	_, err := b.vk.MessagesSend(api.Params{
		"domain":    domain,
		"message":   text,
		"random_id": 0,
	})

	return err
}

func (b *BotAdapter) sendPeerID(text string, peerID int) error {
	b.logger.Info("Sending message to chat",
		zap.Int("peer_id", peerID),
	)

	_, err := b.vk.MessagesSend(api.Params{
		"peer_id":   peerID,
		"message":   text,
		"random_id": 0,
	})

	return err
}

// RegisterAt implements the joe.Adapter interface by emitting the vk API
// events to the given brain.
func (b *BotAdapter) RegisterAt(brain *joe.Brain) {
	go func() {
		err := b.lp.Run()
		if err != nil {
			b.logger.Error("longpoll run failed", zap.Error(err))
		}
	}()

	b.lp.MessageNew(func(object object.MessageNewObject, i int) {
		b.logger.Info("Message received",
			zap.Int("peer_id", object.Message.PeerID),
		)
		brain.Emit(b.dispatch(object))
	})
}

func (b *BotAdapter) dispatch(object object.MessageNewObject) interface{} {
	switch object.Message.Action.Type {
	case "chat_create":
		return ChatCreateEvent{
			Channel: strconv.Itoa(object.Message.PeerID),
			Text:    object.Message.Action.Text,
			Data:    object,
		}

	case "chat_title_update":
		return ChatTitleUpdateEvent{
			Channel: strconv.Itoa(object.Message.PeerID),
			NewText: object.Message.Action.Text,
			Data:    object,
		}

	case "chat_photo_update", "chat_photo_remove":
		return ChatPhotoUpdateEvent{
			Channel:  strconv.Itoa(object.Message.PeerID),
			NewPhoto: object.Message.Action.Photo,
			Removed:  object.Message.Action.Type == "chat_photo_remove",
			Data:     object,
		}

	case "chat_pin_message", "chat_unpin_message":
		return ChatPinUpdateEvent{
			Channel:   strconv.Itoa(object.Message.PeerID),
			UserID:    strconv.Itoa(object.Message.Action.MemberID),
			MessageID: strconv.Itoa(object.Message.Action.ConversationMessageID),
			Unpinned:  object.Message.Action.Type == "chat_unpin_message",
			Data:      object,
		}

	case "chat_invite_user", "chat_invite_user_by_link":
		return UserEnteredChatEvent{
			Channel: strconv.Itoa(object.Message.PeerID),
			UserID:  strconv.Itoa(object.Message.Action.MemberID),
			ByLink:  object.Message.Action.Type == "chat_invite_user_by_link",
			Data:    object,
		}

	case "chat_kick_user":
		return UserLeavedChatEvent{
			Channel: strconv.Itoa(object.Message.PeerID),
			UserID:  strconv.Itoa(object.Message.Action.MemberID),
			Data:    object,
		}

	default:
		return joe.ReceiveMessageEvent{
			Text:     object.Message.Text,
			ID:       strconv.Itoa(object.Message.ID),
			AuthorID: strconv.Itoa(object.Message.FromID),
			Channel:  strconv.Itoa(object.Message.PeerID),
			Data:     object,
		}
	}
}

// Close disconnects the adapter from the vk API.
func (b *BotAdapter) Close() error {
	b.lp.Shutdown()
	return nil
}
