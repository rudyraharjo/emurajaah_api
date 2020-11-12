package repository

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/group"
	"github.com/rudyraharjo/emurojaah/models"
)

type postgreGroupRepository struct {
	DbCon *gorm.DB
}

func NewGroupRepository(db *gorm.DB) group.Repository {
	return &postgreGroupRepository{db}
}

func (repo *postgreGroupRepository) ListAvailableGroupByType(category string) ([]models.Group, error) {

	limitUser := 0
	limitNoGroupIndex := 0
	switch category {
	case "Juz":
		limitUser = 30
		limitNoGroupIndex = 10
	case "Ayat":
		limitUser = 6236
		limitNoGroupIndex = 1
	default:
		limitUser = 604 //Halaman
		limitNoGroupIndex = 5
	}

	var groups []models.Group
	err := repo.DbCon.Table("groups").Where("group_type = ? and current_member < ? and no_group_index > ?", category, limitUser, limitNoGroupIndex).First(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo *postgreGroupRepository) GetAvailableGroupByTypeNew(category string) (*models.Group, error) {

	var groups []models.Group
	now := time.Now()
	groupName := ""
	limitUser := 0
	limitNoGroupIndex := 0

	switch category {
	case "Juz":
		limitUser = 30
		groupName = "Group Juz"
		limitNoGroupIndex = 10
	case "Ayat":
		limitUser = 6236
		groupName = "Group Ayat"
		limitNoGroupIndex = 1
	default:
		limitUser = 604
		groupName = "Group Halaman"
		limitNoGroupIndex = 5
	}

	err := repo.DbCon.Raw(`select G.id, G.no_group_index, G.group_type,G.group_name, G.max_member, count(GM.group_id) as current_member,
	G.created_at, G.updated_at from groups G LEFT JOIN group_members GM on G.id = GM.group_id WHERE G.group_type = ? AND current_member < G.max_member AND G.no_group_index > ? GROUP BY G.id, G.current_member, GM.group_id, G.no_group_index, G.group_type,G.group_name order by created_at asc limit 1`, category, limitNoGroupIndex).Scan(&groups).Error
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {

		getgroups, err2 := repo.GetLastNoIndexGroup(category)

		if err != nil {
			return nil, err2
		}

		newGroup := models.Group{
			NoGroupIndex:  getgroups[0].NoGroupIndex + 1,
			GroupName:     groupName,
			MaxMember:     limitUser,
			CurrentMember: 0,
			GroupType:     category,
			CreatedAt:     now,
		}

		groupID, err3 := repo.CreateGroup(newGroup)
		if err3 != nil {
			return nil, err2
		}

		newGroup.Id = groupID
		return &newGroup, nil
	}

	return &groups[0], nil
}

func (repo *postgreGroupRepository) GetAvailableGroupByID(groupID int) (bool, error) {

	var groups []models.Group
	err := repo.DbCon.Table("groups").Where("id = ?", groupID).First(&groups).Error
	if err != nil {
		return false, err
	}

	if len(groups) > 0 {
		return true, nil
	} else {
		return false, nil
	}

}

func (repo *postgreGroupRepository) GetLastNoIndexGroup(category string) ([]models.Group, error) {

	var groups []models.Group
	err := repo.DbCon.Raw("select * from groups where group_type = ? order by no_group_index desc limit 1", category).Scan(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo *postgreGroupRepository) CreateGroup(g models.Group) (int, error) {
	var newGroup models.Group
	err := repo.DbCon.Table("groups").Create(&g).Scan(&newGroup).Error
	if err != nil {
		return 0, err
	}
	return newGroup.Id, nil
}

func (repo *postgreGroupRepository) ListInactiveGroupMember(groupID int) ([]models.GroupMember, error) {
	var inactiveMembers []models.GroupMember
	err := repo.DbCon.Where("group_id = ? and is_active = 0", groupID).Order("current_index asc").Find(&inactiveMembers).Error
	if err != nil {
		return nil, err
	}
	return inactiveMembers, nil
}

// func (repo *postgreGroupRepository) GetListGroupMembersByIDGroup(groupID int) ([]models.GroupMember, error) {
// 	var groupMembers []models.GroupMember

// 	err := repo.DbCon.Where("group_id = ?", groupID).Order("current_index asc").Find(&groupMembers).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return groupMembers, nil
// }

func (repo *postgreGroupRepository) GetListGroupMembersByIDGroupAndCurrenIndex(groupID int, genNumberIndex int) (int, error) {
	//var groupMembers []models.GroupMember
	var CountgroupMembers int

	fmt.Print("GroupID ", groupID)
	fmt.Print("\n")
	fmt.Print("genNumberIndex ", genNumberIndex)
	fmt.Print("\n")

	err := repo.DbCon.Raw("select count(*) from group_members where group_id = ? and current_index = ?", groupID, genNumberIndex).Count(&CountgroupMembers).Error
	if err != nil {
		return 1, err
	}

	fmt.Print("CountgroupMembers ", CountgroupMembers)
	return CountgroupMembers, nil
}

func (repo *postgreGroupRepository) ListActiveGroupMember(groupID int) ([]models.GroupMember, error) {
	var activeMembers []models.GroupMember

	err := repo.DbCon.Where("group_id = ? and is_active = 1", groupID).Order("current_index asc").Find(&activeMembers).Error
	if err != nil {
		return nil, err
	}
	return activeMembers, nil
}

func (repo *postgreGroupRepository) AddGroupMemberByEmail(groupID int, email string, grouptype string) (*models.GroupMember, int, error) {

	var groupMemberInsert models.GroupMember
	var userAlias []models.UserAliasFull
	var CheckCountGroupAvailable int

	now := time.Now()
	reservedIndex := 0
	max := 0

	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, 1, err
	}

	if err := tx.Table("user_alias").Where("alias = ? ", email).Scan(&userAlias).Error; err != nil {
		tx.Rollback()
		return nil, 1, err
	}

	if len(userAlias) > 0 {

		if strings.Title(grouptype) == "Ayat" {

			max = 6236

			if err := tx.Table("group_members").Where("group_id = ? and group_type = ? ", groupID, strings.Title(grouptype)).Count(&CheckCountGroupAvailable).Error; err != nil {
				tx.Rollback()
				return nil, 1, err
			}

			// fmt.Print("CheckCountGroupAvailable Ayat", CheckCountGroupAvailable)
			// fmt.Print("\n")

			if CheckCountGroupAvailable < max {

				genNumberIndex := repo.generateNumber(groupID, strings.Title(grouptype))

				if genNumberIndex != 0 {
					reservedIndex = genNumberIndex
				}

			} else {
				return nil, 3, tx.Commit().Error
			}

		} else if strings.Title(grouptype) == "Halaman" {
			max = 604

			if err := tx.Table("group_members").Where("group_id = ? and group_type = ? ", groupID, strings.Title(grouptype)).Count(&CheckCountGroupAvailable).Error; err != nil {
				tx.Rollback()
				return nil, 1, err
			}

			if CheckCountGroupAvailable < max {

				genNumberIndex := repo.generateNumber(groupID, strings.Title(grouptype))
				if genNumberIndex != 0 {
					reservedIndex = genNumberIndex
				}

			} else {
				return nil, 3, tx.Commit().Error
			}

		} else {

			max = 30

			if err := tx.Table("group_members").Where("group_id = ? and group_type = ? ", groupID, strings.Title(grouptype)).Count(&CheckCountGroupAvailable).Error; err != nil {
				tx.Rollback()
				return nil, 1, err
			}

			if CheckCountGroupAvailable < max {

				genNumberIndex := repo.generateNumber(groupID, strings.Title(grouptype))
				if genNumberIndex != 0 {
					reservedIndex = genNumberIndex
				}

			} else {
				return nil, 3, tx.Commit().Error
			}

		}

		newMember := models.GroupMember{
			CurrentIndex: reservedIndex,
			UserID:       userAlias[0].UserId,
			GroupID:      groupID,
			GroupType:    strings.Title(grouptype),
			IsDone:       0,
			CreatedAt:    now,
		}

		if err := tx.Create(&newMember).Scan(&groupMemberInsert).Error; err != nil {
			tx.Rollback()
			return nil, 1, err
		}

		if err := tx.Table("groups").Where("id = ?", groupID).Update("current_member", gorm.Expr("current_member + ?", 1)).Error; err != nil {
			tx.Rollback()
			return nil, 1, err
		}

	} else {
		return nil, 2, tx.Commit().Error
	}

	return &groupMemberInsert, 0, tx.Commit().Error
}

func (repo *postgreGroupRepository) AddGroupMemberBulk(reqBulk models.RequestJoinBulkGroup) (int, error) {

	now := time.Now()
	var groupMember []models.GroupMember
	var newMember models.GroupMember
	var groups []models.Group
	var userAlias []models.UserAliasFull
	reservedIndex := 0
	TempCurrentMember := 0
	limitUser := 0
	var CountGroupMember int

	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return 1, err
	}

	if err := tx.Raw(`select G.id, G.no_group_index, G.group_type,G.group_name, G.max_member, count(GM.group_id) as current_member,
	G.created_at, G.updated_at from groups G 
	LEFT JOIN group_members GM on G.id = GM.group_id 
	WHERE G.id = ? AND current_member < G.max_member GROUP BY G.id, G.current_member, GM.group_id, G.no_group_index, G.group_type,G.group_name order by created_at asc limit 1`, reqBulk.GroupID).Scan(&groups).Error; err != nil {
		tx.Rollback()
		return 1, err
	}

	fmt.Print(groups)
	fmt.Print("\n")

	if len(groups) > 0 {
		TempCurrentMember = groups[0].CurrentMember + len(reqBulk.DataBulk)
	} else {
		tx.Rollback()
		return 1, nil
	}

	switch groups[0].GroupType {
	case "Juz":
		limitUser = 30
	case "Ayat":
		limitUser = 6236
	default:
		limitUser = 604
	}

	if len(groups) > 0 && TempCurrentMember <= limitUser {

		if err := tx.Table("group_members").Where("group_id = ? and group_type = ?", groups[0].Id, groups[0].GroupType).Scan(&groupMember).Error; err != nil {
			tx.Rollback()
			return 1, err
		}

		for _, item := range reqBulk.DataBulk {

			genNumberIndex := repo.generateNumberNew(groupMember, groups[0].GroupType)
			if genNumberIndex != 0 {
				reservedIndex = genNumberIndex
			}

			if err := tx.Table("user_alias").Where("alias = ? ", item.Email).Scan(&userAlias).Error; err != nil {
				tx.Rollback()
				return 2, err
			}

			//fmt.Print(userAlias)
			if len(userAlias) > 0 {

				addMember := models.GroupMember{
					CurrentIndex: reservedIndex,
					UserID:       userAlias[0].UserId,
					GroupID:      groups[0].Id,
					GroupType:    groups[0].GroupType,
					IsDone:       0,
					CreatedAt:    now,
				}

				if err := tx.Create(&addMember).Scan(&newMember).Error; err != nil {
					tx.Rollback()
					return 1, err
				}

				if err := tx.Table("group_members").Where("group_id = ? ", groups[0].Id).Count(&CountGroupMember).Error; err != nil {
					tx.Rollback()
					return 1, err
				}

				if err := tx.Table("groups").Where("id = ?", groups[0].Id).Update("current_member", CountGroupMember).Error; err != nil {
					tx.Rollback()
					return 1, err
				}

			} else {

				tx.Rollback()
				return 2, nil

			}
		}

	} else {
		tx.Rollback()
		return 1, nil
	}

	return 0, tx.Commit().Error
}

