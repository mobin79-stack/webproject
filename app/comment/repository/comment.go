package repository

import (
	"context"

	"git.ecobin.ir/ecomicro/template/app/comment/domain"

	"git.ecobin.ir/ecomicro/tooty"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type commentRepository struct {
	Conn *gorm.DB
}

var _ domain.Repository = &commentRepository{}

func NewCommentRepository(dbConnection *gorm.DB) *commentRepository {
	err := dbConnection.AutoMigrate(&Comment{})
	if err != nil {
		panic(err)
	}
	return &commentRepository{dbConnection}
}

func (ur *commentRepository) Create(ctx context.Context, domainComment domain.Comment) (*domain.Comment, error) {
	span := tooty.OpenAnAPMSpan(ctx, "[R] create comment", "repository")
	defer tooty.CloseTheAPMSpan(span)

	commentDao := FromDomainComment(domainComment)
	result := ur.Conn.Debug().Create(&commentDao)
	if result.Error != nil {
		return nil, result.Error
	}
	comment := commentDao.ToDomainComment()
	return &comment, nil
}

func (ur *commentRepository) GetCommentById(ctx context.Context, id uint64) (*domain.Comment, error) {
	span := tooty.OpenAnAPMSpan(ctx, "[R] get comment by id", "repository")
	defer tooty.CloseTheAPMSpan(span)
	var commentDao Comment
	err := ur.Conn.WithContext(ctx).Debug().Where(Comment{Id: id}).Find(&commentDao).Error
	if err != nil {
		return nil, err
	}
	comment := commentDao.ToDomainComment()
	return &comment, nil
}

func (ur *commentRepository) Update(ctx context.Context, condition domain.Comment, domainComment domain.Comment) ([]domain.Comment, error) {
	span := tooty.OpenAnAPMSpan(ctx, "[R] update comment", "repository")
	defer tooty.CloseTheAPMSpan(span)
	var commentArray []Comment
	err := ur.Conn.WithContext(ctx).Debug().Model(&commentArray).Clauses(clause.Returning{}).Where(FromDomainComment(condition)).Updates(FromDomainComment(domainComment)).Error
	if err != nil {
		return []domain.Comment{}, err
	}
	domainComments := make([]domain.Comment, len(commentArray))
	for idx, comment := range commentArray {
		domainComments[idx] = comment.ToDomainComment()
	}
	return domainComments, nil
}