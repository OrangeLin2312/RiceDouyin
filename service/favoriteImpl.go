package service

import (
	"fmt"
	"github.com/luuuweiii/RiceDouyin/dao"
	"gorm.io/gorm"
	"log"
	"sync"
)

type FavoriteService struct {
	FavoriteDao *dao.FavoriteDao
	UserDao     *dao.UserDao
	VideoDao    *dao.VideoDao
	//UserService  *UserService
	VideoService *VideoService
}

func NewFavoriteService(
	FavoriteDao *dao.FavoriteDao,
	UserDao *dao.UserDao,
	VideoDao *dao.VideoDao,
	//UserService *UserService,
	VideoService *VideoService) *FavoriteService {
	return &FavoriteService{
		FavoriteDao: FavoriteDao,
		UserDao:     UserDao,
		VideoDao:    VideoDao,
		//UserService:  UserService,
		VideoService: VideoService,
	}
}

// FavoriteAction:点赞行为,传入登录用户的id和被点赞视频的id
func (favorite *FavoriteService) FavoriteAction(uid int64, vid int64, actiontype int32) error {
	var flag int32
	var db *gorm.DB
	var err error
	//1先更新自己favorite表
	//   a.判断点赞操作actiontype
	//   b.点赞则添加（需不需要检查一下表中原本有没有）flag=1，取消赞则更新actiontype,flag=-1
	//2更新视频Video获赞数
	//   通过vid在Video里对视频获赞数+flag
	//3更新用户User点赞数
	//   通过uid对点赞数+flag
	//4更新用户User获赞数
	//   a.通过vid在Video获取视频发布者v_uid
	//   b.通过v_uid在User对获赞数+flag
	db = favorite.FavoriteDao.DB
	tx := db.Begin() // 开启事务
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()
	if actiontype == 1 {
		flag = 1
		err = favorite.ToFavoriteAction(tx, uid, vid, flag)
		if err != nil {
			log.Println("FavoriteAciton failed")
		}
	} else {
		flag = -1
		err = favorite.ToFavoriteAction(tx, uid, vid, flag)
		if err != nil {
			log.Println("FavoriteAciton failed")
		}
	}
	tx.Commit()
	return err
}

// 连接事务对象和事务操作
func (f *FavoriteService) ToFavoriteAction(tx *gorm.DB, uid int64, vid int64, flag int32) error {
	var err error = nil
	err = f.FavoriteWithTransaction(uid, vid, tx, flag)
	if err != nil {
		log.Println("FavoriteWithTransaction failed")
		tx.Rollback()
		return err
	}
	err = f.VideolikedWithTransaction(vid, tx, flag)
	if err != nil {
		log.Println("VideolikedWithTransaction failed")
		tx.Rollback()
		return err
	}
	err = f.UserlikeWithTransaction(uid, tx, flag)
	if err != nil {
		log.Println("UserlikeWithTransaction failed")
		tx.Rollback()
		return err
	}
	//err = f.UserlikedWithTransaction(vid, tx, flag)
	//if err != nil {
	//	log.Println("UserlikedWithTransaction failed")
	//	tx.Rollback()
	//	return err
	//}
	return nil
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

// 第一步：更新favorite表
func (f *FavoriteService) FavoriteWithTransaction(uid int64, vid int64, tx *gorm.DB, flag int32) error {
	//如果flag==1添加
	fmt.Println(flag)
	if flag == 1 {
		exsit, err := f.FavoriteDao.IsExsitFavorite(uid, vid)
		if err != nil {
			return err
		}
		if exsit == false {
			favorite := dao.Favorite{
				UserId:     uid,
				VideoId:    vid,
				ActionType: 1, //点赞
				//CreateTime: time.Now().Format("2006-01-02 15:04:05"),
				//UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			err = f.FavoriteDao.InsertFavorite(tx, favorite)
			if err != nil {
				log.Println("have a InsertFavorite error")
				return err
			}
		} else {
			err := f.FavoriteDao.UpdateFavorite(tx, uid, vid, 1)
			if err != nil {
				return err
			}
		}
	} else { //如果flag==-1删除
		err := f.FavoriteDao.UpdateFavorite(tx, uid, vid, 2)
		if err != nil {
			log.Println("have a UpdateFavorite error")
			return err
		}
	}
	return nil
}

// 第二步：更新视频获赞数
func (f *FavoriteService) VideolikedWithTransaction(vid int64, tx *gorm.DB, flag int32) error {
	err := f.VideoDao.LikeVideo(tx, vid, int64(flag))
	if err != nil {
		return err
	}
	return nil
}

// 第三步：更新用户点赞数
func (f *FavoriteService) UserlikeWithTransaction(uid int64, tx *gorm.DB, flag int32) error {
	err := f.UserDao.UpdateUserFavoriteCount(tx, uid, int64(flag))
	if err != nil {
		log.Println("UpdateUserFavoriteCount have a error")
		return err
	}
	return nil
}

// 第四步：更新用户获赞数
func (f *FavoriteService) UserlikedWithTransaction(vid int64, tx *gorm.DB, flag int32) error {
	var (
		err    error
		userId int64
	)
	userId, err = f.FavoriteDao.GetUserIdbyVideoId(vid)
	if err != nil {
		log.Println("UserlikedWithTransaction.GetUserIdbyVideoId have a error")
		return err
	}

	err = f.UserDao.UpdateUserTotalFavorited(tx, userId, int64(flag))
	if err != nil {
		log.Println("UserlikedWithTransaction.UpdateUserFavoriteCount have a error")
		return err
	}
	return nil
}
