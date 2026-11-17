package redis

// redis key
// redis key注意使用命名空间的方式，方便查询和拆分
const (
	KeyPrefix          = "bluebell"
	KeyPostTimeZSet    = "post:time"   //zset:帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  //zset:帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset:记录用户及投票的类型,这个 Key 后面还要拼接帖子ID,参数是post_id
	KeyCommunitySetPF  = "community:"  //set 保存每个分区下帖子的id
)

// 给redis加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + ":" + key
}
