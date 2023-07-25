package config

// func TestLoadConfig(t *testing.T) {
// 	type tests struct {
// 		name       string
// 		filePath   string
// 		wantConfig func() Config
// 		wantErr    string
// 	}

// 	testCases := []tests{
// 		{
// 			name:     "valid case",
// 			filePath: "./config.yaml",
// 			wantConfig: func() Config {
// 				var config Config
// 				fileBytes, err := os.ReadFile("./config.yaml")
// 				if err != nil {
// 					log.Fatalln("error reading config file")
// 				}
// 				err = yaml.Unmarshal(fileBytes, &config)
// 				if err != nil {
// 					log.Fatalln("error marshalling config file")
// 				}
// 				return config
// 			},
// 			wantErr: "",
// 		},
// 		{
// 			name:     "file not found",
// 			filePath: ".",
// 			wantConfig: func() Config {
// 				return Config{}
// 			},
// 			wantErr: "Config File \"test\" Not Found in",
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			config, err := LoadConfig(tc.filePath, logger.ZapLogger{})
// 			wantConfig := tc.wantConfig()
// 			if !assert.Equal(t, tc.wantConfig(), config) {
// 				t.Errorf("expected config %v is different from actual config %v", wantConfig, config)
// 			}
// 			if err != nil {
// 				if !strings.ContainsAny(err.Error(), tc.wantErr) {
// 					t.Errorf("expected error %v is different from actual error %v", tc.wantErr, err)
// 				}
// 			}
// 		})
// 	}
// }
