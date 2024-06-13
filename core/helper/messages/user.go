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

	MsgUserPictureUpdateSuccess = "User picture update successful"
	MsgUserPictureUpdateFailed  = "Failed to process user picture update request"

	MsgUserPictureDeleteSuccess = "User picture delete successful"
	MsgUserPictureDeleteFailed  = "Failed to process user picture delete request"
)