func (repo *postgreGroupRepository) AddGroupMemberNew(userId int, category string) (*models.GroupMember, int, error) {

	now := time.Now()
	var newMember models.GroupMember
	var CountUserAlias int
	var CountGroupMember int
	reservedIndex := 0

	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, 1, err
	}

	if err := tx.Table("user_alias").Where("user_id = ? ", userId).Count(&CountUserAlias).Error; err != nil {
		tx.Rollback()
		return nil, 1, err
	}

	if CountUserAlias > 0 {

		groups, err := repo.GetAvailableGroupByTypeNew(category)
		if err != nil {
			tx.Rollback()
			return nil, 1, err
		}

		genNumberIndex := repo.generateNumber(groups.Id, category)
		if genNumberIndex != 0 {
			reservedIndex = genNumberIndex
		}

		addMember := models.GroupMember{
			CurrentIndex: reservedIndex,
			UserID:       userId,
			GroupID:      groups.Id,
			GroupType:    category,
			IsDone:       0,
			CreatedAt:    now,
		}

		if err := tx.Create(&addMember).Scan(&newMember).Error; err != nil {
			tx.Rollback()
			return nil, 1, err
		}

		if err := tx.Table("group_members").Where("group_id = ? ", groups.Id).Count(&CountGroupMember).Error; err != nil {
			tx.Rollback()
			return nil, 1, err
		}

		fmt.Print("CountGroupMember ", CountGroupMember)
		fmt.Print("\n")

		if err := tx.Table("groups").Where("id = ?", groups.Id).Update("current_member", CountGroupMember).Error; err != nil {
			tx.Rollback()
			return nil, 1, err
		}

	} else {
		return nil, 2, tx.Commit().Error
	}

	return &newMember, 0, tx.Commit().Error
}

