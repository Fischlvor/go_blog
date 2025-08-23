package router

type RouterGroup struct {
	BaseRouter
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
}

var RouterGroupApp = new(RouterGroup)
