package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/content"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/spf13/viper"
)

type postgreContentRepository struct {
	DbConn *gorm.DB
}

// NewContentRepository return
func NewContentRepository(DbConn *gorm.DB) content.Repository {
	return &postgreContentRepository{DbConn}
}

func (repo *postgreContentRepository) AddAyatQuran(surah []models.AyatQuran) error {
	for _, s := range surah {
		err := repo.DbConn.Table("quran").Create(&s).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *postgreContentRepository) GetAyatByCategory(category string, index int) ([]models.ResponseQuran, error) {

	var ayat []models.ResponseQuran
	// var NActivity []models.UserActivity
	// var CountActivityUserRead int

	query := fmt.Sprintf("select a.arabic, a.latin, a.translation, a.audio, a.image, a.ayat_sec, b.surah_name, b.surah_type, b.number, b.ayat from quran a "+
		"inner join quran_surah b on a.surah_id = b.id where %s = %d order by a.id asc", category, index)

	err := repo.DbConn.Raw(query).Scan(&ayat).Error
	if err != nil {
		return nil, err
	}

	// CheckUserReadActivity := repo.DbCon.Table("user_activity").Where("status_user_action = ? and user_id = ? ", 4, status).
	// 	Count(&count).Error

	// if err != nil {
	// 	return 0, err
	// }

	// if CountCheckUserActivityStatusFinishReadGroup == 0 {

	// }

	// insUserActivity := models.UserActivity{
	// 	UserId:           params.UserID,
	// 	IdGroupMember:    0,
	// 	GroupId:          0,
	// 	GroupType:        strings.Title(category),
	// 	ContentIndex:     index,
	// 	StatusUserAction: 4,
	// 	Description:      "Reading Surah",
	// 	CreatedAt:        now,
	// }

	// if err := tx.Table("user_activity").Create(&insUserActivity).Scan(&NActivity).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	return ayat, nil
}

// GetSurahByIdSurah Function
func (repo *postgreContentRepository) GetSurahByIDSurah(SurahID int) ([]models.ResponseQuran, error) {

	fmt.Print(SurahID)

	var ayat []models.ResponseQuran

	err := repo.DbConn.Raw(`select a.arabic, a.latin, a.translation, a.audio, a.image, a.ayat_sec, b.surah_name, b.surah_type,b.number, b.ayat from quran a inner join quran_surah b on a.surah_id = b.id where a.surah_id = ? order by a.id asc`, SurahID).Scan(&ayat).Error
	if err != nil {
		return nil, err
	}

	return ayat, nil

}
func (repo *postgreContentRepository) GetQuranByCategoryAndPaging(category string, index int, offset int, limit int) ([]models.ResponseQuran, error) {

	query := fmt.Sprintf("select a.arabic, a.latin, a.translation, a.audio, a.image, a.ayat_sec, b.surah_name, b.surah_type, b.number, b.ayat from quran a "+
		"inner join quran_surah b on a.surah_id = b.id where %s = %d order by a.id asc offset %d limit %d", category, index, offset, limit)
	var ayat []models.ResponseQuran

	err := repo.DbConn.Raw(query).Scan(&ayat).Error
	if err != nil {
		return nil, err
	}
	return ayat, nil
}

func (repo *postgreContentRepository) SaveUserReadingActivity(params models.RequestSaveReadingQuran) ([]models.UserActivity, error) {

	tx := repo.DbConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	now := time.Now()
	CurrentDate := now.Format("2006-01-02")

	var NActivity []models.UserActivity
	//var CheckUserActivityStatusFinishReadGroup []models.UserActivity
	var CountCheckUserActivityStatusFinishReadGroup int

	insUserActivity := models.UserActivity{
		UserId:           params.UserID,
		IdGroupMember:    params.ID,
		GroupId:          params.GroupID,
		GroupType:        strings.Title(params.GroupType),
		ContentIndex:     params.ContentIndex,
		StatusUserAction: 1,
		Description:      "Finish Reading Group",
		CreatedAt:        now,
	}

	if err := tx.Table("user_activity").Create(&insUserActivity).Scan(&NActivity).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Table("group_members").Where("id = ?", params.ID).Update("is_done", 1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if strings.Title(params.GroupType) == "Ayat" {

		fmt.Print("Masuk Ke Ayat")
		fmt.Print("\n")

		if err := tx.Table("user_activity").Where("status_user_action = ? and group_type = ? and group_id = ? and DATE(created_at) = ?", 1, "Ayat", params.GroupID, CurrentDate).Count(&CountCheckUserActivityStatusFinishReadGroup).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if CountCheckUserActivityStatusFinishReadGroup == 6236 {

			var insertedKhatam models.GroupKhatam

			GroupKhatam := models.GroupKhatam{
				IDRefActivity: NActivity[0].Id,
				GroupID:       params.GroupID,
				GroupType:     "Ayat",
				CreatedAt:     now,
			}

			if err := repo.DbConn.Table("group_khatam").Create(&GroupKhatam).Scan(&insertedKhatam).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			var NotifUserKhatam []models.ResponseNotifUserKhatam

			err := repo.DbConn.Raw(`SELECT UA.group_id, UA.group_type, UA.description, USR.full_name, USR_TKN.token_firebase, G.no_group_index as no_group FROM public.user_activity UA LEFT JOIN public.user USR ON UA.user_id = USR.id LEFT JOIN public.user_token USR_TKN ON UA.user_id = USR_TKN.user_id LEFT JOIN public.groups G ON UA.group_id = G.id WHERE UA.status_user_action = ? AND UA.group_type = ? AND UA.group_id = ? GROUP BY UA.group_id, UA.group_type, UA.description, USR.full_name, USR_TKN.token_firebase, G.no_group_index`, 1, "Ayat", params.GroupID).Scan(&NotifUserKhatam).Error
			if err != nil {
				return nil, err
			}

			if len(NotifUserKhatam) > 0 {

				GetToken := []string{}

				URLFcm := viper.GetString(`fcm_key.url`)
				serverKey := viper.GetString(`fcm_key.server_key`)

				for _, item := range NotifUserKhatam {
					GetToken = append(GetToken, ""+item.TokenFirebase+"")
				}

				msgBody := fmt.Sprintf("Group %s ID %d anda sudah melakukan khatam quran untuk hari ini. Yuk bersama sama khatam lagi di hari esok.", params.GroupType, NotifUserKhatam[0].NoGroup)

				notification := models.NotificationSendFcm{
					Title:    "Selamat ..",
					Priority: "high",
					Body:     msgBody,
				}

				PostJSON := models.RequestSendFcm{
					RegistrationIds: GetToken,
					Notification:    notification,
				}

				payloadByte, _ := json.Marshal(PostJSON)

				var payload = bytes.NewReader(payloadByte)
				req, errSend := http.NewRequest("POST", URLFcm, payload)
				req.Header.Set("Authorization", "key="+serverKey)
				req.Header.Set("Content-Type", "application/json")

				if errSend != nil {
					fmt.Println(errSend)
				}

				client := new(http.Client)
				response, errRes := client.Do(req)

				if errRes != nil {
					panic(errRes)
				}

				defer response.Body.Close()

				if response != nil {
					//fmt.Println("response Status:", response.Status)
					//fmt.Println("response Headers:", response.Header)
					body, _ := ioutil.ReadAll(response.Body)
					fmt.Println("response Body:", string(body))
				} else {
					//fmt.Println("response :", response)
				}
			}
		}

	} else if strings.Title(params.GroupType) == "Halaman" {

		if err := tx.Table("user_activity").Where("status_user_action = ? and group_type = ? and group_id = ? and DATE(created_at) = ? ", 1, "Halaman", params.GroupID, CurrentDate).Count(&CountCheckUserActivityStatusFinishReadGroup).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if CountCheckUserActivityStatusFinishReadGroup == 604 {

			var insertedKhatam models.GroupKhatam

			GroupKhatam := models.GroupKhatam{
				IDRefActivity: NActivity[0].Id,
				GroupID:       params.GroupID,
				GroupType:     "Halaman",
				CreatedAt:     now,
			}

			if err := repo.DbConn.Table("group_khatam").Create(&GroupKhatam).Scan(&insertedKhatam).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			var NotifUserKhatam []models.ResponseNotifUserKhatam

			err := repo.DbConn.Raw(`SELECT UA.group_id, UA.group_type, UA.description, USR.full_name, USR_TKN.token_firebase, G.no_group_index as no_group FROM public.user_activity UA LEFT JOIN public.user USR ON UA.user_id = USR.id LEFT JOIN public.user_token USR_TKN ON UA.user_id = USR_TKN.user_id LEFT JOIN public.groups G ON UA.group_id = G.id WHERE UA.status_user_action = ? AND UA.group_type = ? AND UA.group_id = ? GROUP BY UA.group_id, UA.group_type, UA.description, USR.full_name, USR_TKN.token_firebase, G.no_group_index`, 1, "Halaman", params.GroupID).Scan(&NotifUserKhatam).Error
			if err != nil {
				return nil, err
			}

			if len(NotifUserKhatam) > 0 {

				GetToken := []string{}

				URLFcm := viper.GetString(`fcm_key.url`)
				serverKey := viper.GetString(`fcm_key.server_key`)

				for _, item := range NotifUserKhatam {
					GetToken = append(GetToken, ""+item.TokenFirebase+"")
				}

				msgBody := fmt.Sprintf("Group %s ID %d anda sudah melakukan khatam quran untuk hari ini. Yuk bersama sama khatam lagi di hari esok.", params.GroupType, NotifUserKhatam[0].NoGroup)

				notification := models.NotificationSendFcm{
					Title:    "Selamat ..",
					Priority: "high",
					Body:     msgBody,
				}

				PostJSON := models.RequestSendFcm{
					RegistrationIds: GetToken,
					Notification:    notification,
				}

				payloadByte, _ := json.Marshal(PostJSON)

				var payload = bytes.NewReader(payloadByte)
				req, errSend := http.NewRequest("POST", URLFcm, payload)
				req.Header.Set("Authorization", "key="+serverKey)
				req.Header.Set("Content-Type", "application/json")

				if errSend != nil {
					fmt.Println(errSend)
				}

				client := new(http.Client)
				response, errRes := client.Do(req)

				if errRes != nil {
					panic(errRes)
				}

				defer response.Body.Close()

				if response != nil {
					//fmt.Println("response Status:", response.Status)
					//fmt.Println("response Headers:", response.Header)
					body, _ := ioutil.ReadAll(response.Body)
					fmt.Println("response Body:", string(body))
				} else {
					//fmt.Println("response :", response)
				}
			}
		}

	} else if strings.Title(params.GroupType) == "Juz" {

		if err := tx.Table("user_activity").Where("status_user_action = ? and group_type = ? and group_id = ? and DATE(created_at) = ?", 1, "Juz", params.GroupID, CurrentDate).Count(&CountCheckUserActivityStatusFinishReadGroup).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if CountCheckUserActivityStatusFinishReadGroup == 30 {

			var insertedKhatam models.GroupKhatam

			GroupKhatam := models.GroupKhatam{
				IDRefActivity: NActivity[0].Id,
				GroupID:       params.GroupID,
				GroupType:     "Juz",
				CreatedAt:     now,
			}

			if err := repo.DbConn.Table("group_khatam").Create(&GroupKhatam).Scan(&insertedKhatam).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			var NotifUserKhatam []models.ResponseNotifUserKhatam

			err := repo.DbConn.Raw(`SELECT UA.group_id, UA.group_type, UA.description, USR.full_name, USR_TKN.token_firebase, G.no_group_index as no_group FROM public.user_activity UA LEFT JOIN public.user USR ON UA.user_id = USR.id LEFT JOIN public.user_token USR_TKN ON UA.user_id = USR_TKN.user_id LEFT JOIN public.groups G ON UA.group_id = G.id WHERE UA.status_user_action = ? AND UA.group_type = ? AND UA.group_id = ? GROUP BY UA.group_id, UA.group_type, UA.description, USR.full_name, USR_TKN.token_firebase, G.no_group_index`, 1, "Juz", params.GroupID).Scan(&NotifUserKhatam).Error
			if err != nil {
				return nil, err
			}

			if len(NotifUserKhatam) > 0 {

				GetToken := []string{}

				URLFcm := viper.GetString(`fcm_key.url`)
				serverKey := viper.GetString(`fcm_key.server_key`)

				for _, item := range NotifUserKhatam {
					GetToken = append(GetToken, ""+item.TokenFirebase+"")
				}

				msgBody := fmt.Sprintf("Group %s ID %d anda sudah melakukan khatam quran untuk hari ini. Yuk bersama sama khatam lagi di hari esok.", params.GroupType, NotifUserKhatam[0].NoGroup)

				notification := models.NotificationSendFcm{
					Title:    "Selamat ..",
					Priority: "high",
					Body:     msgBody,
				}

				PostJSON := models.RequestSendFcm{
					RegistrationIds: GetToken,
					Notification:    notification,
				}

				payloadByte, _ := json.Marshal(PostJSON)

				var payload = bytes.NewReader(payloadByte)
				req, errSend := http.NewRequest("POST", URLFcm, payload)
				req.Header.Set("Authorization", "key="+serverKey)
				req.Header.Set("Content-Type", "application/json")

				if errSend != nil {
					fmt.Println(errSend)
				}

				client := new(http.Client)
				response, errRes := client.Do(req)

				if errRes != nil {
					panic(errRes)
				}

				defer response.Body.Close()

				if response != nil {
					//fmt.Println("response Status:", response.Status)
					//fmt.Println("response Headers:", response.Header)
					body, _ := ioutil.ReadAll(response.Body)
					fmt.Println("response Body:", string(body))
				} else {
					//fmt.Println("response :", response)
				}
			}
		}
	}

	fmt.Print("CountCheckUserActivityStatusFinishReadGroup => ", CountCheckUserActivityStatusFinishReadGroup)
	fmt.Print("\n")

	return NActivity, tx.Commit().Error

}

func (repo *postgreContentRepository) UpdateQuranPageById(id int, page int) error {
	err := repo.DbConn.Table("quran").Where("id = ?", id).Update("page", page).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreContentRepository) GetQuotes() ([]models.Quote, error) {
	var quotes []models.Quote
	//err := repo.DbConn.Table("quotes").Where("is_active = ?", 1).Find(&quotes).Error
	err := repo.DbConn.Raw(`select * from quotes where is_active = 1 and deleted_at is null order by created_at desc`).Scan(&quotes).Error

	if err != nil {
		return nil, err
	}
	return quotes, nil
}

func (repo *postgreContentRepository) ListUserGroup(userId int) ([]models.ResponseGroupList, error) {
	var groups []models.ResponseGroupList
	err := repo.DbConn.Table("group_members").Select("b.id, b.group_name, b.group_type, b.current_member, group_members.current_index, CASE WHEN group_members.is_done = 1 THEN true ELSE false END AS is_reading").
		Joins("inner join groups b on group_members.group_id = b.id").Where("group_members.user_id = ? and is_active = 1", userId).Scan(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo *postgreContentRepository) GetAllQuranList() ([]models.ResponseAllQuran, error) {

	var SurahSurah []models.ResponseAllQuran

	query := fmt.Sprintf("select * from quran_surah order by id,number asc")

	err := repo.DbConn.Raw(query).Scan(&SurahSurah).Error
	if err != nil {
		return nil, err
	}
	return SurahSurah, nil

}

func (repo *postgreContentRepository) GetSplashScreenList() ([]models.SplashScreen, error) {
	var splashScreen []models.SplashScreen

	query := fmt.Sprintf("select * from splash_screen where is_active = 1 order by position asc")

	err := repo.DbConn.Raw(query).Find(&splashScreen).Error

	if err != nil {
		return nil, err
	}
	return splashScreen, nil
}

func (repo *postgreContentRepository) GetTotalKhatamAllGroup() ([]models.ResponseTotalKhatam, error) {

	var Khatam []models.ResponseTotalKhatam

	errCheckKhatam := repo.DbConn.Raw("select count(t1.id_ref_activity) as count, t2.group_type from group_khatam t1 left join user_activity t2 on t1.id_ref_activity = t2.id group by t2.group_type").Scan(&Khatam).Error

	if errCheckKhatam != nil {
		return nil, errCheckKhatam
	}

	return Khatam, nil
}

func (repo *postgreContentRepository) GetListIbukota() ([]models.ResponseListIbuKota, error) {
	var ScanData []models.ResponseListIbuKota

	var resData []models.ResponseListIbuKota

	err := repo.DbConn.Raw("select prop_key as id,ibukota as name from dukcapil_propinsi order by ibukota asc").Scan(&ScanData).Error

	if err != nil {
		return nil, err
	}

	for _, d := range ScanData {

		var url = "https://api.banghasan.com/sholat/format/json/kota/nama/" + strings.ToLower(d.Name)
		// fmt.Print(url)
		// fmt.Print("\n")
		resp, err := http.Get(url)
		if err != nil {
			fmt.Print("ada err ketika get data kota", err)
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}
		KotaAPIBangsa, err := getKotaAPIBangSa([]byte(body))

		if len(KotaAPIBangsa.Kota) > 1 {

			for _, m := range KotaAPIBangsa.Kota {

				resData = append(resData, models.ResponseListIbuKota{
					ID:   m.ID,
					Name: m.Nama,
				})

			}

		} else {
			resData = append(resData, models.ResponseListIbuKota{
				ID:   KotaAPIBangsa.Kota[0].ID,
				Name: KotaAPIBangsa.Kota[0].Nama,
			})
		}

	}

	return resData, err
}

// getKotaApiBangSa func
func getKotaAPIBangSa(body []byte) (*models.ResponseGetKotaFromAPIBangsa, error) {
	var s = new(models.ResponseGetKotaFromAPIBangsa)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

// GetListProvinces
func (repo *postgreContentRepository) GetListProvinces() ([]models.ResponseListProvinces, error) {

	var ScanData []models.ResponseListProvinces

	err := repo.DbConn.Table("provinces").Scan(&ScanData).Error

	if err != nil {
		return nil, err
	}

	return ScanData, err
}

// GetListCities
func (repo *postgreContentRepository) GetListCities(provinceID int) ([]models.ResponseListCities, error) {

	var ScanData []models.ResponseListCities

	fmt.Print("provinceID => ", provinceID)

	if provinceID != 0 {
		err := repo.DbConn.Table("cities").Where("province_id = ?", provinceID).Find(&ScanData).Error

		if err != nil {
			return nil, err
		}

		return ScanData, err

	} else {

		err := repo.DbConn.Table("cities").Scan(&ScanData).Error

		if err != nil {
			return nil, err
		}

		return ScanData, err

	}

}

func (repo *postgreContentRepository) GetTotalUserByProvince() ([]models.ResponseTotalUserByProvince, bool) {

	var ScanData []models.ResponseTotalUserByProvince

	Query := `select count(USR.id_province) as count, PROV.internation_id, PROV.name from public.user USR 
	RIGHT JOIN provinces PROV ON USR.id_province=PROV.id
	group by PROV.internation_id, PROV.name order by PROV.name asc`

	err := repo.DbConn.Raw(Query).Scan(&ScanData).Error

	if err != nil {
		return nil, false
	}

	return ScanData, true
}

func (repo *postgreContentRepository) GetListAllBoarding() ([]models.BoardingPage, error) {

	var BoardingPage []models.BoardingPage

	query := fmt.Sprintf("select * from boardingpage order by position asc")

	err := repo.DbConn.Raw(query).Find(&BoardingPage).Error

	if err != nil {
		return nil, err
	}
	return BoardingPage, nil
}

func (repo *postgreContentRepository) GetListBoardingPageIsActive() ([]models.BoardingPage, error) {

	var BoardingPage []models.BoardingPage

	query := fmt.Sprintf("select * from boardingpage where is_active = 1 order by position asc")

	err := repo.DbConn.Raw(query).Find(&BoardingPage).Error

	if err != nil {
		return nil, err
	}
	return BoardingPage, nil
}

func (repo *postgreContentRepository) AddBoardingPage(boardingpage models.BoardingPage) ([]models.BoardingPage, error) {

	var BoardingScan []models.BoardingPage

	err := repo.DbConn.Table("boardingpage").Create(&boardingpage).Scan(&BoardingScan).Error
	if err != nil {
		return BoardingScan, err
	}

	return BoardingScan, nil

}

func (repo *postgreContentRepository) UpdateBoardingPage(boardingpage models.BoardingPage) ([]models.BoardingPage, error) {

	var boardingpageScan []models.BoardingPage

	err := repo.DbConn.Table("boardingpage").Where("id = ? ", boardingpage.ID).Update(&boardingpage).Scan(&boardingpageScan).Error
	if err != nil {
		return boardingpageScan, err
	}
	return boardingpageScan, nil
}

func (repo *postgreContentRepository) DeactivateBoardingPage(paramID int) error {

	err := repo.DbConn.Table("boardingpage").Where("id = ? ", paramID).Update("is_active", 0).Error

	if err != nil {
		return err
	}

	return nil
}