func (repo *postgreGroupRepository) AddGroupMember(groupmembers models.GroupMember, groupID int) (*models.GroupMember, error) {

	var newMember models.GroupMember

	var GroupMember []models.GroupMember

	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&groupmembers).Scan(&newMember).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Table("group_members").Where("group_id = ? ", groupID).Scan(&GroupMember).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Table("groups").Where("id = ?", groupID).Update("current_member", len(GroupMember)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// if err := tx.Table("groups").Where("id = ?", groupID).Update("current_member", gorm.Expr("current_member + ?", 1)).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	return &newMember, tx.Commit().Error

}

func (repo *postgreGroupRepository) generateNumberNew(groupMember []models.GroupMember, categoryType string) int {

	min := 1
	max := 0
	genNumberIndex := 0
	mTempCurrentIndex := make(map[int]bool)

	for i := 0; i < len(groupMember); i++ {
		mTempCurrentIndex[groupMember[i].CurrentIndex] = true
	}

	if len(mTempCurrentIndex) > 0 {
		if categoryType == "Ayat" {

			max = 6236

			for {

				rand.Seed(time.Now().UnixNano())
				genNumberIndex = rand.Intn(max-min+1) + min

				if _, ok := mTempCurrentIndex[genNumberIndex]; ok {
					continue
				} else {
					break
				}
			}

		} else if categoryType == "Halaman" {

			max = 604

			for {

				rand.Seed(time.Now().UnixNano())
				genNumberIndex = rand.Intn(max-min+1) + min

				if _, ok := mTempCurrentIndex[genNumberIndex]; ok {
					continue
				} else {
					break
				}
			}

		} else {

			max = 30

			for {

				rand.Seed(time.Now().UnixNano())
				genNumberIndex = rand.Intn(max-min+1) + min

				if _, ok := mTempCurrentIndex[genNumberIndex]; ok {
					continue
				} else {
					break
				}
			}

		}
	} else {

		if categoryType == "Ayat" {

			max = 6236
			rand.Seed(time.Now().UnixNano())
			genNumberIndex = rand.Intn(max-min+1) + min

		} else if categoryType == "Halaman" {

			max = 604

			rand.Seed(time.Now().UnixNano())
			genNumberIndex = rand.Intn(max-min+1) + min

		} else {

			max = 30

			rand.Seed(time.Now().UnixNano())
			genNumberIndex = rand.Intn(max-min+1) + min

		}
	}

	return genNumberIndex

}

