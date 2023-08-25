package dao

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

var db *gorm.DB

func clearTables(db *gorm.DB) {
	db.Exec("truncate favorite")
}
func setupDB() {
	var err error
	db, err = Init()
	fmt.Printf("初始化数据库成功-----\n")
	if err != nil {
		panic(err)
	}
}
func TestMain(m *testing.M) {
	//clearTables(db)
	setupDB()
	m.Run()
	//clearTables(db)
}

func TestAllFavoriteDao(t *testing.T) {
	t.Run("insert", testInertFt)
	t.Run("getlist", testGetVideoIdList)
	t.Run("queryIsF", testQueryIsFavorite)
	t.Run("exsit", testIsExsitFavorite)
	t.Run("update", testUpdateFavorite)
	//	t.Run("video",TestGetUserIdbyVideoId)

}

func testInertFt(t *testing.T) {
	//db := setupDB()
	fdao := NewFavoriteDao(db)
	f := Favorite{
		VideoId:    2344,
		UserId:     1234,
		ActionType: 1,
	}
	err := fdao.InsertFavorite(db, f)
	if err != nil {
		t.Errorf("error of InsertFavorite %v", err)
	}
}

//	func TestGetUserIdbyVideoId(t *testing.T) {
//		db := setupDB()
//		fdao := NewFavoriteDao(db)
//		_,err:=fdao.GetUserIdbyVideoId(123)
//		if err != nil {
//			t.Errorf("error of InsertFavorite %v", err)
//		}
//	}
func testGetVideoIdList(t *testing.T) {
	//db := setupDB()
	fdao := NewFavoriteDao(db)
	list, err := fdao.GetVideoIdList(1234)
	if list[0] != 2344 || err != nil {
		t.Errorf("error of GetVideoIdList %v", err)
	}
}
func testIsExsitFavorite(t *testing.T) {
	//db := setupDB()
	fdao := NewFavoriteDao(db)
	exsit, err := fdao.IsExsitFavorite(0, 2344)
	if exsit == true || err != nil {
		t.Errorf("error of IsExsitFavorite %v", err)
	}
	exsit, err = fdao.IsExsitFavorite(1234, 2344)
	if exsit == false || err != nil {
		t.Errorf("error of IsExsitFavorite %v", err)
	}
}
func testUpdateFavorite(t *testing.T) {
	//db := setupDB()
	fdao := NewFavoriteDao(db)
	err := fdao.UpdateFavorite(db, 1, 2344, 2)
	if err == nil {
		t.Errorf("error of UpdateFavorite %v", err)
	}
	err = fdao.UpdateFavorite(db, 1234, 2344, 2)
	if err != nil {
		t.Errorf("error of UpdateFavorite %v", err)
	}
}
func testQueryIsFavorite(t *testing.T) {
	//db := setupDB()
	fdao := NewFavoriteDao(db)
	ik, err := fdao.QueryIsFavorite(1234, 2344)
	if ik != true || err != nil {
		t.Errorf("error of QueryIsFavorite %v", err)
	}
	ik, err = fdao.QueryIsFavorite(124, 2344)
	if ik == true || err != nil {
		t.Errorf("error of QueryIsFavorite %v", err)
	}
}
