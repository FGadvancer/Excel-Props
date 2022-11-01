package constant

import "errors"

const (
	ArgsErr          = 100
	NotUserInfo      = 101
	PasswordError    = 102
	ParseTokenFailed = 103
)

type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}

var (
	ErrTokenExpired             = ErrInfo{701, TokenExpiredMsg.Error()}
	ErrTokenInvalid             = ErrInfo{702, TokenInvalidMsg.Error()}
	ErrTokenMalformed           = ErrInfo{703, TokenMalformedMsg.Error()}
	ErrTokenNotValidYet         = ErrInfo{704, TokenNotValidYetMsg.Error()}
	ErrTokenUnknown             = ErrInfo{705, TokenUnknownMsg.Error()}
	ErrTokenKicked              = ErrInfo{706, TokenUserKickedMsg.Error()}
	ErrTokenDifferentPlatformID = ErrInfo{707, TokenDifferentPlatformIDMsg.Error()}
	ErrTokenDifferentUserID     = ErrInfo{708, TokenDifferentUserIDMsg.Error()}
)
var (
	ParseTokenMsg               = errors.New("parse token failed")
	TokenExpiredMsg             = errors.New("token is timed out, please log in again")
	TokenInvalidMsg             = errors.New("token has been invalidated")
	TokenNotValidYetMsg         = errors.New("token not active yet")
	TokenMalformedMsg           = errors.New("that's not even a token")
	TokenUnknownMsg             = errors.New("couldn't handle this token")
	TokenUserKickedMsg          = errors.New("user has been kicked")
	TokenDifferentPlatformIDMsg = errors.New("different platformID")
	TokenDifferentUserIDMsg     = errors.New("different userID")
	AccessMsg                   = errors.New("no permission")
	StatusMsg                   = errors.New("status is abnormal")
	DBMsg                       = errors.New("db failed")
	ArgsMsg                     = errors.New("args failed")
	CallBackMsg                 = errors.New("callback failed")
	InvitationMsg               = errors.New("invitationCode error")

	ThirdPartyMsg = errors.New("third party error")
)
