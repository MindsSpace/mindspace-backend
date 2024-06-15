package messages

const (
	MsgUserAuthenticateSuccess = "User authentication successful"
	MsgUserAuthenticateFailed  = "Failed to process user authentication request"
	MsgUserWrongPassword       = "Entered password is incorrect"

	MsgUsersFetchSuccess = "Users fetched successfully"
	MsgUsersFetchFailed  = "Failed to fetch users"
	MsgUserFetchSuccess  = "User fetched successfully"
	MsgUserFetchFailed   = "Failed to fetch user"

	MsgUserUpdateSuccess = "User update successful"
	MsgUserUpdateFailed  = "Failed to process user update request"

	MsgUserDeleteSuccess = "User delete successful"
	MsgUserDeleteFailed  = "Failed to process user delete request"

	MsgUserPointAddSuccess = "User point add successful"
	MsgUserPointAddFailed  = "Failed to process user point add request"

	MsgUserAvatarUpdateSuccess = "User avatar update successful"
	MsgUserAvatarUpdateFailed  = "Failed to process user avatar update request"

	MsgUserAvatarDeleteSuccess = "User avatar delete successful"
	MsgUserAvatarDeleteFailed  = "Failed to process user avatar delete request"
)
