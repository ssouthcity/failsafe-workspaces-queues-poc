package newsfeed

type Category string

const (
	News          Category = "News"
	Videos        Category = "Videos"
	SocialMedia   Category = "Social Media"
	ServerUpdates Category = "Server Updates"
	Community     Category = "Community"
)

func (c Category) String() string {
	return string(c)
}
