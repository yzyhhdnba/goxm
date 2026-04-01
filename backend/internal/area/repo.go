package area

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Area{})
}

func (r *Repository) SeedDefaults(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "slug"}},
			DoNothing: true,
		}).
		Create(defaultAreas()).
		Error
}

func (r *Repository) ListActive(ctx context.Context) ([]Area, error) {
	var areas []Area
	if err := r.db.WithContext(ctx).
		Where("status = ?", StatusActive).
		Order("sort_order ASC").
		Order("id ASC").
		Find(&areas).
		Error; err != nil {
		return nil, err
	}
	return areas, nil
}

func defaultAreas() []Area {
	return []Area{
		{ID: 1, Name: "动画", Slug: "anime", SortOrder: 10, Status: StatusActive},
		{ID: 2, Name: "番剧", Slug: "bangumi", SortOrder: 20, Status: StatusActive},
		{ID: 3, Name: "音乐", Slug: "music", SortOrder: 30, Status: StatusActive},
		{ID: 4, Name: "游戏", Slug: "game", SortOrder: 40, Status: StatusActive},
		{ID: 5, Name: "科技", Slug: "tech", SortOrder: 50, Status: StatusActive},
		{ID: 6, Name: "生活", Slug: "life", SortOrder: 60, Status: StatusActive},
	}
}
