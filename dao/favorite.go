package dao

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type Favorite struct {
	gorm.Model
	//Id         int64 `gorm:"id"`
	VideoId    int64 `gorm:"video_id"`
	UserId     int64 `gorm:"user_id"`
	ActionType int64 `gorm:"action_type"`
	//CreateTime string `gorm:"create_time"`
	//UpdateTime string `gorm:"update_time"`
}

func NewFavoriteDao(db *gorm.DB) *FavoriteDao {
	return &FavoriteDao{DB: db}
}

type FavoriteDao struct {
	DB *gorm.DB
}

func (Favorite) TableName() string {
	return "favorite"
}

// 根据videoId得到userId
func (fD *FavoriteDao) GetUserIdbyVideoId(vid int64) (int64, error) {
	var user User
	err := fD.DB.Model(Video{}).Where("video_id=?", vid).First(&user).Error
	if err != nil {
		log.Println("have a GetUserIdbyVideoId error:", err)
		return -1, err
	}
	return user.UserId, err
}

// 添加一条点赞记录
func (fD *FavoriteDao) InsertFavorite(tx *gorm.DB, favorite Favorite) error {
	err := tx.Model(Favorite{}).Create(&favorite).Error
	if err != nil {
		log.Println("have a dao.InsertFavorite error:", err)
	}
	return err
}

//// 删除一条点赞记录
//func DelFavorite(tx *gorm.DB, uid int64, vid int64) error {
//	err := tx.Model(Favorite{}).Where("user_id=? and video_id=?", uid, vid).Delete(Favorite{}).Error
//	if err != nil {
//		log.Println("have a dao.DelFavorite error:", err)
//	}
//	return err
//}

// 查询点赞记录是否存在
func (fD *FavoriteDao) IsExsitFavorite(uid int64, vid int64) (bool, error) {
	var favorite Favorite
	err := fD.DB.Model(Favorite{}).Where("user_id=? and video_id=?", uid, vid).First(&favorite).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//log.Println(err)
			return false, nil
		}
		log.Println("have a dao.IsExsitFavorite error:", err)
		return false, err
	}
	return true, nil
}

// 更新点赞
func (fD *FavoriteDao) UpdateFavorite(tx *gorm.DB, uid int64, vid int64, actionType int32) error {
	res := tx.Model(Favorite{}).Where("user_id=? and video_id=?", uid, vid).
		Update("action_type", actionType)
	//.Update("update_time", time.Now().Format("2006-01-02 15:04:05"))
	err := res.Error
	if err != nil {
		log.Println("have a dao.UpdateFavorite error:", err)
		return err
	} else if res.RowsAffected == 0 {
		log.Println("havr no update")
		return errors.New("not exsit or have a same data")
	}
	return nil
}

// 获取视频id列表
func (fD *FavoriteDao) GetVideoIdList(uid int64) ([]int64, error) {
	var videoIdList []int64
	err := fD.DB.Model(&Favorite{}).Where("user_id=?", uid).Pluck("video_id", &videoIdList).Error
	if err != nil {
		log.Println("have a dao.GetVideoIdList error:", err)
	}
	return videoIdList, err
}

// 查询是否处于uid用户给vid视频点赞状态
func (fD *FavoriteDao) QueryIsFavorite(uid int64, vid int64) (bool, error) {
	var favorite Favorite
	err := fD.DB.Model(Favorite{}).Where("user_id=? and video_id=?", uid, vid).First(&favorite).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		log.Println("have a dao.QueryIsFavorite error:", err)
		return false, err
	} else {
		actiontype := favorite.ActionType
		if actiontype == 1 {
			return true, nil
		} else {
			return false, nil
		}
	}
}
