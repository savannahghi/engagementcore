package rest

import (
	"net/http"
	"testing"
)

func Test_getplayMP4QueryParam(t *testing.T) {

	type args struct {
		r *http.Request
	}
	r, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?persistent=BOTH&playMP4=false", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	r2, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?persistent=BOTH&playMP4=true", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: " happy: getplayMP4param success",
			args: args{
				r: r,
			},

			wantErr: false,
		},
		{
			name: " happy: getplayMP4param success",
			args: args{
				r: r2,
			},

			wantErr: false,
		},
		{
			name: "sad: getplayMP4param failed",
			args: args{
				r: nil,
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getplayMP4QueryParam(tt.args.r, "playMP4")
			if (err != nil) != tt.wantErr {
				t.Errorf("getplayMP4QueryParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_getOptionalVisibilityQueryParam(t *testing.T) {

	type args struct {
		r *http.Request
	}
	r, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?visibility=HIDE&playMP4=false", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	r2, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?visibility=SHOW&playMP4=true", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: " happy: getOptionalVisibilityQueryParam success",
			args: args{
				r: r,
			},

			wantErr: false,
		},
		{
			name: " happy: getOptionalVisibilityQueryParam success",
			args: args{
				r: r2,
			},

			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getOptionalVisibilityQueryParam(tt.args.r, "visibility")
			if (err != nil) != tt.wantErr {
				t.Errorf("getOptionalVisibilityQueryParam error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_getStringVar(t *testing.T) {

	type args struct {
		r     *http.Request
		value string
	}
	r, err := http.NewRequest(http.MethodGet,
		"/feed/nudges/{nudgeID}/unresolve/", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: " sad: getStringVar failed",
			args: args{
				r:     r,
				value: "visibility",
			},

			wantErr: true,
		},
		{
			name: " sad: getStringVar failed",
			args: args{
				r:     r,
				value: "expired",
			},

			wantErr: true,
		},

		{
			name: " sad: getStringVar failed",
			args: args{
				r:     nil,
				value: "expired",
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getStringVar(tt.args.r, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStringVarerror = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_getRequiredBooleanFilterQueryParam(t *testing.T) {

	type args struct {
		r     *http.Request
		value string
	}
	r, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?persistent=BOTH&playMP4=true", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	r2, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?persistent=ANY&playMP4=true", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: " happy: getRequiredBooleanFilterQueryParam",
			args: args{
				r:     r,
				value: "persistent",
			},

			wantErr: false,
		},
		{
			name: " sad: getRequiredBooleanFilterQueryParam failed",
			args: args{
				r:     r,
				value: "expired",
			},

			wantErr: true,
		},
		{
			name: " sad: getRequiredBooleanFilterQueryParam failed",
			args: args{
				r:     r2,
				value: "persistent",
			},

			wantErr: true,
		},
		{
			name: " sad: getRequiredBooleanFilterQueryParam failed",
			args: args{
				r:     nil,
				value: "",
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getRequiredBooleanFilterQueryParam(tt.args.r, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStringVarerror = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_getOptionalStatusQueryParam(t *testing.T) {

	type args struct {
		r     *http.Request
		value string
	}
	r, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?persistent=BOTH&playMP4=true", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: " happy: getOptionalStatusQueryParam",
			args: args{
				r:     r,
				value: "persistent",
			},

			wantErr: false,
		},
		{
			name: " happy: getOptionalStatusQueryParam passed",
			args: args{
				r:     r,
				value: "expired",
			},

			wantErr: false,
		},
		{
			name: " sad: getOptionalStatusQueryParam failed",
			args: args{
				r:     nil,
				value: "persistent",
			},

			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getOptionalStatusQueryParam(tt.args.r, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOptionalStatusQueryParam = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_getUIDFlavourAndIsAnonymous(t *testing.T) {

	type args struct {
		r *http.Request
	}
	r, err := http.NewRequest(http.MethodGet,
		"localhost/feed/?persistent=BOTH&playMP4=true", nil)
	if err != nil {
		t.Errorf("error in the request %v", err)
		return
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: " sad: getUIDFlavourAndIsAnonymous failed",
			args: args{
				r: r,
			},

			wantErr: true,
		},
		{
			name: " sad: getUIDFlavourAndIsAnonymous failed",
			args: args{
				r: nil,
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := getUIDFlavourAndIsAnonymous(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUIDFlavourAndIsAnonymous = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
