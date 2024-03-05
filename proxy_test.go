package main

func (suite *Tests) Test_proxyTheRequest() {

	supplied_headers := map[string]string{
		"X-Forwarded-For": "127.0.0.1",
		"Content-Type":    "application/json",
		"Method":          "POST",
	}

	tests := []struct {
		name    string
		query   string
		host    string
		path    string
		headers map[string]string
		wantErr bool
	}{
		{
			name: "test_empty",
			query: `query {
				__type(name: "Query") {
					name
				}
			}`,
			host:    "https://telegram-bot.app",
			path:    "/v1/graphql",
			headers: supplied_headers,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
		})
	}
}
