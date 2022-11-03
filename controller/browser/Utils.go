package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"errors"

	"gorm.io/gorm"
)

func CheckPermission(user model_v1.User, project_key string, perms []int) (*model_v1.Project, error) {
	project := model_v1.Project{}
	result := database.DB.Where("key = ?", project_key).First(&project)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("no project")
		} else {
			return nil, errors.New("server error")
		}
	}

	if project.UserRID == user.RID {
		return &project, nil
	}

	member := model_v1.Member{}
	result = database.DB.Where("project_r_id = ? AND member_r_id = ?", project.Key, user.RID).First(&member)
	if result.Error != nil { // 该用户不是该项目的成员
		return nil, errors.New("project no member")
	} else {
		for _, perm := range perms {
			if perm == member.Status {
				return &project, nil
			}
		}
		// 该成员没有上传权限
		return nil, errors.New("no permission")
	}
}
