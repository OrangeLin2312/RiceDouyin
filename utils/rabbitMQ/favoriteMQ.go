package rabbitMQ

import (
	"encoding/json"
	"fmt"
	"github.com/luuuweiii/RiceDouyin/config"
	"github.com/luuuweiii/RiceDouyin/dao"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
)

type FavoriteMq struct {
	conn        *amqp.Connection
	FavoriteDao *dao.FavoriteDao
	UserDao     *dao.UserDao
	VideoDao    *dao.VideoDao
}

func NewFavoriteMq(conn *amqp.Connection, UserDao *dao.UserDao,
	VideoDao *dao.VideoDao, f *dao.FavoriteDao) *FavoriteMq {
	return &FavoriteMq{
		conn:        conn,
		FavoriteDao: f,
		UserDao:     UserDao,
		VideoDao:    VideoDao,
	}
}

func ChannelClose(ch *amqp.Channel) {
	ch.Close()
}

func (f *FavoriteMq) FavoritePublish(body []byte) {
	ch, err := f.conn.Channel()
	if err != nil {
		fmt.Printf("open a channel failed, err:%v\n", err)
		panic(err)
	}
	// 3. 要发送，我们必须声明要发送到的队列。
	q, err := ch.QueueDeclare(
		"favorite", // name
		true,       // 持久的
		false,      // delete when unused
		false,      // 独有的
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		fmt.Printf("declare a queue failed, err:%v\n", err)
		return
	}

	// 4. 然后我们可以将消息发布到声明的队列
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // 立即
		false,  // 强制
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		fmt.Printf("publish a message failed, err:%v\n", err)
		return
	}
	log.Printf(" [x] Sent %s", body)
}

func (f *FavoriteMq) FavoriteWorker() {
	//声明一个channal
	ch, err := f.conn.Channel()
	if err != nil {
		fmt.Printf("open a worker channel failed, err:%v\n", err)
		return
	}
	// 声明一个queue
	q, err := ch.QueueDeclare(
		"favorite", // name
		true,       // 声明为持久队列
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Printf("ch.Qos() failed, err:%v\n", err)
		return
	}

	// 立即返回一个Delivery的通道
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // 注意这里传false,关闭自动消息确认
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Printf("ch.Consume failed, err:%v\n", err)
		return
	}

	// 开启循环不断地消费消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)
			var favoriteMsg config.RmqMessage
			err := json.Unmarshal(d.Body, &favoriteMsg)
			if err != nil {
				fmt.Println("have a Unmarshal error")
				return
			}
			if favoriteMsg.ActionType == 1 {
				f.FavoriteTureWorker(favoriteMsg)
			} else {
				f.FavoriteFalseWorker(favoriteMsg)
			}
			d.Ack(false) // 手动传递消息确认
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func (f *FavoriteMq) FavoriteTureWorker(msg config.RmqMessage) error {
	db := f.FavoriteDao.DB
	tx := db.Begin() // 开启事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()
	var flag int32 = 1
	err := f.ToFavoriteAction(tx, msg.UserId, msg.VideoId, flag)
	if err != nil {
		log.Println("FavoriteAciton failed")
	}
	tx.Commit()
	return err
}
func (f *FavoriteMq) FavoriteFalseWorker(msg config.RmqMessage) error {
	db := f.FavoriteDao.DB
	tx := db.Begin() // 开启事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()
	var flag int32 = -1
	err := f.ToFavoriteAction(tx, msg.UserId, msg.VideoId, flag)
	if err != nil {
		log.Println("FavoriteAciton failed")
	}
	tx.Commit()
	return err
}

// 连接事务对象和事务操作
func (f *FavoriteMq) ToFavoriteAction(tx *gorm.DB, uid int64, vid int64, flag int32) error {
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

// 第一步：更新favorite表
func (f *FavoriteMq) FavoriteWithTransaction(uid int64, vid int64, tx *gorm.DB, flag int32) error {
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
			}
			err = f.FavoriteDao.InsertFavorite(tx, favorite)
			if err != nil {
				fmt.Println("have a InsertFavorite error")
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
			fmt.Println("have a UpdateFavorite error")
			return err
		}
	}
	return nil
}

// 第二步：更新视频获赞数
func (f *FavoriteMq) VideolikedWithTransaction(vid int64, tx *gorm.DB, flag int32) error {
	err := f.VideoDao.LikeVideo(tx, vid, int64(flag))
	if err != nil {
		return err
	}
	return nil
}

// 第三步：更新用户点赞数
func (f *FavoriteMq) UserlikeWithTransaction(uid int64, tx *gorm.DB, flag int32) error {
	err := f.UserDao.UpdateUserFavoriteCount(tx, uid, int64(flag))
	if err != nil {
		fmt.Println("UpdateUserFavoriteCount have a error")
		return err
	}
	return nil
}

// 第四步：更新用户获赞数
func (f *FavoriteMq) UserlikedWithTransaction(vid int64, tx *gorm.DB, flag int32) error {
	var (
		err    error
		userId int64
	)
	userId, err = f.FavoriteDao.GetUserIdbyVideoId(vid)
	if err != nil {
		fmt.Println("UserlikedWithTransaction.GetUserIdbyVideoId have a error")
		return err
	}

	err = f.UserDao.UpdateUserTotalFavorited(tx, userId, int64(flag))
	if err != nil {
		fmt.Println("UserlikedWithTransaction.UpdateUserFavoriteCount have a error")
		return err
	}
	return nil
}
