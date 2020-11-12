package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/group"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreGroupRepository
type postgreGroupAdmRepository struct {
	DbConn *gorm.DB
}

// NewAdmGroupRepository DB
func NewAdmGroupRepository(DbConn *gorm.DB) group.Repository {
	return &postgreGroupAdmRepository{DbConn}
}

func (repo *postgreGroupAdmRepository) DeleteDuplicateGroupMembers() (int, error) {
	var groupMembers []models.GroupMember
	var groups []models.Group

	tx := repo.DbConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return 0, err
	}

	if err := tx.Raw(`SELECT a.* FROM group_members a, group_members b WHERE a.id > b.id AND (a.group_id = b.group_id and a.current_index = b.current_index) ORDER BY a.current_index ASC`).Scan(&groupMembers).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(groupMembers) > 0 {

		for i := 0; i < len(groupMembers); i++ {

			if err := tx.Raw(`SELECT * FROM groups WHERE id = ?`, groupMembers[i].GroupID).Scan(&groups).Error; err != nil {
				tx.Rollback()
				return 0, err
			}

			if len(groups) > 0 {

				for J := 0; J < len(groups); J++ {
					// println(J, groups[J].Id)
					if err := tx.Table("groups").Where("id = ?", groups[J].Id).Update("current_member", gorm.Expr("current_member - ?", 1)).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
				}

			}
		}

		if err := tx.Raw(`DELETE FROM group_members a USING group_members b WHERE a.id < b.id AND (a.group_id = b.group_id and a.current_index = b.current_index)`).Scan(&groupMembers).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		println("Delete Duplicate & Berkurang Di Table groups")

	} else {
		return 2, nil
	}

	return 1, tx.Commit().Error
	//return groupMembers, nil
}

func (repo *postgreGroupAdmRepository) GetListGroups() ([]models.GroupWithStatus, error) {
	var groupsWithStatus []models.GroupWithStatus

	//query := fmt.Sprintf("select * from groups where group_type in('Ayat', 'Halaman', 'Juz') order by group_name, no_group_index asc")

	// Query := fmt.Sprintf("select G.id, G.no_group_index, G.group_type,G.group_name,G.max_member,G.current_member, count(CASE WHEN GM.is_done = 1 THEN 1 END) as selesai, count(CASE WHEN GM.is_done = 0 THEN 0 END) as belom from groups G LEFT JOIN group_members GM on G.id = GM.group_id GROUP BY G.id, G.no_group_index, G.group_type,G.group_name,G.max_member,G.current_member")

	Query := fmt.Sprintf(`select G.id, G.no_group_index, G.group_type,G.group_name, G.max_member, count(GM.group_id) as current_member,
	count(CASE WHEN GM.is_done = 1 THEN 1 END) as selesai, 
	count(CASE WHEN GM.is_done = 0 THEN 0 END) as belom,
	G.created_at from groups G 
	LEFT JOIN group_members GM on G.id = GM.group_id 
	GROUP BY G.id, G.current_member, GM.group_id, G.no_group_index, G.group_type,G.group_name order by G.id asc`)

	err := repo.DbConn.Raw(Query).Find(&groupsWithStatus).Error

	if err != nil {
		return nil, err
	}
	return groupsWithStatus, nil
}

func (repo *postgreGroupAdmRepository) GetListGroupMember(groupID int) ([]models.ResponseGroupMember, error) {
	var groupMembers []models.ResponseGroupMember

	query := fmt.Sprintf("SELECT gm.user_id, gm.current_index, gm.is_done, gm.group_type, gp.no_group_index, usr.full_name, usr.gender, q.juz_id, q.surah_id, q.ayat_sec, q.page, qs.surah_name, usr.profile_picture,usr_als.alias as email FROM public.group_members gm LEFT JOIN public.user usr ON gm.user_id=usr.id LEFT JOIN public.groups gp ON gm.group_id=gp.id LEFT JOIN public.user_alias usr_als ON gm.user_id=usr_als.user_id LEFT JOIN public.quran q ON gm.current_index=q.id LEFT JOIN public.quran_surah qs ON q.surah_id=qs.id WHERE gm.group_id = ? order by gm.current_index asc")

	err := repo.DbConn.Raw(query, groupID).Find(&groupMembers).Error

	if err != nil {
		return groupMembers, err
	}

	return groupMembers, err
}

func (repo *postgreGroupAdmRepository) GetLastNoIndexGroup(category string) ([]models.Group, error) {

	var groups []models.Group
	err := repo.DbConn.Raw("select * from groups where group_type = ? order by no_group_index desc limit 1", category).Scan(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo *postgreGroupAdmRepository) CreateGroup(g models.Group) ([]models.Group, error) {
	var newGroup []models.Group
	err := repo.DbConn.Table("groups").Create(&g).Scan(&newGroup).Error
	if err != nil {
		return newGroup, err
	}
	return newGroup, nil
}

func (repo *postgreGroupAdmRepository) GetUserBelomBaca(paramID int) ([]models.GroupMemberReadNotIsDone, error) {
	var Members []models.GroupMemberReadNotIsDone

	Query := `select a.user_id, a.is_done, b.full_name, c.token_firebase from public.group_members a
	left join public.user b on a.user_id = b.id
	left join public.user_token c on a.user_id = c.user_id
	where group_id = ? and is_done = ? group by a.user_id, a.is_done, b.full_name, c.token_firebase`

	err := repo.DbConn.Raw(Query, paramID, 0).Scan(&Members).Error

	if err != nil {
		return Members, err
	}
	return Members, nil
}

func (repo *postgreGroupAdmRepository) InsertHistory(history models.History) error {
	var newhistory models.History
	err := repo.DbConn.Table("tbl_history").Create(&history).Scan(&newhistory).Error
	if err != nil {
		return err
	}
	return err
}
