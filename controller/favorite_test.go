package controller

//func TestFavorite(t *testing.T) {
//	e := newExpect(t)
//
//	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
//	feedResp.Value("status_code").Number().Equal(0)
//	feedResp.Value("video_list").Array().Length().Gt(0)
//	firstVideo := feedResp.Value("video_list").Array().First().Object()
//	videoId := firstVideo.Value("id").Number().Raw()
//
//	userId, token := getTestUserToken(testUserA, e)
//
//	favoriteResp := e.POST("/douyin/favorite/action/").
//		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).
//		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 1).
//		Expect().
//		Status(http.StatusOK).
//		JSON().Object()
//	favoriteResp.Value("status_code").Number().Equal(0)
//
//	favoriteListResp := e.GET("/douyin/favorite/list/").
//		WithQuery("token", token).WithQuery("user_id", userId).
//		WithFormField("token", token).WithFormField("user_id", userId).
//		Expect().
//		Status(http.StatusOK).
//		JSON().Object()
//	favoriteListResp.Value("status_code").Number().Equal(0)
//	for _, element := range favoriteListResp.Value("video_list").Array().Iter() {
//		video := element.Object()
//		video.ContainsKey("id")
//		video.ContainsKey("author")
//		video.Value("play_url").String().NotEmpty()
//		video.Value("cover_url").String().NotEmpty()
//	}
//}
