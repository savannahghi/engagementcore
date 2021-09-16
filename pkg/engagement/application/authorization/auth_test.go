package authorization

import (
	"testing"

	"github.com/savannahghi/profileutils"
)

func TestCheckPemissions(t *testing.T) {
	type args struct {
		subject string
		input   profileutils.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy: permission passed",
			args: args{
				subject: "view",
				input: profileutils.PermissionInput{
					Resource: "action_view",
					Action:   "view",
				},
			},
			wantErr: false,
		},
		{
			name: "happy: permission passed",
			args: args{
				subject: "delete",
				input: profileutils.PermissionInput{
					Resource: "delete_item",
					Action:   "delete",
				},
			},
			wantErr: false,
		},
		{
			name: "happy: permission passed",
			args: args{
				subject: "view",
				input: profileutils.PermissionInput{
					Resource: "resolve_item",
					Action:   "resolve",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckPemissions(tt.args.subject, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPemissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestCheckAuthorization(t *testing.T) {
	type args struct {
		subject string
		input   profileutils.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy: permission passed",
			args: args{
				subject: "can publish",
				input: profileutils.PermissionInput{
					Resource: "publish_item",
					Action:   "create",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckAuthorization(tt.args.subject, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPemissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestIsAuthorized(t *testing.T) {
	type args struct {
		user       *profileutils.UserInfo
		permission profileutils.PermissionInput
	}

	user1 := &profileutils.UserInfo{
		DisplayName: "",
		Email:       "timhealthcloud.co.ke",
		PhoneNumber: "+254714568338",
		PhotoURL:    "",
		ProviderID:  "",
		UID:         "",
	}
	user2 := &profileutils.UserInfo{}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy: permission passed",
			args: args{
				user: user1,
				permission: profileutils.PermissionInput{
					Resource: "publish_item",
					Action:   "create",
				},
			},
			wantErr: false,
		},

		{
			name: "sad: permission denied",
			args: args{
				user: user2,
				permission: profileutils.PermissionInput{
					Resource: "publish_item",
					Action:   "create",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := IsAuthorized(tt.args.user, tt.args.permission)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPemissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
