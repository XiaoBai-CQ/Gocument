package dao

import (
	"Gocument/app/api/global"
	"Gocument/app/api/internal/model"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CheckUserInRedis(username string) (bool, error) {
	redisName := "user:" + username

	_, err := global.RedisDB.Get(global.Ctx, redisName).Result()
	if err == nil {
		//键值对存在
		return true, nil
	} else if err != redis.Nil {
		//发生其他错误
		global.Logger.Error("redis query wrong")
		return false, err
	}
	//键值对不存在
	return false, nil
}

func CheckUserInMysql(username string) (bool, error) {
	var user model.User
	result := global.MysqlDB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			//未查询到记录
			global.Logger.Debug("user does not exist")
			return false, nil
		} else {
			//mysql错误(第三方错误需要转化)
			global.Logger.Error("Mysql failed to query existing user", zap.Error(result.Error))
			return false, result.Error
		}
	} else {
		return true, nil
	}
}

func CheckUserAndFilename(filename, username string) (*model.File, error) {
	var file model.File

	err := global.MysqlDB.Where("file_name = ? AND username = ?", filename, username).First(&file).Error
	if err != nil {
		return &model.File{}, err
	}
	return &file, nil
}
