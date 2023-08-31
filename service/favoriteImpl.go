package service

import (
	"encoding/json"
	"github.com/luuuweiii/RiceDouyin/config"
	"github.com/luuuweiii/RiceDouyin/dao"
	"github.com/luuuweiii/RiceDouyin/utils/rabbitMQ"
	"log"
	"sync"
)

type FavoriteService struct {
	Rmq         *rabbitMQ.FavoriteMq
	FavoriteDao *dao.FavoriteDao
	UserDao     *dao.UserDao
	VideoDao    *dao.VideoDao
	//UserService  *UserService
	VideoService *VideoService
}

func NewFavoriteService(
	Rmq *rabbitMQ.FavoriteMq,
	FavoriteDao *dao.FavoriteDao,
	UserDao *dao.UserDao,
	VideoDao *dao.VideoDao,
	//UserService *UserService,
	VideoService *VideoService) *FavoriteService {
	return &FavoriteService{
		Rmq:         Rmq,
		FavoriteDao: FavoriteDao,
		UserDao:     UserDao,
		VideoDao:    VideoDao,
		//UserService:  UserService,
		VideoService: VideoService,
	}
}

// FavoriteAction:点赞行为,传入登录用户的id和被点赞视频的id
// 1先更新自己favorite表
//
//	//   a.判断点赞操作actiontype
//	//   b.点赞则添加（需不需要检查一下表中原本有没有）flag=1，取消赞则更新actiontype,flag=-1
//	//2更新视频Video获赞数
//	//   通过vid在Video里对视频获赞数+flag
//	//3更新用户User点赞数
//	//   通过uid对点赞数+flag
//	//4更新用户User获赞数
//	//   a.通过vid在Video获取视频发布者v_uid
//	//   b.通过v_uid在User对获赞数+flag
func (favorite *FavoriteService) FavoriteAction(uid int64, vid int64, actiontype int32) error {
	var err error

	msg := config.RmqMessage{
		UserId:     uid,
		VideoId:    vid,
		ActionType: actiontype,
	}
	msgbody, err1 := json.Marshal(msg)
	//rabbitMQ.RmqFavoriteTrue.FavoritePublish(msgbody)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	favorite.Rmq.FavoritePublish(msgbody)
	return err
}

// FavoriteList：喜欢列表，传入当前用户（非登录用户）的id，返回该用户的喜欢列表
func (favorite *FavoriteService) FavoriteList(uid int64) ([]Video, error) {
	//1先从数据库找到对应uid的点赞视频videoIdList
	videoIdList, err := favorite.FavoriteDao.GetVideoIdList(uid)
	if err != nil {
		log.Println("GetVideoIdList have a error:", err)
		return nil, err
	}

	//2用videoIdList去并发操作getVideo()
	var wg sync.WaitGroup
	videoList := make([]Video, len(videoIdList))
	for index, vid := range videoIdList {
		wg.Add(1)
		go favorite.GetVideoList(index, vid, videoList, &wg)
	}

	wg.Wait()

	return videoList, nil
}

// 获取点赞状态
func (favorite *FavoriteService) IsFavorite(uid int64, vid int64) (bool, error) {
	var (
		exsit bool
		err   error
	)
	exsit, err = favorite.FavoriteDao.QueryIsFavorite(uid, vid)
	if err != nil {
		log.Println("have a QueryIsFavorite error")
		return exsit, err
	}
	return exsit, nil
}

// GetVideoList获取视频列表
func (f *FavoriteService) GetVideoList(index int, vid int64, list []Video, wg *sync.WaitGroup) {
	defer wg.Done()
	//var mu sync.Mutex
	video, err := f.VideoService.GetVideo(vid)
	if err != nil {
		return
	}
	//mu.Lock()
	list[index] = video
	//mu.Unlock()
}
