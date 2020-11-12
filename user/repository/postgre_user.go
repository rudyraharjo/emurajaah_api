package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/rudyraharjo/emurojaah/user"
)

type postgreUserRepository struct {
	DbConn *gorm.DB
}

// NewUserRepository DB
func NewUserRepository(DbConn *gorm.DB) user.Repository {
	return &postgreUserRepository{DbConn}
}

func (repo *postgreUserRepository) RetrieveCredentialByAlias(alias string) (*models.UserCredential, error) {

	var cred models.UserCredential

	errQry := repo.DbConn.Table("user_alias").Select("user_id, credential").Where("alias = ?", alias).First(&cred).Error
	if errQry != nil {
		return nil, errQry
	}

	return &cred, nil
}

func (repo *postgreUserRepository) IsUserExist(id int) (int, error) {
	var count int
	err := repo.DbConn.Table("user").Where("id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 2, nil
	}
	return 1, nil
}

func (repo *postgreUserRepository) IsAliasExist(alias string) (int, error) {
	var count int
	err := repo.DbConn.Table("user_alias").Where("alias = ?", alias).Count(&count).Error
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 2, nil
	}
	return 1, nil
}

func (repo *postgreUserRepository) AddUser(user models.UserBasicInfo) (int, error) {
	var insertedUser models.UserBasicInfo

	errInsert := repo.DbConn.Table("user").Create(&user).Scan(&insertedUser).Error
	if errInsert != nil {
		return 0, errInsert
	}

	return insertedUser.Id, nil
}

func (repo *postgreUserRepository) AddUserAlias(alias []models.UserAliasFull) error {
	err := repo.DbConn.Table("user_alias").Create(&alias).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreUserRepository) AddSingleUserAlias(alias models.UserAliasFull) error {
	err := repo.DbConn.Table("user_alias").Create(&alias).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreUserRepository) AddUserToken(token models.UserToken) error {
	err := repo.DbConn.Table("user_token").Create(&token).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreUserRepository) AddSurvey(survey models.UserSurvey) error {
	err := repo.DbConn.Table("user_survey").Create(&survey).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreUserRepository) GetUserIdByAlias(alias string) (int, error) {
	var userAlias models.UserAliasFull
	err := repo.DbConn.Table("user_alias").Where("alias = ?", alias).First(&userAlias).Error
	if err != nil {
		return 0, err
	}
	return userAlias.UserId, nil
}

func (repo *postgreUserRepository) GetUserBasicInfoById(userId int) (*models.ResponseUserBasicInfo, error) {
	var userBasic models.ResponseUserBasicInfo
	//err := repo.DbConn.Table("user").Where("id = ?", userId).Find(&userBasic).Error

	err := repo.DbConn.Raw("SELECT USR.id,USR.full_name,USR.gender,USR.birth_date,USR.birth_place,USR.profile_picture,USR.created_at,USR.updated_at,USR.id_city,USR.id_province,USR.address, initcap(CT.name) as city, initcap(PROV.name) as province FROM public.user USR LEFT JOIN cities CT ON USR.id_city = CT.id LEFT JOIN provinces PROV ON USR.id_province = PROV.id where USR.id = ?", userId).Scan(&userBasic).Error
	if err != nil {
		return nil, err
	}
	return &userBasic, nil
}

