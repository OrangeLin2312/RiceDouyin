package service

type Video struct {
	id             int64  // 视频唯一标识
	author         User   // 视频作者信息
	play_url       string // 视频播放地址
	cover_url      string // 视频封面地址
	favorite_count int64  // 视频的点赞总数
	comment_count  int64  // 视频的评论总数
	is_favorite    bool   // true-已点赞，false-未点赞
	title          string // 视频标题
}

type User struct {
	id               int64  // 用户id
	name             string // 用户名称
	follow_count     int64  // 关注总数
	follower_count   int64  // 粉丝总数
	is_follow        bool   // true-已关注，false-未关注
	avatar           string //用户头像
	background_image string //用户个人页顶部大图
	signature        string //个人简介
	total_favorited  int64  //获赞数量
	work_count       int64  //作品数量
	favorite_count   int64  //点赞数量
}
type favoriteService interface {
	//主业务
	//点赞行为,传入登录用户的id和被点赞视频的id以及行为状态
	FavoriteAction(uid int64, vid int64, actiontype int32) error
	//获取喜欢列表，传入当前用户（非登录用户）的id，返回该用户的喜欢列表
	FavoriteList(uid int64) ([]Video, error)
	//其他业务
	//IsFavorite(uid int64, vid int64) (bool, error)
}
