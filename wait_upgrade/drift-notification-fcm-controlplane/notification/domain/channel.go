package domain

type Channel string

const (
	ChannelFCM   Channel = "fcm"
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelInApp Channel = "in_app"
)

func (c Channel) Valid() bool {
	switch c {
	case ChannelFCM, ChannelEmail, ChannelSMS, ChannelInApp:
		return true
	default:
		return false
	}
}
