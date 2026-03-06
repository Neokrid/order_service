package errors

import (
	"github.com/Neokrid/order_service/pkg/errors"
)

var (
	InvalidAuthorizationHeader = errors.NewUnauthorizedError("invalid authorization header", "invalid_authorization_header")
	InvalidTokenError          = errors.NewUnauthorizedError("invalid token", "invalid_token")

	UserNotFound           = errors.NewInvalidDataError("user not found", "user_not_found")
	IncorrectPassword      = errors.NewUnauthorizedError("incorrect password", "incorrect_password")
	ConfirmCodeAlreadySend = errors.NewInvalidDataError("confirm code already send", "confirm_code_already_send")
	ConfirmCodeNotExist    = errors.NewInvalidDataError("confirm code not exist", "confirm_code_not_exist")
	ConfirmCodeIncorrect   = errors.NewInvalidDataError("confirm code incorrect", "confirm_code_incorrect")

	TokenClaimsError = errors.NewInvalidDataError("bad token claims", "bad_token_claims")
	TokensDontMatch  = errors.NewInvalidDataError("tokens dont match", "tokens_dont_match")
	TokenDontExist   = errors.NewInvalidDataError("token dont exist", "token_dont_exist")

	NoNewPassword = errors.NewBadRequestError("no new password", "no_new_password")
	NotUnique     = errors.NewInvalidDataError("not unique", "not_unique")

	NoPermissionsRequest  = errors.NewBadRequestError("no permission to respond to this request", "no_permission_to_respond_to_this_request")
	InvalidStatus         = errors.NewBadRequestError("invalid status", "invalid_status")
	FriendRequestNotFound = errors.NewInvalidDataError("friend request not found", "friend_request_not_found")
)