func (repo *postgreGroupRepository) generateNumber(groupID int, categoryType string) int {

	min := 1
	max := 0
	genNumberIndex := 0
	reservedIndex := 0
	var CheckCountGenNumber int

	if categoryType == "Ayat" {

		max = 6236

		for {

			rand.Seed(time.Now().UnixNano())
			genNumberIndex = rand.Intn(max-min+1) + min

			err := repo.DbCon.Table("group_members").Where("group_id = ? and current_index = ?", groupID, genNumberIndex).Count(&CheckCountGenNumber).Error
			if err != nil {
				return 0
			}

			if CheckCountGenNumber == 0 && genNumberIndex <= 6236 {
				reservedIndex = genNumberIndex
				break
			}

		}

	} else if categoryType == "Halaman" {

		max = 604

		for {

			rand.Seed(time.Now().UnixNano())
			genNumberIndex = rand.Intn(max-min+1) + min

			err := repo.DbCon.Table("group_members").Where("group_id = ? and current_index = ?", groupID, genNumberIndex).Count(&CheckCountGenNumber).Error
			if err != nil {
				return 0
			}

			if CheckCountGenNumber == 0 && genNumberIndex <= 604 {
				reservedIndex = genNumberIndex
				break
			}

		}

	} else {

		max = 30

		for {

			rand.Seed(time.Now().UnixNano())
			genNumberIndex = rand.Intn(max-min+1) + min

			err := repo.DbCon.Table("group_members").Where("group_id = ? and current_index = ?", groupID, genNumberIndex).Count(&CheckCountGenNumber).Error
			if err != nil {
				return 0
			}

			if CheckCountGenNumber == 0 && genNumberIndex <= 30 {
				reservedIndex = genNumberIndex
				break
			}

		}

	}

	return reservedIndex

	// rand.Seed(time.Now().UnixNano())
	// return rand.Intn(max-min+1) + min
}

