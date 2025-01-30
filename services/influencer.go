package services

import (
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/claimDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/claimVerificationDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/influencerDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/influencerTopicDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/topicDao"
	"github.com/abdulmalikraji/verify-influencers-backend/dto"
	"github.com/gofiber/fiber/v2"
)

type InfluencerService interface {
	GetInfluencer(ctx *fiber.Ctx, request dto.GetInfluencerRequest) (dto.GetInfluencerResponse, int, error)
}

type influencerService struct {
	claimDao             claimDao.DataAccess
	influencerDao        influencerDao.DataAccess
	influencerTopicDao   influencerTopicDao.DataAccess
	topicDao             topicDao.DataAccess
	claimVerificationDao claimVerificationDao.DataAccess
}

func NewInfluencerService(
	claimDao claimDao.DataAccess,
	influencerDao influencerDao.DataAccess,
	influencerTopicDao influencerTopicDao.DataAccess,
	topicDao topicDao.DataAccess,
	claimVerificationDao claimVerificationDao.DataAccess,
) InfluencerService {
	return &influencerService{
		claimDao:             claimDao,
		influencerDao:        influencerDao,
		influencerTopicDao:   influencerTopicDao,
		topicDao:             topicDao,
		claimVerificationDao: claimVerificationDao,
	}
}

func (s influencerService) GetInfluencer(ctx *fiber.Ctx, request dto.GetInfluencerRequest) (dto.GetInfluencerResponse, int, error) {

	// Get influencer from the database
	influencer, err := s.influencerDao.FindById(*request.ID)
	if err != nil {
		return dto.GetInfluencerResponse{}, fiber.StatusInternalServerError, err
	}

	// Get topics associated with the influencer
	influencerTopics, err := s.influencerTopicDao.FindAllByInfluencerId(*request.ID)
	if err != nil {
		return dto.GetInfluencerResponse{}, fiber.StatusInternalServerError, err
	}
	var topics []string
	for _, influencerTopic := range influencerTopics {
		topic, err := s.topicDao.FindById(influencerTopic.ID)
		if err != nil {
			return dto.GetInfluencerResponse{}, fiber.StatusInternalServerError, err
		}
		topics = append(topics, topic.Name)
	}

	// Get claims associated with the influencer
	madeClaims, err := s.claimDao.FindAllByInfluencerId(*request.ID)
	if err != nil {
		return dto.GetInfluencerResponse{}, fiber.StatusInternalServerError, err
	}

	var claims []dto.InfluencerClaim
	score := 0.0

	for _, claim := range madeClaims {
		verification, err := s.claimVerificationDao.FindByClaimId(claim.ID)
		if err != nil {
			return dto.GetInfluencerResponse{}, fiber.StatusInternalServerError, err
		}

		score += verification.Score

		claims = append(claims, dto.InfluencerClaim{
			Claim:       claim.ParsedClaim,
			ClaimURL:    claim.SourceURL,
			Proof:       verification.Comment,
			ProofSource: verification.Evidence,
			ProofURL:    verification.SourceUrl,
			Status:      verification.Status,
			Score:       verification.Score,
			Topic:       claim.Topic,
			ClaimedAt:   claim.ClaimedAt,
		})

	}

	response := dto.GetInfluencerResponse{
		Name:       influencer.Name,
		Username:   influencer.Username,
		Followers:  influencer.Followers,
		URL:        influencer.URL,
		Bio:        influencer.Bio,
		TrustScore: (score * 100) / float64(len(claims)),
		Topics:     topics,
		Claims:     claims,
	}

	return response, fiber.StatusOK, nil
}
