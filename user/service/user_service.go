package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/rudyraharjo/emurojaah/models"
	"github.com/rudyraharjo/emurojaah/user"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo user.Repository
}

func NewUSerService(ur user.Repository) user.Service {
	return &userService{ur}
}

//HashPassword - Encrypt password
func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash - validate encrypt password
func (s *userService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *userService) LoginUser(params models.RequestLogin) (bool, int) {

	storedCred, errQry := s.userRepo.RetrieveCredentialByAlias(strings.ToLower(params.Alias))
	if errQry != nil {
		return false, 0
	}

	isPassMatch := s.CheckPasswordHash(params.Password, storedCred.Credential)
	if !isPassMatch {
		return false, 0
	}

	return true, storedCred.UserID
}

func (s *userService) LoginUserWithGoogle(params models.RequestLogin) (bool, int) {

	storedCred, errQry := s.userRepo.RetrieveCredentialByAlias(strings.ToLower(params.Alias))
	if errQry != nil {
		return false, 0
	}

	return true, storedCred.UserID
}

func (s *userService) LoginUserAdmin(params models.RequestLogin) (bool, int) {

	storedCred, errQry := s.userRepo.RetrieveCredentialAdminByAlias(strings.ToLower(params.Alias))
	if errQry != nil {
		return false, 0
	}

	isPassMatch := s.CheckPasswordHash(params.Password, storedCred.Credential)
	if !isPassMatch {
		return false, 0
	}

	return true, storedCred.UserID

}

func (s *userService) CheckIsUserExist(id int) int {
	exCode, err := s.userRepo.IsUserExist(id)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return exCode
}

func (s *userService) CheckTokenFirebase(id int) int {

	exCode, err := s.userRepo.IsTokenExits(id)

	if err != nil {
		fmt.Println(err)
		return 0
	}
	return exCode
}

func (s *userService) CheckIsAliasExist(alias string) int {
	exCode, err := s.userRepo.IsAliasExist(alias)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return exCode
}

