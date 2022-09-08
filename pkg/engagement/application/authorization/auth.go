package authorization

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/profileutils"
)

// CheckPemissions is used to check whether the permissions of a subject are set
func CheckPemissions(subject string, input profileutils.PermissionInput) (bool, error) {
	enforcer := &casbin.Enforcer{}
	ok, err := enforcer.Enforce(subject, input.Resource, input.Action)
	if err != nil {
		return false, fmt.Errorf("unable to check permissions %w", err)
	}
	if ok {
		return true, nil
	}
	return false, nil
}

// CheckAuthorization is used to check the user permissions
func CheckAuthorization(subject string, permission profileutils.PermissionInput) (bool, error) {
	isAuthorized, err := CheckPemissions(subject, permission)
	if err != nil {
		return false, fmt.Errorf("internal server error: can't authorize user: %w", err)
	}

	if !isAuthorized {
		return false, nil
	}

	return true, nil
}

// IsAuthorized checks if the subject identified by their email has permission to access the
// specified resource
// currently only known internal anonymous users and external API Integrations emails are checked, internal and default logged in users
// have access by default.
// for subjects identified by their phone number normalize the phone and omit the first (+) character
func IsAuthorized(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error) {
	if user.PhoneNumber != "" && converterandformatter.StringSliceContains(profileutils.AuthorizedPhones, user.PhoneNumber) {
		return CheckAuthorization(user.PhoneNumber[1:], permission)
	}
	if user.Email != "" && converterandformatter.StringSliceContains(profileutils.AuthorizedEmails, user.Email) {
		return CheckAuthorization(user.Email, permission)

	}
	return false, nil
}
