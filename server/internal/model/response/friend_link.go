package response

import "server/internal/model/database"

type FriendLinkInfo struct {
	List  []database.FriendLink `json:"list"`
	Total int64                 `json:"total"`
}
