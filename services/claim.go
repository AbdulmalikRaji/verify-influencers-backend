package services

import (
	"fmt"
	"time"

	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/claimDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/claimVerificationDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/influencerDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/influencerTopicDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/topicDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"github.com/abdulmalikraji/verify-influencers-backend/dto"
	"github.com/abdulmalikraji/verify-influencers-backend/pkg/gemini"
	"github.com/abdulmalikraji/verify-influencers-backend/pkg/serper"
	"github.com/abdulmalikraji/verify-influencers-backend/pkg/twitter"
	"github.com/abdulmalikraji/verify-influencers-backend/utils"
	"github.com/abdulmalikraji/verify-influencers-backend/utils/enums"
	"github.com/gofiber/fiber/v2"
)

type ClaimService interface {
	FindInfluencerClaims(ctx *fiber.Ctx, request dto.FindInfluencerClaimsRequest) (dto.FindInfluencerClaimsResponse, int, error)
}

type claimService struct {
	claimDao             claimDao.DataAccess
	influencerDao        influencerDao.DataAccess
	influencerTopicDao   influencerTopicDao.DataAccess
	topicDao             topicDao.DataAccess
	claimVerificationDao claimVerificationDao.DataAccess
}

func NewClaimService(
	claimDao claimDao.DataAccess,
	influencerDao influencerDao.DataAccess,
	influencerTopicDao influencerTopicDao.DataAccess,
	topicDao topicDao.DataAccess,
	claimVerificationDao claimVerificationDao.DataAccess,
) ClaimService {
	return &claimService{
		claimDao:             claimDao,
		influencerDao:        influencerDao,
		influencerTopicDao:   influencerTopicDao,
		topicDao:             topicDao,
		claimVerificationDao: claimVerificationDao,
	}
}

func (s *claimService) FindInfluencerClaims(ctx *fiber.Ctx, request dto.FindInfluencerClaimsRequest) (dto.FindInfluencerClaimsResponse, int, error) {

	var claims []dto.Claim
	var influencer models.Influencer

	if request.Source == 1 {

		claimTweets, err := twitter.GetTwitterClaimsV2(request.Username, utils.ConvertTimeToXFormat(request.StartDate), utils.ConvertTimeToXFormat(request.EndDate))
		if err != nil {
			return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
		}

		influencer, err = s.influencerDao.FindByUsername(request.Username)
		if err != nil && err.Error() == "record not found" {

			user, err := twitter.GetTwitterUserByUsername(request.Username)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}

			newInfluencer := models.Influencer{
				Name:           user.Name,
				Username:       request.Username,
				Platform:       "X",
				CreatedAt:      time.Now(),
				LastModifiedAt: time.Now(),
				TrustScore:     0,
				Followers:      user.UserPublicMetrics.Followers,
				URL:            user.URL,
				Bio:            user.Description,
				DelFlg:         false,
			}

			influencer, err = s.influencerDao.Insert(newInfluencer)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}

		} else if err != nil {
			return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
		}

		for _, tweet := range claimTweets {
			tweetTime, err := utils.ParseTweetTime(tweet.CreatedAt)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}

			parsedClaim, err := gemini.ExtractClaim(tweet.Text)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}
			if parsedClaim == "" {
				continue
			}

			topic, err := gemini.ExtractTopic(parsedClaim)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}

			claim := models.Claim{
				Content:      tweet.Text,
				ParsedClaim:  parsedClaim,
				Source:       "tweet",
				ClaimedAt:    tweetTime,
				InfluencerID: influencer.ID,
				Topic:        topic,
				SourceURL:    fmt.Sprintf("https://x.com/%s/status/%s", request.Username, tweet.ID),
			}

			// use claim model in anoter func for analysis and verification here
			madeClaim, err := s.claimDao.Insert(claim)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}
			err = s.AnalyzeAndVerifyClaim(madeClaim)
			if err != nil {
				return dto.FindInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}

			tweetClaim := dto.Claim{
				Raw:       tweet.Text,
				Source:    1,
				Timestamp: tweetTime,
				Claim:     parsedClaim,
				Topic:     topic,
				//InfluencerID: influencer.ID,
			}

			claims = append(claims, tweetClaim)
		}
		if len(claims) < 1 {
			return dto.FindInfluencerClaimsResponse{}, fiber.StatusNotFound, fmt.Errorf("No claims found for the specified username")
		}
	}

	return dto.FindInfluencerClaimsResponse{Claims: claims, Username: request.Username}, fiber.StatusOK, nil
}

func (s *claimService) AnalyzeAndVerifyClaim(claim models.Claim) error {
	// Step 1: Extract Topic

	// Check if topic exists, else create
	existingTopic, err := s.topicDao.FindByName(claim.Topic)
	if err != nil && err.Error() == "record not found" {
		newTopic := models.Topic{
			Name:      claim.Topic,
			CreatedAt: time.Now(),
			DelFlg:    false,
		}
		existingTopic, err = s.topicDao.Insert(newTopic)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Verify Claim and Get Score
	verificationResult, err := serper.VerifyClaim(claim.ParsedClaim)
	if err != nil {
		return err
	}

	result, err := gemini.GetScore(verificationResult.SearchParameters.Query, verificationResult.ResultStr)
	if err != nil {
		return err
	}

	// Determine Status Based on Score
	trustScore := result.Score // AI provides a score between 0 and 1
	var status string

	if trustScore >= 0.75 {
		status = enums.Verified
	} else if trustScore >= 0.4 {
		status = enums.Questionable
	} else {
		status = enums.Debunked
	}

	// Store Verification Result
	claimVerification := models.ClaimVerification{
		ClaimID:    claim.ID,
		VerifiedBy: result.BestResult.PublicationInfo,
		Status:     status,
		Score:      trustScore,
		Evidence:   result.BestResult.Title,
		Comment:    result.BestResult.Snippet,
		SourceUrl:  result.BestResult.Link,
		CreatedAt:  time.Now(),
		DelFlg:     false,
	}
	_, err = s.claimVerificationDao.Insert(claimVerification)
	if err != nil {
		return err
	}

	// Link Influencer to Topic
	influencerTopic := models.InfluencerTopic{
		InfluencerID: claim.InfluencerID,
		TopicID:      existingTopic.ID,
		CreatedAt:    time.Now(),
		DelFlg:       false,
	}
	_, err = s.influencerTopicDao.Insert(influencerTopic)
	if err != nil {
		return err
	}

	return nil
}
