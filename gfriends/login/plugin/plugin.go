package plugin

type AuthRsp struct {
	ErrCode           int    `json:"errcode"`
	ErrMsg            string `json:"errmsg"`
	StrOpenid         string `json:"openid"`
	StrUnionid        string `json:"unionid"`
	StrRefreshToken   string `json:"refresh_token"`
	StrAccessToken    string `json:"access_token"`
	AccessTokenExpire int64  `json:"expires_in"`
}

type RefreshRsp struct {
	ErrCode           int    `json:"errcode"`
	ErrMsg            string `json:"errmsg"`
	StrOpenid         string `json:"openid"`
	StrUnionid        string `json:"unionid"`
	StrAccessToken    string `json:"access_token"`
	AccessTokenExpire int64  `json:"expires_in"`
}

type LoginHandler interface {
	AuthorziationCode(strCode string) (*AuthRsp, error)
	RefreshAccessToken(refreshToken string) (*RefreshRsp, error)
}
