package service

import (
	"github.com/sankethkini/NewsLetter-Backend/internal/service/admin"
	newsletter "github.com/sankethkini/NewsLetter-Backend/internal/service/news_letter"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/subscription"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/user"
)

type Registry struct {
	UserService         user.Service
	SubscriptionService subscription.Service
	AdminService        admin.Service
	NewsService         newsletter.Service
}

func NewRegistry(us user.Service, sb subscription.Service, ad admin.Service, nw newsletter.Service) *Registry {
	return &Registry{
		UserService:         us,
		SubscriptionService: sb,
		AdminService:        ad,
		NewsService:         nw,
	}
}
