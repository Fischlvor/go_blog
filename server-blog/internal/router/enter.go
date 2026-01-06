package router

type RouterGroup struct {
	BaseRouter
	AuthRouter
	UserRouter
	ImageRouter
	ArticleRouter
	CommentRouter
	AdvertisementRouter
	FriendLinkRouter
	FeedbackRouter
	WebsiteRouter
	ConfigRouter
	AIChatRouter
	AIManagementRouter
	EmojiRouter
	PublicEmojiRouter
	ResourceRouter
}

var RouterGroupApp = new(RouterGroup)
