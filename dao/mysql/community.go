package mysql

import (
	"bluebell/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	err = db.Model(&models.Community{}).
		Select("community_id", "community_name").
		Find(&communityList).Error
	if err != nil {
		return nil, err
	}
	if len(communityList) == 0 {
		zap.L().Warn("there is no community in db")
	}
	return nil, err
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail{})
	err = db.Model(&models.CommunityDetail{}).Select("community_id,community_name", "introduction", "create_time").Where("community_id=?", id).Find(community).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrorInvalidID
		}
		return nil, err
	}
	return community, err
}
