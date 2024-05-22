package ext

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	innerWhere string
	where      *uint
	values     map[string]string
	order      *uint
	page       *Page
	model      interface{}
	db         *gorm.DB
}

// Add 新增
func (s *Service) Add(value interface{}) error {
	tx := s.db.Begin()
	if err := tx.Create(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Edit 修改
func (s *Service) Edit(value interface{}) error {
	tx := s.db.Begin()
	if err := tx.Save(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Delete 删除
func (s *Service) Delete(id uint) error {
	tx := s.db.Begin()
	if err := tx.Delete(&s.model, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Take 查询
func (s *Service) Take(id uint, dest interface{}) error {
	tx := s.db.Begin()

	if err := tx.
		Preload(clause.Associations).
		Take(dest, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Fetch 查询
func (s *Service) Fetch(datas interface{}, total *int64) error {
	scopes := []func(*gorm.DB) *gorm.DB{}

	if s.innerWhere != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(s.innerWhere)
		})
	}

	if s.where != nil {
		query, vals, err := WhereBuild(s.db, *s.where, s.values)
		if err != nil {
			return err
		}
		if query != "" {
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(query, vals...)
			})
		}
	}

	if err := s.db.Scopes(scopes...).
		Model(s.model).
		Count(total).Error; err != nil {
		return err
	}

	if s.order != nil {
		s := OrderByBuild(s.db, *s.order)
		if s != "" {
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Order(s)
			})
		}
	}

	if err := s.db.Scopes(scopes...).
		Preload(clause.Associations).
		Offset(s.page.Index * s.page.Size).
		Limit(s.page.Size).
		Find(datas).Error; err != nil {
		return err
	}

	return nil
}