func (s *userService) RegisterUser(registerParams models.RequestRegister) (int, error) {

	now := time.Now()
	layout := "2006-01-02"
	dates, errT := time.Parse(layout, registerParams.BirthDate)

	if errT != nil {
		return 0, errT
	}

	encPass, errEnc := s.HashPassword(registerParams.Password)
	if errEnc != nil {
		return 0, errEnc
	}

	// ------- TEST SEND EMAIL -------- //

	// to := []string{"yapmak.id@gmail.com"}
	// cc := []string{"rraharjo.rudy@gmail.com", "try.kdigitech@gmail.com"}

	// subject := "Test Send Email"
	// message := "Pesan ini hanya test"

	// err := s.sendMail(to, cc, subject, message)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// log.Println("mail sent!")

	// ------- END TEST SEND EMAIL -------- //

	userInfo := models.UserBasicInfo{
		FullName:       registerParams.FullName,
		Gender:         registerParams.Gender,
		IDProvince:     registerParams.IDProvince,
		IDCity:         registerParams.IDCity,
		BirthDate:      dates,
		BirthPlace:     "",
		ProfilePicture: "",
		CreatedAt:      now,
	}

	userId, errInst := s.userRepo.AddUser(userInfo)
	if errInst != nil {
		return 0, errInst
	}

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 14
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	GenActivationCode := b.String()

	userAlias := models.UserAliasFull{
		UserId:         userId,
		AliasType:      "Email",
		Alias:          strings.ToLower(registerParams.Email),
		Credential:     encPass,
		CredentialType: "PASS",
		IsVerified:     0,
		ActivationCode: GenActivationCode,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	errIns := s.userRepo.AddSingleUserAlias(userAlias)
	if errIns != nil {
		return 0, errIns
	}

	userToken := models.UserToken{
		UserID:        userId,
		TokenFirebase: registerParams.TokenFirebase,
		CreatedAt:     now,
	}

	errTkn := s.userRepo.AddUserToken(userToken)

	if errTkn != nil {
		return 0, errTkn
	}

	return userId, nil
}

func (s *userService) sendMail(to []string, cc []string, sub string, msg string) error {

	ConfigSMTPHost := viper.GetString(`config_smtp.host`)
	ConfigSMTPPort := viper.GetInt(`config_smtp.port`)
	ConfigEmail := viper.GetString(`config_smtp.email`)
	ConfigPassword := viper.GetString(`config_smtp.password`)

	body := "From: " + ConfigEmail + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + sub + "\n\n" +
		msg

	auth := smtp.PlainAuth("", ConfigEmail, ConfigPassword, ConfigSMTPHost)
	smtpAddr := fmt.Sprintf("%s:%d", ConfigSMTPHost, ConfigSMTPPort)

	err := smtp.SendMail(smtpAddr, auth, ConfigEmail, append(to, cc...), []byte(body))

	if err != nil {
		return err
	}

	return nil

}

func (s *userService) AddTokenUser(id int, token string) (int, error) {

	now := time.Now()

	userToken := models.UserToken{
		UserID:        id,
		TokenFirebase: token,
		CreatedAt:     now,
	}

	errTkn := s.userRepo.AddUserToken(userToken)

	if errTkn != nil {
		return 0, errTkn
	}

	return id, nil

}

func (s *userService) CheckSurahReadingIsDone(params models.ReadIsDone) (int, error) {

	paramCheckReadSurah := models.ReadIsDone{
		UserID:  params.UserID,
		GroupID: params.GroupID,
		Index:   params.Index,
	}

	isDone, err := s.userRepo.CheckSurahIsDone(paramCheckReadSurah)

	if err != nil {
		fmt.Println(err)
		return isDone, err
	}

	return isDone, err
}

func (s *userService) DeleteTokenFirebase(idUser int) error {

	err := s.userRepo.DeleteTokenFirebase(idUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) AddJuzzSurvey(reqParams models.RequestJuzzSurvey) error {
	survey := models.UserSurvey{
		UserId:         reqParams.UserId,
		SurveyAnswer:   reqParams.Answer,
		SurveyType:     "hafal-juzz",
		SurveyQuestion: "Berapa jumlah juzz yang telah Anda hafal?",
		CreatedDate:    time.Now().Add(time.Hour * 7),
	}

	err := s.userRepo.AddSurvey(survey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *userService) SurveyJuzzGrouping() ([]models.ResponseSurveyJuzzGrouping, error) {

	data := s.userRepo.GetSurveyJuzzGrouping()

	if len(data) == 0 {
		return data, nil
	}
	return data, nil
}

func (s *userService) GetUserAdminInfoByAlias(alias string) (*models.UserAdminAliasFull, *models.ResponseUserBasicInfo, error) {

	userAdm, err := s.userRepo.GetUserAdminByAlias(alias)
	if err != nil {
		return nil, nil, err
	}

	usr, err2 := s.userRepo.GetUserBasicInfoById(userAdm.UserID)
	if err2 != nil {
		fmt.Println(err2)
		return nil, nil, err2
	}

	fmt.Print(usr)

	return userAdm, usr, nil
}

func (s *userService) GetUserAdminList() []models.UserAdminAliasFull {

	data, err := s.userRepo.GetListUserAdminByAlias()
	if err != nil {
		return nil
	}
	return data

}

func (s *userService) GetUserMemberList() []models.UserMemberAliasFull {

	data, err := s.userRepo.GetUserMemberList()

	if err != nil {
		return nil
	}

	return data
}

func (s *userService) GetUserBasicInfoByAlias(alias string) (*models.ResponseUserBasicInfo, error) {
	userId, err1 := s.userRepo.GetUserIdByAlias(alias)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}

	usr, err2 := s.userRepo.GetUserBasicInfoById(userId)
	if err2 != nil {
		fmt.Println(err2)
		return nil, err2
	}

	return usr, nil
}

func (s *userService) GetUserBasicInfoById(userId int) (*models.ResponseUserBasicInfo, error) {
	usr, err := s.userRepo.GetUserBasicInfoById(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return usr, nil
}

func (s *userService) GetPointReward() (int, []models.UserReward, error) {

	data, err := s.userRepo.GetPointGroupByType()
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	totalPoint := 0
	for _, d := range data {
		totalPoint += d.Point
	}

	return totalPoint, data, nil

}

func (s *userService) GetUserPointReward(userId int) (int, []models.UserReward, error) {
	//var reward []models.UserReward

	data, err := s.userRepo.GetUserPointGroupByType(userId)

	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	totalPoint := 0

	for _, d := range data {
		totalPoint += d.UserPoint
	}

	return totalPoint, data, nil
}

func (s *userService) GetUserAlias(userId int) []models.UserAlias {
	data, err := s.userRepo.GetUserAliasById(userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *userService) GetPublicGroupStatistic(userID int) []models.ResponseTotalKhatam {
	totalkhatam, err := s.userRepo.GetPublicStatOfReadQuranByUserID(userID)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return totalkhatam
}

func (s *userService) GetUserReadStatistic(userId int) []models.ResponsePersonalReadStatus {
	data, err := s.userRepo.GetUserStatOfReadQuran(userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *userService) UpdateUserProfile(userInfo models.RequestEditProfile) error {

	layout := "2006-01-02"
	dates, _ := time.Parse(layout, userInfo.BirthDate)

	usr := models.UserBasicInfo{
		Id:             userInfo.UserId,
		FullName:       userInfo.FullName,
		Gender:         userInfo.Gender,
		ProfilePicture: userInfo.ProfilePicture,
		Address:        userInfo.Address,
		IDProvince:     userInfo.IDProvince,
		IDCity:         userInfo.IDCity,
		BirthDate:      dates,
		BirthPlace:     userInfo.BirthPlace,
		UpdatedAt:      time.Now(),
	}
	err := s.userRepo.UpdateProfile(usr)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *userService) LoginWithRequest(params models.RequestLogin) (int, *models.ResponseLoginSuccess, *models.ResponseLoginFail) {
	payloadByte, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return 0, nil, nil
	}

	APIAddress := "http://localhost:8081/api/login"

	var payload = bytes.NewReader(payloadByte)
	req, errReq := http.NewRequest("POST", APIAddress, payload)
	req.Header.Add("Content-Type", "application/json")

	if errReq != nil {
		fmt.Println(errReq)
		return 0, nil, nil
	}

	client := new(http.Client)
	response, errRes := client.Do(req)

	if errRes != nil {
		fmt.Println(errRes)
		return 0, nil, nil
	}

	if response != nil {
		defer response.Body.Close()

		if response.StatusCode == 200 {
			var resp models.ResponseLoginSuccess
			errP := json.NewDecoder(response.Body).Decode(&resp)

			if errP != nil {
				fmt.Println(errP)
				return 0, nil, nil
			}

			return response.StatusCode, &resp, nil

		} else {
			var resp models.ResponseLoginFail
			errP := json.NewDecoder(response.Body).Decode(&resp)

			if errP != nil {
				fmt.Println(errP)
				return 0, nil, nil
			}

			return response.StatusCode, nil, &resp
		}

	}

	return 0, nil, nil
}

func (s *userService) GetTotalMember() (int, error) {

	data, err := s.userRepo.GetTotalMember()

	if err != nil {
		return 0, err
	}

	return data, nil
}
