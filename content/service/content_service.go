package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/splashscreen"
	"github.com/rudyraharjo/emurojaah/admin/termcondition"

	"github.com/rudyraharjo/emurojaah/admin/banner"
	"github.com/rudyraharjo/emurojaah/content"
	"github.com/rudyraharjo/emurojaah/group"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/rudyraharjo/emurojaah/user"
)

type contentService struct {
	contentRepo       content.Repository
	userRepo          user.Repository
	groupRepo         group.Repository
	bannerRepo        banner.Repository
	splashscreenRepo  splashscreen.Repository
	termconditionRepo termcondition.Repository
}

// NewContentService func
func NewContentService(repo content.Repository,
	repoUser user.Repository,
	repoGroup group.Repository,
	bannerRepo banner.Repository,
	splashscreenRepo splashscreen.Repository,
	termconditionRepo termcondition.Repository,
) content.Service {
	return &contentService{repo, repoUser, repoGroup, bannerRepo, splashscreenRepo, termconditionRepo}
}

func (s *contentService) AddQuranSurah(reqParams models.AddQuranRequest) error {
	surahId := reqParams.AyatId
	datas := reqParams.Surah
	var finalSurah []models.AyatQuran
	now := time.Now()

	for _, data := range datas {
		no, _ := strconv.Atoi(data.Nomor)

		finalSurah = append(finalSurah, models.AyatQuran{
			Arabic:      data.Ar,
			Latin:       data.Tr,
			AyatSec:     no,
			Translation: data.Id,
			SurahId:     surahId,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	err := s.contentRepo.AddAyatQuran(finalSurah)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *contentService) GetQuranFromAPI() error {
	for i := 1; i < 115; i++ {
		apiAddress := fmt.Sprintf("https://al-quran-8d642.firebaseio.com/surat/%d.json?print=pretty", i)
		req, errReq := http.NewRequest("GET", apiAddress, nil)
		if errReq != nil {
			fmt.Println(errReq)
			return errReq
		}

		client := new(http.Client)
		res, errRes := client.Do(req)

		if errRes != nil {
			fmt.Println(errReq)
			return errRes
		}

		defer res.Body.Close()

		var surah []models.SurahRequest
		errP := json.NewDecoder(res.Body).Decode(&surah)

		if errP != nil {
			fmt.Println(errP)
			return errP
		}

		quran := models.AddQuranRequest{
			AyatId: i,
			Surah:  surah,
		}

		_ = s.AddQuranSurah(quran)
	}
	return nil
}

func (s *contentService) HandlerAutoUpdateQuranPage() {
	for i := 1; i < 605; i++ {
		err := s.GetQuranPageFromAPI(i)
		if err != nil {
			break
		}
	}
}

func (s *contentService) HandlerUpdateIDCityAndIDProv() {

	_, err := s.userRepo.GetNameProvinceFromUser()

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Print(data)

}

func (s *contentService) GetQuranPageFromAPI(index int) error {
	apiAddress := fmt.Sprintf("http://api.alquran.cloud/v1/page/%d/en.asad", index)
	fmt.Printf(apiAddress)
	req, errReq := http.NewRequest("GET", apiAddress, nil)
	if errReq != nil {
		fmt.Println(errReq)
		return errReq
	}

	client := new(http.Client)
	res, errRes := client.Do(req)

	if errRes != nil {
		fmt.Println(errReq)
		return errRes
	}

	defer res.Body.Close()

	var page models.ResponseFromApiPageQuran
	errP := json.NewDecoder(res.Body).Decode(&page)

	if errP != nil {
		fmt.Println(errP)
		return errP
	}

	for _, ayat := range page.Data.Ayahs {
		fmt.Printf("%d : %d \n", index, ayat.Number)
		errU := s.UpdateQuranPageById(ayat.Number, index)
		if errU != nil {
			fmt.Println(errU)
			return errU
		}
	}

	return nil
}

func (s *contentService) GetSurahByCategory(params models.RequestQuran) []models.ResponseQuran {

	col := ""
	switch strings.ToLower(params.Type) {
	case "juz":
		col = "juz_id"
	case "surat":
		col = "surah_id"
	case "halaman":
		col = "page"
	case "ayat":
		col = "a.id"
	}

	data, err := s.contentRepo.GetAyatByCategory(col, params.Index)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

func (s *contentService) GetQuranByCategoryWithPaging(params models.RequestQuranPaging) []models.ResponseQuran {
	col := ""
	switch strings.ToLower(params.Type) {
	case "juz":
		col = "juz_id"
	case "surat":
		col = "surah_id"
	case "halaman":
		col = "page"
	case "ayat":
		col = "a.id"
	}

	data, err := s.contentRepo.GetQuranByCategoryAndPaging(col, params.Index, params.Offset, params.Limit)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

func (s *contentService) UpdateQuranPageById(id int, page int) error {
	err := s.contentRepo.UpdateQuranPageById(id, page)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *contentService) GetQuoteList() []models.Quote {
	quotes, err := s.contentRepo.GetQuotes()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return quotes
}

func (s *contentService) GetListUserGroup(userId int) []models.ResponseGroupList {
	data, err := s.contentRepo.ListUserGroup(userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}
func (s *contentService) GetPublicGroupStatistic(userId int) []models.ResponseGroupList {
	data, err := s.contentRepo.ListUserGroup(userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *contentService) HandlerHomePageContentGlobalGroupStatus(userID int) []models.ResponseTotalKhatam {

	totalkhatam, err := s.userRepo.GetPublicStatOfReadQuran()

	if err != nil {
		fmt.Print(err)
	}

	return totalkhatam

}

func (s *contentService) HandlerHomePageContentBanner(userID int) []models.Banner {
	banners, err := s.bannerRepo.GetBannerListActive()

	if err != nil {
		fmt.Println(err)
		return banners
	}
	return banners
}

func (s *contentService) HandlerHomePageContentGroups(userID int) []models.ResponseGroupUserJoined {
	groups, err := s.groupRepo.ListAllGroupIsJoined(userID)
	if err != nil {
		fmt.Println(err)
		return groups
	}
	return groups
}

func (s *contentService) HandlerHomePageContentQuotes(userID int) []models.Quote {

	quotes, err := s.contentRepo.GetQuotes()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return quotes

}

func (s *contentService) HandlerHomePageContent(userID int) models.ResponseHomePageContent {

	// banners := s.GetBannerList()
	quotes := s.GetQuoteList()
	banners, _ := s.bannerRepo.GetBannerListActive()
	// groups := s.GetListUserGroup(userID)
	//groups, _ := s.groupRepo.ListAllGroup()
	groups, _ := s.groupRepo.ListAllGroupIsJoined(userID)
	totalkhatam, err := s.userRepo.GetPublicStatOfReadQuran()

	if err != nil {
		fmt.Print(err)
	}

	resp := models.ResponseHomePageContent{
		Banner:              banners,
		Groups:              groups,
		Quotes:              quotes,
		ResponseTotalKhatam: totalkhatam,
	}
	return resp
}

func (s *contentService) GetAllQuran() []models.ResponseAllQuran {

	all_surah, err := s.contentRepo.GetAllQuranList()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return all_surah

}

func (s *contentService) GetSurahByIDSurah(SurahID int) []models.ResponseQuran {

	data, err := s.contentRepo.GetSurahByIDSurah(SurahID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *contentService) GetTermCondition() []models.TermCondition {

	res, err := s.termconditionRepo.GetListTermConditionsIsActived()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res

}

func (s *contentService) GetSplashScreenList() []models.SplashScreen {

	res, err := s.splashscreenRepo.GetSplashScreenListActive()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (s *contentService) GetTotalKhatamAllGroup() ([]models.ResponseTotalKhatam, error) {

	data, err := s.contentRepo.GetTotalKhatamAllGroup()
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *contentService) SaveUserReadingActivity(params models.RequestSaveReadingQuran) error {

	data, err := s.contentRepo.SaveUserReadingActivity(params)

	fmt.Print(data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *contentService) GetListIbukota() []models.ResponseListIbuKota {
	data, err := s.contentRepo.GetListIbukota()

	if err != nil {
		return data
	}

	return data
}

func (s *contentService) GetListProvinces() []models.ResponseListProvinces {

	data, err := s.contentRepo.GetListProvinces()

	if err != nil {
		return data
	}

	return data
}

func (s *contentService) GetListCities(provinceID int) []models.ResponseListCities {

	data, err := s.contentRepo.GetListCities(provinceID)

	if err != nil {
		return data
	}

	return data

}

func (s *contentService) GetTotalUserByProvince() ([]models.ResponseTotalUserByProvince, bool) {

	data, IsSuccess := s.contentRepo.GetTotalUserByProvince()

	if !IsSuccess {
		return data, false
	}

	return data, true

}

func (s *contentService) GetListAllBoarding() []models.BoardingPage {
	res, err := s.contentRepo.GetListAllBoarding()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (s *contentService) GetListBoardingPageIsActive() []models.BoardingPage {
	res, err := s.contentRepo.GetListBoardingPageIsActive()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (s *contentService) AddBoardingPage(boardingpage models.BoardingPage) ([]models.BoardingPage, error) {

	now := time.Now()

	addBoardingPage := models.BoardingPage{
		Title:       boardingpage.Title,
		Description: boardingpage.Description,
		ImageURL:    boardingpage.ImageURL,
		Position:    boardingpage.Position,
		IsActive:    1,
		CreatedAt:   now,
	}

	data, err := s.contentRepo.AddBoardingPage(addBoardingPage)

	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

func (s *contentService) UpdateBoardingPage(boardingpage models.BoardingPage) ([]models.BoardingPage, error) {

	now := time.Now()

	modelUpdateBoardingPage := models.BoardingPage{
		ID:          boardingpage.ID,
		Title:       boardingpage.Title,
		Description: boardingpage.Description,
		ImageURL:    boardingpage.ImageURL,
		Position:    boardingpage.Position,
		IsActive:    1,
		UpdatedAt:   now,
	}

	data, err := s.contentRepo.UpdateBoardingPage(modelUpdateBoardingPage)

	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

func (s *contentService) DeactivateBoardingPage(paramID int) error {
	err := s.contentRepo.DeactivateBoardingPage(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