func (repo *postgreUserRepository) GetUserPointGroupByType(userId int) ([]models.UserReward, error) {
	var reward []models.UserReward

	// err := repo.DbConn.Raw("select group_type as type, sum(point) as point from user_activity where user_id = ? group by group_type", userId).Scan(&reward).Error
	// if err != nil {
	// 	return nil, err
	// }

	var rewardAyat []models.UserRewardAyat
	var rewardHalaman []models.UserRewardHalaman
	var rewardJuz []models.UserRewardJuz

	errAyat := repo.DbConn.Raw("SELECT ua.group_type as type, p.point, sum(p.point) point FROM user_activity ua LEFT JOIN point p ON ua.group_type = p.type where ua.user_id = ? AND ua.group_type = ? group by ua.group_type, p.point", userId, "Ayat").Scan(&rewardAyat).Error
	if errAyat != nil {
		return nil, errAyat
	}

	errHalaman := repo.DbConn.Raw("SELECT ua.group_type as type, p.point, sum(p.point) as user_point FROM user_activity ua LEFT JOIN point p ON ua.group_type = p.type where ua.user_id = ? AND ua.group_type = ? group by ua.group_type, p.point", userId, "Halaman").Scan(&rewardHalaman).Error
	if errHalaman != nil {
		return nil, errHalaman
	}

	errJuz := repo.DbConn.Raw("SELECT ua.group_type as type, p.point, sum(p.point) as user_point FROM user_activity ua LEFT JOIN point p ON ua.group_type = p.type where ua.user_id = ? AND ua.group_type = ? group by ua.group_type, p.point", userId, "Juz").Scan(&rewardJuz).Error
	if errJuz != nil {
		return nil, errJuz
	}

	if len(rewardAyat) > 0 {

		reward = append(reward, models.UserReward{
			Type:      rewardAyat[0].Type,
			Point:     rewardAyat[0].Point,
			UserPoint: rewardAyat[0].UserPoint,
		})

	} else {

		reward = append(reward, models.UserReward{
			Type:      "Ayat",
			Point:     1,
			UserPoint: 0,
		})
	}

	if len(rewardHalaman) > 0 {
		reward = append(reward, models.UserReward{
			Type:      rewardHalaman[0].Type,
			Point:     rewardHalaman[0].Point,
			UserPoint: rewardHalaman[0].UserPoint,
		})
	} else {
		reward = append(reward, models.UserReward{
			Type:      "Halaman",
			Point:     11,
			UserPoint: 0,
		})
	}

	if len(rewardJuz) > 0 {
		reward = append(reward, models.UserReward{
			Type:      rewardJuz[0].Type,
			Point:     rewardJuz[0].Point,
			UserPoint: rewardJuz[0].UserPoint,
		})
	} else {
		reward = append(reward, models.UserReward{
			Type:      "Juz",
			Point:     208,
			UserPoint: 0,
		})
	}

	// reward = append(reward, models.UserReward{
	// 	Type:      "Khatam Ayat",
	// 	Point:     6240,
	// 	UserPoint: 0,
	// })

	// reward = append(reward, models.UserReward{
	// 	Type:      "Khatam Halaman",
	// 	Point:     6240,
	// 	UserPoint: 0,
	// })

	// reward = append(reward, models.UserReward{
	// 	Type:      "Khatam Juz",
	// 	Point:     6240,
	// 	UserPoint: 0,
	// })

	return reward, nil
}

func (repo *postgreUserRepository) GetUserAliasById(userId int) ([]models.UserAlias, error) {
	var alias []models.UserAlias
	err := repo.DbConn.Table("user_alias").Where("user_id = ?", userId).Find(&alias).Error
	if err != nil {
		return nil, err
	}
	return alias, nil
}

func (repo *postgreUserRepository) GetPublicStatOfReadQuran() ([]models.ResponseTotalKhatam, error) {

	var Khatam []models.ResponseTotalKhatam

	errGroupKhatam := repo.DbConn.Raw("select count(group_type) as count, group_type from group_khatam group by group_type").Scan(&Khatam).Error

	if errGroupKhatam != nil {
		return nil, errGroupKhatam
	}
	return Khatam, nil
}

