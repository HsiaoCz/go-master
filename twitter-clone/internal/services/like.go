package services

import (
	"errors"
	"twitter-clone/internal/database"
	"twitter-clone/internal/models"
)

type LikeService struct{}

func NewLikeService() *LikeService {
	return &LikeService{}
}

func (s *LikeService) LikeTweet(userID, tweetID uint) error {
	// Check if tweet exists
	var tweet models.Tweet
	if err := database.DB.First(&tweet, tweetID).Error; err != nil {
		return errors.New("tweet not found")
	}

	// Check if already liked
	var existingLike models.Like
	result := database.DB.Where("user_id = ? AND tweet_id = ?", userID, tweetID).First(&existingLike)
	if result.Error == nil {
		return errors.New("tweet already liked")
	}

	// Create new like
	like := &models.Like{
		UserID:  userID,
		TweetID: tweetID,
	}

	return database.DB.Create(like).Error
}

func (s *LikeService) UnlikeTweet(userID, tweetID uint) error {
	result := database.DB.Where("user_id = ? AND tweet_id = ?", userID, tweetID).Delete(&models.Like{})
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (s *LikeService) GetTweetLikes(tweetID uint) ([]models.Like, error) {
	var likes []models.Like
	err := database.DB.Where("tweet_id = ?", tweetID).
		Preload("User").
		Find(&likes).Error
	return likes, err
}