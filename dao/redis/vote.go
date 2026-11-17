package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//本项目使用简化版的投票分数
//投一票加432分 86400/200=432，一天86400秒，200张赞成票可以给你的帖子续一天，时间戳增大

/*
投票的几种情况
direction=1
1.之前没有投过票，现在投赞成票 -->更新分数和投票记录 差值的绝对值：1 +432
2.之前投反对票，现在改投赞成票 -->更新分数和投票记录 差值的绝对值：2 +432*2
direction=0
1.之前投过反对票，现在要取消投票 -->更新分数和投票记录 差值的绝对值：1 +432

2.之前投过赞成票，现在要取消投票 -->更新分数和投票记录 差值的绝对值：1 -432
direction=-1
1.之前没有投过票，现在反对成票 -->更新分数和投票记录 差值的绝对值：1 -432
2.之前投赞成票，现在改投反对成票 -->更新分数和投票记录 差值的绝对值：2 -432*2

投票的限制：
每个帖子自发表之日起一个星期内允许投票，超过一个星期就不允许再投票了
1.到期之后将redis中保存的赞成票数及反对票数储存到mysql表中
2.到期之后删除那个KeyPostVotedZSetPF
*/
const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrPostNotFound   = errors.New("帖子不存在或未写入 Redis")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	//事务操作。这两个操作要么同时成功要么同时失败
	//创建事务 Pipeline
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数，投一票续432秒
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}
func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	//根据帖子 postID，从 Redis 的一个 ZSet 中获取该帖子的 Score（分数），并赋值给 postTime,(postTime：帖子发布时间戳)
	postTime, err := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Result()
	if err == redis.Nil {
		return ErrPostNotFound
	}
	if err != nil {
		return err
	}
	//要转换成 Unix 时间戳？因为 Unix 时间戳是一个整数，方便存储、传输和比较。
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2和3需要放到一个pipeline事务中操作
	//2.更新帖子的分数
	//先查userID之前的投票记录，KeyPostVotedZSetPF = "post:voted:"这个要求传入一个参数postID，帖子ID
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	//这次投票值跟之前一样，提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1 //增加分数。
	} else {
		op = -1 //减少分数
	}
	//原来是反对票:ov = -1 现在改成赞成票:value = 1
	//因为:1 > -1(true)所以：op = 1表示增加分数。
	//再看反过来的情况：
	//原来是赞成票ov = 1，现在改成反对票：value = -1
	//因为： -1 > 1(false) 所以：op = -1表示减少分数。
	//例如：
	//ov = -1 ,value = 1 diff := math.Abs(-1 - 1)=2，op=1,,op*diff*scorePerVote=1*2*432
	diff := math.Abs(ov - value)    //计算两次投票的差值
	pipeline := client.TxPipeline() //创建一个事务
	//给帖子postID增加或减少分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	//3.记录用户为该帖子投票的数据
	if value == 0 { //如果用户取消投票
		//就把这个用户从投票记录中删除。
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else { //如果这个用户原来已经存在，ZAdd会更新他的Score，而不是再增加一个重复成员。
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err = pipeline.Exec()
	return err
}
