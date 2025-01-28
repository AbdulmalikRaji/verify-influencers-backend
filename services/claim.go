package services

import (
	"fmt"
	"time"

	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/claimDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/influencerDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"github.com/abdulmalikraji/verify-influencers-backend/dto"
	"github.com/abdulmalikraji/verify-influencers-backend/pkg/gemini"
	"github.com/abdulmalikraji/verify-influencers-backend/pkg/twitter"
	"github.com/abdulmalikraji/verify-influencers-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClaimService interface {
	GetInfluencerClaims(ctx *fiber.Ctx, request dto.GetInfluencerClaimsRequest) (dto.GetInfluencerClaimsResponse, int, error)
}

type claimService struct {
	claimDao      claimDao.DataAccess
	influencerDao influencerDao.DataAccess
}

func NewClaimService(claimDao claimDao.DataAccess, influencerDao influencerDao.DataAccess) ClaimService {
	return &claimService{
		claimDao:      claimDao,
		influencerDao: influencerDao,
	}
}

func (s *claimService) GetInfluencerClaims(ctx *fiber.Ctx, request dto.GetInfluencerClaimsRequest) (dto.GetInfluencerClaimsResponse, int, error) {

	var claims []dto.Claim
	var influencer models.Influencer

	if request.Source == 1 {

		claimTweets, err := twitter.GetTwitterClaimsV2(request.Username, utils.ConvertTimeToXFormat(request.StartDate), utils.ConvertTimeToXFormat(request.EndDate))
		if err != nil {
			return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
		}

		influencer, err = s.influencerDao.FindByUsername(request.Username)
		if err != nil && err == gorm.ErrRecordNotFound {

			user, err := twitter.GetTwitterUserByUsername(request.Username)
			if err != nil {
				return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
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
				DelFlg:         false,
			}
			influencer, err = s.influencerDao.Insert(newInfluencer)
			if err != nil {
				return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}
		}
		if err != nil {
			return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
		}

		for _, tweet := range claimTweets {
			tweetTime, err := utils.ParseTweetTime(tweet.CreatedAt)
			if err != nil {
				return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
			}

			parsedClaim, err := gemini.ExtractClaim(tweet.Text)
			if err!= nil {
                return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
            }

			claim := models.Claim{
				Content:      tweet.Text,
				ParsedClaim: parsedClaim,
				Source:       "tweet",
				ClaimedAt:    tweetTime,
				InfluencerID: influencer.ID,
				SourceURL:    fmt.Sprintf("https://x.com/%s/status/%s", request.Username, tweet.ID),
			}

			err = s.claimDao.Insert(claim)
			if err!= nil {
                return dto.GetInfluencerClaimsResponse{}, fiber.StatusInternalServerError, err
            }

			tweetClaim := dto.Claim{
				Raw:       tweet.Text,
				Source:    1,
				Timestamp: tweetTime,
				//InfluencerID: influencer.ID,
			}

			claims = append(claims, tweetClaim)
		}
	}

	return dto.GetInfluencerClaimsResponse{Claims: claims, Username: request.Username}, fiber.StatusOK, nil
}
