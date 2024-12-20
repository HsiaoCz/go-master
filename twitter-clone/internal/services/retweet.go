package services

import (
	"errors"
	"twitter-clone/internal/database"
	"twitter-clone/internal/models"
)

type RetweetService struct{}

func NewRetweetService() *RetweetService {
	return &RetweetService{}
}

func (s *RetweetService) CreateRetweet(userID, tweetID uint) error {
	// Check if tweet exists
	var tweet models.Tweet
	if err := database.DB.First(&tweet, tweetID).Error; err != nil {
		return errors.New("tweet not found")
	}

	// Check if already retweeted
	var existingRetweet models.Retweet
	result := database.DB.Where("user_id = ? AND tweet_id = ?", userID, tweetID).First(&existingRetweet)
	if result.Error == nil {
		return errors.New("tweet already retweeted")
	}

	// Create new retweet
	retweet := &models.Retweet{
		UserID:  userID,
		TweetID: tweetID,
	}

	return database.DB.Create(retweet).Error
}

func (s *RetweetService) UndoRetweet(userID, tweetID uint) error {
	result := database.DB.Where("user_id = ? AND tweet_id = ?", userID, tweetID).Delete(&models.Retweet{})
	if result.RowsAffected == 0 {
		return errors.New("retweet not found")
	}
	return nil
}

func (s *RetweetService) GetTweetRetweets(tweetID uint) ([]models.Retweet, error) {
	var retweets []models.Retweet
	err := database.DB.Where("tweet_id = ?", tweetID).
		Preload("User").
		Find(&retweets).Error
	return retweets, err
}