func (repo *postgreUserRepository) GetPublicStatOfReadQuranByUserID(userID int) ([]models.ResponseTotalKhatam, error) {

	var Khatam []models.ResponseTotalKhatam
	var CountUserKhatam []models.ResponseUserTotalKhatam
	var resKhatam []models.ResponseTotalKhatam

	// resKhatam = append(resKhatam, models.ResponseTotalKhatam{
	// 	Count:     d.Count,
	// 	GroupType: d.GroupType,
	// 	IsJoined:  true,
	// })

	errGroupKhatam := repo.DbConn.Raw("SELECT group_type,(select count(*) from group_khatam where group_type = groups.group_type) as count FROM groups group by group_type").Scan(&Khatam).Error

	if errGroupKhatam != nil {
		return nil, errGroupKhatam
	}

	for _, d := range Khatam {

		if d.Count > 0 {

			errUserIsJoinedGroup := repo.DbConn.Raw("select A.group_id, COUNT(B.group_type) as count from group_khatam A left join user_activity B on A.group_id = B.group_id and A.group_type = B.group_type where B.user_id = ? and B.group_type = ? group by A.group_id", userID, d.GroupType).Scan(&CountUserKhatam).Error

			if errUserIsJoinedGroup != nil {
				return nil, errUserIsJoinedGroup
			}

			if len(CountUserKhatam) != 0 {
				resKhatam = append(resKhatam, models.ResponseTotalKhatam{
					Count:     d.Count,
					GroupType: d.GroupType,
					IsJoined:  true,
				})
			} else {
				resKhatam = append(resKhatam, models.ResponseTotalKhatam{
					Count:     d.Count,
					GroupType: d.GroupType,
					IsJoined:  false,
				})
			}

		} else {

			resKhatam = append(resKhatam, models.ResponseTotalKhatam{
				Count:     d.Count,
				GroupType: d.GroupType,
				IsJoined:  false,
			})
		}

	}

	return resKhatam, nil
}

