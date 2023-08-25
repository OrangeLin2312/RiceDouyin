package service

import (
	"github.com/luuuweiii/RiceDouyin/dao"
	"gorm.io/gorm"
	"testing"
)

var (
	db   *gorm.DB
	fdao *dao.FavoriteDao
)

func clearTables() {
	db.Exec("truncate favorite")
}
func TestMain(m *testing.M) {
	Setup()
	//clearTables()
	m.Run()
	//clearTables()
}
func Setup() {
	var err error
	db, err = dao.Init()
	if err != nil {
		panic(err)
	}

	fdao = dao.NewFavoriteDao(db)
}

func TestFavoriteServiceIml_FavoriteAction(t *testing.T) {
	fser := NewFavoriteService(fdao, nil, nil, nil)
	err := fser.FavoriteAction(123, 345, 1)
	if err != nil {
		t.Errorf("error of FavoriteAction %v", err)
	}
	err = fser.FavoriteAction(123, 345, 2)
	if err != nil {
		t.Errorf("error of FavoriteAction %v", err)
	}
}

//func TestFavoriteServiceIml_FavoriteList(t *testing.T) {
//	fser := NewFavoriteService(fdao, nil, nil, nil)
//	//暂时先—忽略一下
//	_, err := fser.FavoriteList(123)
//	if err != nil {
//		t.Errorf("error of InsertFavorite %v", err)
//	}
//}
