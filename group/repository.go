package group

import "github.com/rudyraharjo/emurojaah/models"

type Repository interface {

	//create
	CreateGroup(g models.Group) (int, error)

	AddGroupMember(groupmembers models.GroupMember, groupId int) (*models.GroupMember, error)
	AddGroupMemberNew(userId int, category string) (*models.GroupMember, int, error)

	AddGroupMemberBulk(reqParams models.RequestJoinBulkGroup) (int, error)

	AddGroupMemberByEmail(groupID int, email string, grouptype string) (*models.GroupMember, int, error)

	// GetNoGroupIndex
	GetNoGroupIndex(groupID int) ([]models.Group, error)

	// boolean
	IsMemberExist(groupId int, userId int) (int, error)
	IsUserHasGroupCategory(userId int, category string) (int, error)

	// list
	GetAllGroupMembersWithType() ([]models.GroupMemberWithType, error)
	ListAvailableGroupByType(category string) ([]models.Group, error)

	GetAvailableGroupByTypeNew(category string) (*models.Group, error)

	GetLastNoIndexGroup(category string) ([]models.Group, error)

	ListInactiveGroupMember(groupId int) ([]models.GroupMember, error)
	ListActiveGroupMember(groupId int) ([]models.GroupMember, error)
	//GetListGroupMembersByIDGroup(groupId int) ([]models.GroupMember, error)
	GetListGroupMembersByIDGroupAndCurrenIndex(groupId int, genNumberIndex int) (int, error)

	ListAllGroup() ([]models.ResponseAllgroup, error)
	ListAllGroupIsJoined(userID int) ([]models.ResponseGroupUserJoined, error)

	GetListGroupType(userID int, groupType string, Offset int, Limit int) ([]models.ResponseListGroupTypeWithNoUrut, error)

	ListUserGroup(userId int) ([]models.ResponsGroupeGroupMembersListWithNoUrut, error)
	ListGroupMemberWithName(groupId int) ([]models.ResponseGroupMemberList, error)
	ListGroupMemberWithNameAndOffsetLimit(groupId int, offset int, limit int) ([]models.ResponseGroupMemberList, error)
	ListGroupMemberByStatusAndPaging(groupId int, offset int, limit int, status int) ([]models.ResponseGroupMemberList, error)

	GetListGroupMembersByUserIDAndGroupID(userID int, groupID int) (int, error)
	GetAllGroupMembersReadNotIsDone() ([]models.GroupMemberReadNotIsDone, error)

	// update
	//UpdateInactiveMember(groupId int, oldUserId int, newUserId int) (*models.GroupMember, error)
	SetGroupMemberAsInactive(groupId int, userId int) error
	UpdateGroupMemberReadingStatusToZero() error
	UpdateGroupMemberReadingIndex(ID int, index int) error

	// counter
	TotalGroupMemberByStatus(groupID int, status int) (int, error)
	TotalGroupMember(groupID int) (int, error)

	ExitGroupAndInsertActivity(groupID int, userID int) error
	LeaveReadingGroupAndInsertActivity(params models.RequestLeaveReading) error

	InsertHistory(History models.History) error
}
