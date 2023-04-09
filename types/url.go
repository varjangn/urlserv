package types

type URL struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	ShortId   string `json:"short_id"`
	Long      string `json:"long_url"`
	CreatedAt int64  `json:"created_ts"`
}

func NewURL(user *User, shortId, long string) *URL {
	return &URL{
		Id:        0,
		UserId:    user.Id,
		ShortId:   shortId,
		Long:      long,
		CreatedAt: 0,
	}
}
