package controller

import (
	"github.com/luuuweiii/RiceDouyin/dao"
	"github.com/luuuweiii/RiceDouyin/service"
)

type Handler struct {
	FavoriteHandler *FavoriteHandler
	VideoHandler    *VideoHandler
	UserHandler     *UserHandler
	//别的模块handler
}

func NewHandler(
	FavoriteHandler *FavoriteHandler,
	VideoHandler *VideoHandler,
	UserHandler *UserHandler,
) *Handler {
	return &Handler{
		FavoriteHandler: FavoriteHandler,
		VideoHandler:    VideoHandler,
		UserHandler:     UserHandler,
	}
}
func BuildInjector() (*Handler, error) {
	db, err := dao.Init()
	if err != nil {
		return nil, err
	}

	//dao
	favoriteDao := dao.NewFavoriteDao(db)
	userDao := dao.NewUserDao(db)
	videoDao := dao.NewVideoDao(db)

	//service
	userService := service.NewUserService(userDao)
	videoService := service.NewVideoService(videoDao)
	favoriteService := service.NewFavoriteService(favoriteDao, userDao, videoDao, videoService)

	//handler
	favoriteHandler := NewFavoriteHandler(favoriteService)
	userHandler := NewUserHandler(userService)
	videoHandler := NewVideoHandler(videoService)

	//处理层
	handler := NewHandler(favoriteHandler, videoHandler, userHandler)
	return handler, nil
}