func (repo *postgreGroupRepository) GetNoGroupIndex(groupID int) ([]models.Group, error) {
	var group []models.Group

	err := repo.DbCon.Table("groups").Where("id = ?", groupID).Scan(&group).Error
	if err != nil {
		return group, err
	}

	return group, nil
}

// func (repo *postgreGroupRepository) UpdateInactiveMember(groupId int, oldUserId int, newUserId int) (*models.GroupMember, error) {

// 	var updatedMember models.GroupMember

// 	tx := repo.DbCon.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	if err := tx.Error; err != nil {
// 		return nil, err
// 	}

// 	if err := tx.Table("group_members").Where("group_id = ? and user_id = ?", groupId, oldUserId).Updates(models.GroupMember{IsActive: 1, UserId: newUserId, UpdatedAt: time.Now()}).Scan(&updatedMember).Error; err != nil {
// 		fmt.Print("update group_members ", err)
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	if err := tx.Table("groups").Where("id = ?", groupId).Update("current_member", gorm.Expr("current_member + ?", 1)).Error; err != nil {
// 		fmt.Print("update groups current_member + 1")
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	return &updatedMember, tx.Commit().Error

// 	// Default Query

// 	// err := repo.DbCon.Table("group_members").Where("group_id = ? and user_id = ?", groupId, oldUserId).
// 	// 	Updates(models.GroupMember{IsActive: 1, UserId: newUserId, UpdatedAt: time.Now()}).Scan(&updatedMember).Error
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// }

func (repo *postgreGroupRepository) SetGroupMemberAsInactive(groupId int, userId int) error {

	var countGroupMember int
	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Table("group_members").Where("group_id = ? and user_id = ?", groupId, userId).Count(&countGroupMember).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("groups").Where("id = ?", groupId).Update("current_member", gorm.Expr("current_member - ?", countGroupMember)).Error; err != nil {
		tx.Rollback()
		return err
	}

	queryDeleteGroupMember := fmt.Sprintf(`DELETE FROM public.group_members WHERE group_id = $1 and user_id = $2;`)

	if err := repo.DbCon.Exec(queryDeleteGroupMember, groupId, userId).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

	// Default Query

	// err := repo.DbCon.Table("group_members").Where("group_id = ? and user_id = ?", groupId, userId).
	// 	Update("is_active", 0).Error
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (repo *postgreGroupRepository) ExitGroupAndInsertActivity(groupID int, userID int) error {
	var GroupMembers []models.GroupMember
	now := time.Now()

	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Table("group_members").Where("group_id = ? and user_id = ?", groupID, userID).Scan(&GroupMembers).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(GroupMembers) != 0 {

		if err := tx.Table("groups").Where("id = ?", groupID).Update("current_member", gorm.Expr("current_member - ?", len(GroupMembers))).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, GroupMember := range GroupMembers {

			insUserActivity := models.UserActivity{
				UserId:           GroupMember.UserID,
				IdGroupMember:    GroupMember.ID,
				GroupId:          GroupMember.GroupID,
				GroupType:        strings.Title(GroupMember.GroupType),
				ContentIndex:     GroupMember.CurrentIndex,
				StatusUserAction: 3,
				Description:      "Leave Group",
				CreatedAt:        now,
			}

			if err := tx.Table("user_activity").Create(&insUserActivity).Error; err != nil {
				tx.Rollback()
				return err
			}

		}

		queryDeleteGroupMember := fmt.Sprintf(`DELETE FROM public.group_members WHERE group_id = $1 and user_id = $2;`)

		if err := repo.DbCon.Exec(queryDeleteGroupMember, groupID, userID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (repo *postgreGroupRepository) LeaveReadingGroupAndInsertActivity(params models.RequestLeaveReading) error {

	tx := repo.DbCon.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	now := time.Now()

	insUserActivity := models.UserActivity{
		UserId:           params.UserID,
		IdGroupMember:    params.ID,
		GroupId:          params.GroupID,
		GroupType:        strings.Title(params.GroupType),
		ContentIndex:     params.ContentIndex,
		StatusUserAction: 2,
		Description:      "Leave Reading Group",
		CreatedAt:        now,
	}

	queryDeleteGroupMember := fmt.Sprintf(`DELETE FROM public.group_members WHERE id = $1;`)

	if err := repo.DbCon.Exec(queryDeleteGroupMember, params.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("groups").Where("id = ?", params.GroupID).Update("current_member", gorm.Expr("current_member - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("user_activity").Create(&insUserActivity).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *postgreGroupRepository) IsMemberExist(groupId int, userId int) (int, error) {
	var count int
	err := repo.DbCon.Table("group_members").Where("group_id = ? and user_id = ?", groupId, userId).Count(&count).Error
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 2, nil
	}
	return 1, nil
}

func (repo *postgreGroupRepository) IsUserHasGroupCategory(userId int, category string) (int, error) {
	var count int
	err := repo.DbCon.Raw("select a.group_id from group_members a inner join "+
		"groups b on a.group_id = b.id where b.group_type = ? and a.user_id = ? and a.is_active = 1", category, userId).Count(&count).Error
	if err != nil {
		return 2, err
	}

	if count == 0 {
		return 0, nil
	}
	return 1, nil
}

// ListUserGroup function
func (repo *postgreGroupRepository) ListUserGroup(userID int) ([]models.ResponsGroupeGroupMembersListWithNoUrut, error) {

	var groups []models.ResponsGroupeGroupMembersListWithNoUrut

	err := repo.DbCon.Raw("select gm.user_id,gm.group_id,gm.group_type, gp.no_group_index as no_urut from group_members gm LEFT JOIN groups gp ON gm.group_id=gp.id where user_id = ? group by gm.user_id,gm.group_id,gm.group_type, gp.no_group_index order by gp.no_group_index", userID).Scan(&groups).Error

	if err != nil {
		return nil, err
	}

	// for i := 0; i < len(sCan); i++ {

	// 	if sCan[i].GroupType == tempTypeGroup {
	// 		noUrut++
	// 	} else {
	// 		tempTypeGroup = sCan[i].GroupType
	// 		noUrut = 1
	// 	}

	// 	groups = append(groups, models.ResponsGroupeGroupMembersListWithNoUrut{
	// 		NoUrut:    noUrut,
	// 		UserID:    sCan[i].UserID,
	// 		GroupID:   sCan[i].GroupID,
	// 		GroupType: sCan[i].GroupType,
	// 	})

	// }

	return groups, nil
}

func (repo *postgreGroupRepository) GetListGroupMembersByUserIDAndGroupID(userID int, groupID int) (int, error) {

	// var groupMembers []models.GroupMember
	var CountgroupMembers int

	//err := repo.DbCon.Table("group_members").Where("user_id = ? and group_id = ? ", userID, groupID).Scan(&groupMembers).Error
	err := repo.DbCon.Table("group_members").Where("user_id = ?", userID).Count(&CountgroupMembers).Error

	if err != nil {
		return 0, err
	}

	return CountgroupMembers, nil
}

func (repo *postgreGroupRepository) ListAllGroup() ([]models.ResponseAllgroup, error) {
	var AllGroups []models.ResponseAllgroup

	err := repo.DbCon.Raw("select group_type, group_name from groups group by group_type,group_name").Scan(&AllGroups).Error

	if err != nil {
		return nil, err
	}

	return AllGroups, nil
}

func (repo *postgreGroupRepository) ListAllGroupIsJoined(userID int) ([]models.ResponseGroupUserJoined, error) {

	var scanGroup []models.ResponseAllgroup
	var responseGroup []models.ResponseGroupUserJoined

	err := repo.DbCon.Raw("select group_type, group_name from groups group by group_type,group_name").Scan(&scanGroup).Error

	if err != nil {
		return nil, err
	}

	for _, data := range scanGroup {
		var count int

		//fmt.Print("Check Loopin Group => ", scanGroup)
		checkCount := repo.DbCon.Table("group_members").Where("user_id = ? and group_type = ? ", userID, data.GroupType).Count(&count).Error

		if checkCount != nil {
			return nil, err
		}

		if count != 0 {

			responseGroup = append(responseGroup, models.ResponseGroupUserJoined{
				GroupType: data.GroupType,
				GroupName: data.GroupName,
				IsJoined:  true,
			})

		} else {
			responseGroup = append(responseGroup, models.ResponseGroupUserJoined{
				GroupType: data.GroupType,
				GroupName: data.GroupName,
				IsJoined:  true,
			})
		}
	}

	return responseGroup, nil
}

func (repo *postgreGroupRepository) GetListGroupType(userID int, groupType string, offset int, limit int) ([]models.ResponseListGroupTypeWithNoUrut, error) {

	//var sCan []models.ResponseListGroupType
	var data []models.ResponseListGroupTypeWithNoUrut

	// tempIDGroup := 0
	// noUrut := 0

	if limit == 0 {
		limit = 10
	}

	err := repo.DbCon.Raw(`SELECT gm.id, gm.group_id, gm.user_id,gm.current_index, CASE WHEN gm.is_done = 1 THEN true ELSE false END as is_reading, gp.group_name, gp.no_group_index as no_urut, gp.group_type, q.juz_id, q.surah_id, q.ayat_sec, qs.asma, qs.surah_name FROM group_members gm LEFT JOIN groups gp ON gm.group_id=gp.id LEFT JOIN quran q ON gm.current_index=q.id LEFT JOIN quran_surah qs ON q.surah_id=qs.id where gm.user_id = ? and gp.group_type = ? order by gm.is_done,gp.no_group_index,gm.current_index offset ? limit ?`, userID, strings.Title(groupType), offset, limit).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil

}

func (repo *postgreGroupRepository) ListGroupMemberWithName(groupId int) ([]models.ResponseGroupMemberList, error) {
	var data []models.ResponseGroupMemberList

	rawQuery := `select a.user_id, b.full_name, b.gender, b.profile_picture, a.current_index, a.is_done from group_members a inner join public.user b on a.user_id = b.id where a.group_id = 3 order by a.current_index asc`

	err := repo.DbCon.Raw(rawQuery, groupId).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

// ListGroupMemberWithNameAndOffsetLimit func
func (repo *postgreGroupRepository) ListGroupMemberWithNameAndOffsetLimit(groupID int, offset int, limit int) ([]models.ResponseGroupMemberList, error) {
	var data []models.ResponseGroupMemberList

	rawQuery := `select a.id, a.user_id, b.full_name, b.gender, b.profile_picture, a.current_index, a.is_done from group_members a inner join public.user b on a.user_id = b.id where a.group_id = ? order by a.current_index asc offset ? limit ?`

	err := repo.DbCon.Raw(rawQuery, groupID, offset, limit).
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *postgreGroupRepository) ListGroupMemberByStatusAndPaging(groupId int, offset int, limit int, status int) ([]models.ResponseGroupMemberList, error) {

	var data []models.ResponseGroupMemberList
	err := repo.DbCon.Raw("select a.id, a.user_id, b.full_name, b.gender, b.profile_picture, a.current_index, a.is_done from group_members a inner join public.user b on a.user_id = b.id where a.group_id = ? and a.is_done = ? order by a.current_index asc offset ? limit ?", groupId, status, offset, limit).Scan(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *postgreGroupRepository) TotalGroupMemberByStatus(groupID int, status int) (int, error) {
	var count int
	err := repo.DbCon.Table("group_members").Where("group_id = ? and is_done = ? ", groupID, status).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *postgreGroupRepository) TotalGroupMember(groupID int) (int, error) {
	var count int
	err := repo.DbCon.Table("group_members").Where("group_id = ?", groupID).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *postgreGroupRepository) GetAllGroupMembersWithType() ([]models.GroupMemberWithType, error) {
	var data []models.GroupMemberWithType
	err := repo.DbCon.Raw("select a.id, a.group_id, b.group_type, a.user_id, a.current_index, a.is_done, " +
		"a.created_at, a.updated_at from group_members a inner join groups b on a.group_id = b.id").Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *postgreGroupRepository) GetAllGroupMembersReadNotIsDone() ([]models.GroupMemberReadNotIsDone, error) {

	var data []models.GroupMemberReadNotIsDone

	query := `select a.user_id, a.is_done, b.full_name, c.token_firebase from public.group_members a
	left join public.user b on a.user_id = b.id
	left join public.user_token c on a.user_id = c.user_id
	where is_done = 0 group by a.user_id, a.is_done, b.full_name, c.token_firebase`

	//query := `select * from public.fn_get_all_group_member_readnot_isdone()`

	err := repo.DbCon.Raw(query).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *postgreGroupRepository) UpdateGroupMemberReadingStatusToZero() error {
	err := repo.DbCon.Table("group_members").Update("is_done", 0).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreGroupRepository) InsertHistory(history models.History) error {
	var newhistory models.History
	err := repo.DbCon.Table("tbl_history").Create(&history).Scan(&newhistory).Error
	if err != nil {
		return err
	}
	return err
}

func (repo *postgreGroupRepository) UpdateGroupMemberReadingIndex(ID int, index int) error {

	sqlUpdate := fmt.Sprintf(`UPDATE group_members set current_index = $2, updated_at = $3, is_done = $4 WHERE id = $1`)
	err := repo.DbCon.Exec(sqlUpdate, ID, index, time.Now(), 0).Error

	//err := repo.DbCon.Table("group_members").Where("id = ? ", ID).Update("current_index", index).Error
	if err != nil {
		return err
	}
	return nil
}
