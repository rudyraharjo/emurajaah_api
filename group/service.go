package group

import "github.com/rudyraharjo/emurojaah/models"

type Service interface {

	// main controller
	JoinGroup(userId int, category string) (*models.GroupMember, []models.Group, bool)

	JoinGroupNew(userId int, category string) (*models.GroupMember, []models.Group, bool, int)

	JoinBulkGroup(reqBulk models.RequestJoinBulkGroup) ([]models.Group, bool, int)

	JoinGroupByEmail(groupID int, email string, grouptype string) (*models.GroupMember, []models.Group, bool, int)

	GetAvailableGroupByType(category string) *models.Group
	ExitGroup(groupId int, userId int) error

	LeaveReadingGroup(reqParams models.RequestLeaveReading) error

	HandleUpdateMemberReadingIndex()

	HandleNotifUserReadIsNotDone()

	// create
	CreateGroup(category string, nogroupindex int) *models.Group

	// get list
	GetAllGroupMembersWithType() []models.GroupMemberWithType
	GetListUserGroup(userId int) ([]models.ResponsGroupeGroupMembersListWithNoUrut, error)

	GetListGroupMembersByUserIDAndGroupID(userID int, groupID int) int

	GetListGroupType(userID int, GroupType string, Offset int, Limit int) ([]models.ResponseListGroupTypeWithNoUrut, error)

	GetListAllGroup() ([]models.ResponseAllgroup, error)

	GetListGroupMemberWithName(groupId int) []models.ResponseGroupMemberList
	GetListGroupMemberWithNameAndOffsetLimit(groupId int, offset int, limit int) []models.ResponseGroupMemberList
	GetListMemberGroupByStatusAndPaging(groupId int, offset int, limit int, status int) []models.ResponseGroupMemberList

	GetAllReadNotIsDone() []models.GroupMemberReadNotIsDone

	// update

	//UpdateInactiveMember(groupId int, oldUserId int, newUserId int) (*models.GroupMember, error)
	SetGroupMemberAsInactive(groupId int, userId int) error

	UpdateGroupMemberReadingStatusToZero() error

	// counter
	TotalGroupMemberByStatus(groupId int, status int) int
	TotalGroupMember(groupId int) int

	// UpdateGroupMemberCounter(groupId int) bool
	// ReduceGroupMemberCounter(groupId int) error
}
