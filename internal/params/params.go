package params

import "github.com/vrischmann/envconfig"

type params struct {
	Username string `envconfig:"default="`
	Password string `envconfig:"default="`
	BaseURL  string `envconfig:"default="`
}

func Parse() (*params, error) {
	var p params
	if err := envconfig.InitWithPrefix(&p, "APP"); err != nil {
		return nil, err
	}
	return &p, nil
}
