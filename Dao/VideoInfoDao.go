package Dao

import (
	"time"
)

type TableVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"` //视频名，5.23添加
}

type Video struct {
	TableVideo
	Author        User  `json:"author"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount  int64 `json:"comment_count"`
	IsFavorite    bool  `json:"is_favorite"`
}
