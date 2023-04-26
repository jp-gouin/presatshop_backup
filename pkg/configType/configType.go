package configType

type ConfigDB struct {
	User    string
	DBName  string
	Passwd  string
	Addr    string
	DumpDir string
}
type ConfigS3 struct {
	URL       string
	AccessKey string
	SecretKey string
}
type ConfigSlack struct {
	ChannelID string
	APIToken  string
}

type Config struct {
	DB    ConfigDB
	S3    ConfigS3
	Slack ConfigSlack
}