func (repo *postgreUserRepository) GetUserStatOfReadQuran(userId int) ([]models.ResponsePersonalReadStatus, error) {

	var data []models.ResponsePersonalReadStatus
	err := repo.DbConn.Raw(`SELECT group_type as type,(select count(*) from user_activity where group_type = groups.group_type and user_id = ? and status_user_action = ?) as count FROM groups group by group_type`, userId, 1).Scan(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *postgreUserRepository) UpdateProfile(userInfo models.UserBasicInfo) error {

	err := repo.DbConn.Table("user").Where("id = ?", userInfo.Id).Update(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreUserRepository) IsTokenExits(id int) (int, error) {

	var count int

	err := repo.DbConn.Table("user_token").Where("user_id = ?", id).Count(&count).Error

	//err := repo.DbConn.Raw(`select * from public.fn_get_all_user_token_is_exits(?,?)`, id, token).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil

}

func (repo *postgreUserRepository) CheckSurahIsDone(params models.ReadIsDone) (int, error) {

	var count int
	err := repo.DbConn.Table("group_members").Where("user_id = ? and group_id = ? and current_index = ? and is_done = 1", params.UserID, params.GroupID, params.Index).Count(&count).Error

	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, err
	}

	return count, err
}

func (repo *postgreUserRepository) DeleteTokenFirebase(idUser int) error {

	err := repo.DbConn.Table("user_token").Where("user_id = ?", idUser).Delete(&idUser).Error

	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreUserRepository) GetUserAdminByAlias(alias string) (*models.UserAdminAliasFull, error) {
	var userAdmin models.UserAdminAliasFull

	err := repo.DbConn.Table("user_admin").Where("alias = ? OR username = ?", alias, alias).Find(&userAdmin).Error

	if err != nil {
		return nil, err
	}

	return &userAdmin, nil
}

func (repo *postgreUserRepository) RetrieveCredentialAdminByAlias(alias string) (*models.UserCredential, error) {

	var cred models.UserCredential

	err := repo.DbConn.Table("user_admin").Select("user_id, credential").Where("alias = ? OR username = ?", alias, alias).First(&cred).Error
	if err != nil {
		return nil, err
	}

	return &cred, nil
}

func (repo *postgreUserRepository) GetListUserAdminByAlias() ([]models.UserAdminAliasFull, error) {

	var ListUserAdmin []models.UserAdminAliasFull

	err := repo.DbConn.Raw("select * From user_admin order by id desc").Scan(&ListUserAdmin).Error

	if err != nil {
		return nil, err
	}

	return ListUserAdmin, nil

}

func (repo *postgreUserRepository) GetUserMemberList() ([]models.UserMemberAliasFull, error) {
	var listusermember []models.UserMemberAliasFull

	// err := repo.DbConn.Raw("select t1.id, t1.alias as email, t1.created_at, t2.full_name, t2.address, t2.province from public.user_alias t1 left join public.user t2 on t1.user_id = t2.id order by id asc").Scan(&listusermember).Error

	err := repo.DbConn.Raw("select USR.id,USR.full_name as fullname,USR.created_at, USR_ALS.alias as email, CT.name as cityname,PROV.name as provname from public.user USR left join public.user_alias USR_ALS on USR.id = USR_ALS.user_id left join public.cities CT on USR.id_city = CT.id left join public.provinces PROV on USR.id_province = PROV.id").Scan(&listusermember).Error

	if err != nil {
		return nil, err
	}

	return listusermember, nil

}

// GetPointGroupByType
func (repo *postgreUserRepository) GetPointGroupByType() ([]models.UserReward, error) {
	var reward []models.UserReward
	err := repo.DbConn.Raw("select group_type as type, sum(point) as point from user_activity group by group_type").Scan(&reward).Error
	if err != nil {
		return nil, err
	}
	return reward, nil
}

func (repo *postgreUserRepository) GetTotalMember() (int, error) {
	var count int

	err := repo.DbConn.Table("user_alias").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *postgreUserRepository) GetListUserMember() (int, error) {
	var count int

	return count, nil
}

func (repo *postgreUserRepository) GetSurveyJuzzGrouping() []models.ResponseSurveyJuzzGrouping {
	var resSurveyJuz []models.ResponseSurveyJuzzGrouping
	var scanResponseSurveyJuz []models.ResponseSurveyJuz
	var tempJuzzIDAns []int
	//kembar := false
	errDatas := repo.DbConn.Raw("select * from user_survey order by id asc").Scan(&scanResponseSurveyJuz).Error
	if errDatas != nil {
		return nil
	}

	if len(scanResponseSurveyJuz) != 0 {
		for _, d := range scanResponseSurveyJuz {

			s := strings.Split(d.SurveyAnswer, ",")

			for i := 0; i < len(s); i++ {
				no, _ := strconv.Atoi(s[i])
				tempJuzzIDAns = append(tempJuzzIDAns, no)
			}

		}

		dupMap := dupCount(tempJuzzIDAns)

		for h := 1; h <= 30; h++ {
			resSurveyJuz = append(resSurveyJuz, models.ResponseSurveyJuzzGrouping{
				JuzID:    h,
				CountJuz: 0,
			})
		}
		// for JuzID, Count := range dupMap {
		// 	resSurveyJuz = append(resSurveyJuz, models.ResponseSurveyJuzzGrouping{
		// 		JuzID:    JuzID,
		// 		CountJuz: Count,
		// 	})
		// }

		for i := 0; i < len(resSurveyJuz); i++ {

			for JuzID, Count := range dupMap {
				if resSurveyJuz[i].JuzID == JuzID {
					resSurveyJuz[i].CountJuz = Count
				}
			}
		}

	}

	return resSurveyJuz

}

func dupCount(list []int) map[int]int {

	dupFreq := make(map[int]int)

	for _, item := range list {

		_, exist := dupFreq[item]

		if exist {
			dupFreq[item]++ // Tambah 1 jika ada
		} else {
			dupFreq[item] = 1 // start 1
		}
	}
	return dupFreq
}

func (repo *postgreUserRepository) GetNameProvinceFromUser() ([]models.User, error) {

	var ScanUSer []models.User
	var ScanCity []models.ResponseListCities

	err := repo.DbConn.Table("user").Scan(&ScanUSer).Error
	if err != nil {
		return ScanUSer, err
	}

	if len(ScanUSer) > 0 {

		for _, item := range ScanUSer {

			QueryLike := `select * from cities where name LIKE '%' || $1 || '%'`

			errCity := repo.DbConn.Raw(QueryLike, strings.ToUpper(item.Province)).Scan(&ScanCity).Error

			if errCity != nil {
				fmt.Print(errCity)
			}

			if len(ScanCity) > 0 {
				sqlUpdate := fmt.Sprintf(`UPDATE public.user set id_city = $2, id_province = $3 WHERE id = $1`)
				errUpdate := repo.DbConn.Exec(sqlUpdate, item.ID, ScanCity[0].ID, ScanCity[0].ProvinceID).Error

				if errUpdate != nil {
					fmt.Print("Ada Error saat update ", errUpdate)
				}
			}

		}

	}

	return ScanUSer, nil

}
