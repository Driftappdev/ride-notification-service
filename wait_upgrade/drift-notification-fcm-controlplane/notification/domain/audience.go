package domain

import "errors"

type AudienceType string

const (
	AudienceToken     AudienceType = "token"
	AudienceTokens    AudienceType = "tokens"
	AudienceTopic     AudienceType = "topic"
	AudienceCondition AudienceType = "condition"
	AudienceUser      AudienceType = "user"
	AudienceSegment   AudienceType = "segment"
)

type Audience struct {
	Type       AudienceType
	Token      string
	Tokens     []string
	Topic      string
	Condition  string
	UserID     string
	SegmentID  string
}

func (a Audience) Validate() error {
	switch a.Type {
	case AudienceToken:
		if a.Token == "" { return errors.New("audience token is required") }
	case AudienceTokens:
		if len(a.Tokens) == 0 { return errors.New("audience tokens are required") }
	case AudienceTopic:
		if a.Topic == "" { return errors.New("audience topic is required") }
	case AudienceCondition:
		if a.Condition == "" { return errors.New("audience condition is required") }
	case AudienceUser:
		if a.UserID == "" { return errors.New("audience user_id is required") }
	case AudienceSegment:
		if a.SegmentID == "" { return errors.New("audience segment_id is required") }
	default:
		return errors.New("invalid audience type")
	}
	return nil
}
