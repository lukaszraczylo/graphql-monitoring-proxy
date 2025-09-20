package main

func (suite *Tests) Test_extractClaimsFromJWTHeader() {
	jwt_token_for_tests := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiSGFzdXJhIjp7IngtaGFzdXJhLWFsbG93ZWQtcm9sZXMiOlsiZ3Vlc3QiLCJ1c2VyIiwiZ3JvdXBhZG1pbiIsInBheWFkbWluIl0sIngtaGFzdXJhLWRlZmF1bHQtcm9sZSI6Imd1ZXN0IiwieC1oYXN1cmEtdXNlci1pZCI6IjE2NyIsIngtaGFzdXJhLXVzZXItdXVpZCI6ImRkM2U2ZTM1LTA0MDktNDNiMC1iZmYxLWNlZjNjNmVkNWYxMCJ9LCJpc3MiOiJBdXRoU2VydmljZSIsImV4cCI6MTY5NjgwMTcyNiwibmJmIjoxNjk2NTg1NzI2LCJpYXQiOjE2OTY1ODU3MjZ9.dsJ5JKzG5tXOlqeZ_Gfe2XC-vyrcwtYwOGfhvt8q9UY"

	type args struct {
		authorization string
	}

	tests := []struct {
		name           string
		args           args
		wantUsr        string
		wantRole       string
		jwt_token_path string
		jwt_role_path  string
	}{
		{
			name:     "test_empty",
			wantUsr:  "-",
			wantRole: "-",
		},
		{
			name: "test_invalid_path",
			args: args{
				authorization: jwt_token_for_tests,
			},
			wantUsr:        "-",
			wantRole:       "-",
			jwt_token_path: "invalid",
		},
		{
			name: "test_invalid_role_path",
			args: args{
				authorization: jwt_token_for_tests,
			},
			wantUsr:       "-",
			wantRole:      "-",
			jwt_role_path: "invalid",
		},
		{
			name: "test_valid",
			args: args{
				authorization: jwt_token_for_tests,
			},
			wantUsr:        "167",
			wantRole:       "guest",
			jwt_token_path: "Hasura.x-hasura-user-id",
			jwt_role_path:  "Hasura.x-hasura-default-role",
		},
		{
			name: "test_invalid_token",
			args: args{
				authorization: "invalid",
			},
			wantUsr:  "-",
			wantRole: "-",
		},
		{
			name: "test_invalid_three_part_token",
			args: args{
				authorization: "invalid.threepart.token",
			},
			wantUsr:  "-",
			wantRole: "-",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if len(tt.jwt_token_path) > 0 {
				cfg.Client.JWTUserClaimPath = tt.jwt_token_path
			}
			if len(tt.jwt_role_path) > 0 {
				cfg.Client.JWTRoleClaimPath = tt.jwt_role_path
			}
			gotUsr, gotRole := extractClaimsFromJWTHeader(tt.args.authorization)
			suite.Equal(tt.wantUsr, gotUsr, "Unexpected user ID")
			suite.Equal(tt.wantRole, gotRole, "Unexpected role")
		})
	}
}
