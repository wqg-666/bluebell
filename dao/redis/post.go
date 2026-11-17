package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//根据用户选择的排序方式（时间排序或分数排序），从 Redis 中查询对应的帖子 ID，并按照分页要求返回帖子 ID 列表。

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//1.根据用户请求中携带的order参数确定要查询的order key参数

	// 默认按照时间排序,就把Key设置为帖子时间对应的 Redis Key。
	key := getRedisKey(KeyPostTimeZSet)
	// 如果用户选择分数排序
	if p.Order == models.OrderScore {
		//就把Key修改为帖子分数对应的 Redis Key。
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.根据排序集合和分页参数获取帖子 ID
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64,0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	//查找key中分数显示是1的元素的数量->统计每篇帖子投赞成票的票数
	//	//ZCOUNT key min max
	//	v1 := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	//使用pipeline一次发送多条命令，减少RTT，减少网络连接的次数
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		//v := cmder.(*redis.IntCmd).Val() 就是从 Redis 执行的 ZCOUNT 命令结果中，取出赞成票的数量。
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	//2.确定查询的索引qishidian
	start := (page - 1) * size
	end := start + size - 1
	//3.ZREVRANGE按分数从大到小查询
	return client.ZRevRange(key, start, end).Result()
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//使用 zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	//针对新的 zset 按之前的逻辑取顺序

	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算
		pipeline := client.Pipeline()
		//key	存储交集结果的目标 Redis Key
		//redis.ZStore{Aggregate: "MAX"}	指定交集后，相同成员的分数如何计算
		//cKey, orderKey	参与求交集的两个 ZSet
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) //zinterstore计算
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
