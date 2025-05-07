package config

type AuthCodeConfig struct {
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	SignName        string `json:"sign_name" yaml:"sign_name"`
	TemplateCode    string `json:"template_code" yaml:"template_code"`
}
