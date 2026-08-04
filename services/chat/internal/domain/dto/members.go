package dto

type GetMembersDBResponse struct {
	Members []*MemberDB
	IDs     []int64
}
