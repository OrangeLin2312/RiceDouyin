package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/luuuweiii/RiceDouyin/service"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type LikeListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}
type FavoriteHandler struct {
	FavoriteSerive *service.FavoriteService
}

func NewFavoriteHandler(favoriteservice *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{FavoriteSerive: favoriteservice}
}

// FavoriteAction no practical effect, just check if token is valid
func (fh *FavoriteHandler) FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	if user, exist := usersLoginInfo[token]; exist {
		var (
			videoId    int64
			userId     int64
			actionType int32
		)

		userId = user.Id
		strvideoId := c.Query("video_id")
		videoId, _ = strconv.ParseInt(strvideoId, 10, 64)
		stractionType := c.Query("action_type")
		actionTypetmp, _ := strconv.ParseInt(stractionType, 10, 32)
		actionType = int32(actionTypetmp)

		err := fh.FavoriteSerive.FavoriteAction(userId, videoId, actionType)
		if err != nil {
			log.Println("favorite.FavoriteAciton failed")
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteAction failed",
			})
		} else {
			log.Println("favorite.FavoriteAciton success")
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "FavoriteAction success"})
		}
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
	}
}

// FavoriteList all users have same favorite video list
func (fh *FavoriteHandler) FavoriteList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; exist {
		var (
			pageuserId int64
			//userId     int64
			err error
		)
		//userId = user.Id
		strpguserId := c.Query("user_id ")
		pageuserId, _ = strconv.ParseInt(strpguserId, 10, 64)
		videoList, err := fh.FavoriteSerive.FavoriteList(pageuserId)
		if err != nil {
			log.Println("favorite.FavoriteList failed")
			c.JSON(http.StatusOK, LikeListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Get favoritelist failed",
				},
			})
		} else {
			log.Println("favorite.FavoriteList success")
			c.JSON(http.StatusOK, LikeListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Get favoritelist success",
				},
				VideoList: videoList, //因为结构体在controller里，之后我感觉肯定要改
			})
		}

	} else {
		c.JSON(http.StatusOK, LikeListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
	}
}
