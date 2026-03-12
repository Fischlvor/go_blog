package advertisement

import (
	"context"

	"server-blog-v2/config"
	"server-blog-v2/internal/repo"
	"server-blog-v2/internal/usecase"
	"server-blog-v2/internal/usecase/input"
	"server-blog-v2/internal/usecase/output"
)

type useCase struct {
	cfg  *config.Config
	repo repo.AdvertisementRepo
}

// New 创建 Advertisement UseCase。
func New(cfg *config.Config, repo repo.AdvertisementRepo) usecase.Advertisement {
	return &useCase{
		cfg:  cfg,
		repo: repo,
	}
}

func (u *useCase) GetInfo(ctx context.Context) (*output.AdvertisementInfo, error) {
	ads, total, err := u.repo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]output.AdvertisementItem, len(ads))
	for i, ad := range ads {
		// 拼接完整 URL
		adImage := ad.AdImage
		if u.cfg.Qiniu.Domain != "" && adImage != "" && adImage[0] != 'h' {
			protocol := "http://"
			if u.cfg.Qiniu.UseHTTPS {
				protocol = "https://"
			}
			adImage = protocol + u.cfg.Qiniu.Domain + "/" + adImage
		}

		items[i] = output.AdvertisementItem{
			ID:       ad.ID,
			AdName:   ad.AdName,
			AdImage:  adImage,
			AdLink:   ad.AdLink,
			AdType:   ad.AdType,
			IsActive: ad.IsActive,
		}
	}

	return &output.AdvertisementInfo{
		List:  items,
		Total: total,
	}, nil
}

func (u *useCase) List(ctx context.Context, params input.ListAdvertisements) (*output.ListResult[output.AdvertisementItem], error) {
	offset := (params.Page - 1) * params.PageSize

	ads, total, err := u.repo.List(ctx, offset, params.PageSize, params.Title)
	if err != nil {
		return nil, err
	}

	items := make([]output.AdvertisementItem, len(ads))
	for i, ad := range ads {
		// 拼接完整 URL
		adImage := ad.AdImage
		if u.cfg.Qiniu.Domain != "" && adImage != "" && adImage[0] != 'h' {
			protocol := "http://"
			if u.cfg.Qiniu.UseHTTPS {
				protocol = "https://"
			}
			adImage = protocol + u.cfg.Qiniu.Domain + "/" + adImage
		}

		items[i] = output.AdvertisementItem{
			ID:       ad.ID,
			AdName:   ad.AdName,
			AdImage:  adImage,
			AdLink:   ad.AdLink,
			AdType:   ad.AdType,
			IsActive: ad.IsActive,
		}
	}

	return &output.ListResult[output.AdvertisementItem]{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}
