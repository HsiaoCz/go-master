package services

import (
	"errors"
	"twitter-clone/internal/database"
	"twitter-clone/internal/models"
)

type TweetService struct{}

func NewTweetService() *TweetService {
	return &TweetService{}
}

func (s *TweetService) CreateTweet(userID uint, content string, parentID *uint) (*models.Tweet, error) {
	tweet := &models.Tweet{
		UserID:   userID,
		Content:  content,
		ParentID: parentID,
	}

	result := database.DB.Create(tweet)
	if result.Error != nil {
		return nil, result.Error
	}

	// Load the user relationship
	database.DB.Preload("User").First(tweet, tweet.ID)
	return tweet, nil
}

func (s *TweetService) GetTweet(tweetID uint) (*models.Tweet, error) {
	var tweet models.Tweet
	result := database.DB.Preload("User").
		Preload("Likes").
		Preload("Retweets").
		First(&tweet, tweetID)

	if result.Error != nil {
		return nil, result.Error
	}
	return &tweet, nil
}

func (s *TweetService) UpdateTweet(tweetID, userID uint, content string) (*models.Tweet, error) {
	var tweet models.Tweet
	if err := database.DB.First(&tweet, tweetID).Error; err != nil {
		return nil, err
	}

	if tweet.UserID != userID {
		return nil, errors.New("unauthorized to update this tweet")
	}

	tweet.Content = content
	if err := database.DB.Save(&tweet).Error; err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (s *TweetService) DeleteTweet(tweetID, userID uint) error {
	var tweet models.Tweet
	if err := database.DB.First(&tweet, tweetID).Error; err != nil {
		return err
	}

	if tweet.UserID != userID {
		return errors.New("unauthorized to delete this tweet")
	}

	return database.DB.Delete(&tweet).Error
}

func (s *TweetService) GetUserTweets(userID uint) ([]models.Tweet, error) {
	var tweets []models.Tweet
	result := database.DB.Where("user_id = ?", userID).
		Preload("User").
		Preload("Likes").
		Preload("Retweets").
		Find(&tweets)

	if result.Error != nil {
		return nil, result.Error
	}
	return tweets, nil
}

func (s *TweetService) GetHomeTimeline(userID uint, page, limit int) ([]models.Tweet, error) {
	var tweets []models.Tweet
	offset := (page - 1) * limit

	// Get tweets from users that the current user follows and their own tweets
	err := database.DB.
		Preload("User").
		Preload("Likes").
		Preload("Retweets").
		Table("tweets").
		Joins("LEFT JOIN follows ON tweets.user_id = follows.following_id").
		Where("follows.follower_id = ? OR tweets.user_id = ?", userID, userID).
		Where("tweets.parent_id IS NULL"). // Exclude replies
		Order("tweets.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tweets).Error

	return tweets, err
}

func (s *TweetService) GetUserTimeline(userID uint, page, limit int) ([]models.Tweet, error) {
	var tweets []models.Tweet
	offset := (page - 1) * limit

	// Get tweets from the specified user
	err := database.DB.
		Preload("User").
		Preload("Likes").
		Preload("Retweets").
		Where("user_id = ?", userID).
		Where("parent_id IS NULL"). // Exclude replies
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tweets).Error

	return tweets, err
}

func (s *TweetService) GetReplies(tweetID uint, page, limit int) ([]models.Tweet, error) {
	var replies []models.Tweet
	offset := (page - 1) * limit

	err := database.DB.
		Preload("User").
		Preload("Likes").
		Preload("Retweets").
		Where("parent_id = ?", tweetID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&replies).Error

	return replies, err
}

// Add this method to get timeline with additional metadata
type TimelineResponse struct {
	Tweets    []models.Tweet `json:"tweets"`
	HasMore   bool           `json:"has_more"`
	TotalPage int            `json:"total_page"`
	Page      int            `json:"page"`
}

func (s *TweetService) GetTimelineWithMetadata(userID uint, page, limit int, timelineType string) (*TimelineResponse, error) {
	var total int64
	var tweets []models.Tweet
	var err error

	// Count total tweets
	query := database.DB.Model(&models.Tweet{})
	if timelineType == "home" {
		query = query.
			Joins("LEFT JOIN follows ON tweets.user_id = follows.following_id").
			Where("follows.follower_id = ? OR tweets.user_id = ?", userID, userID).
			Where("tweets.parent_id IS NULL")
	} else if timelineType == "user" {
		query = query.
			Where("user_id = ?", userID).
			Where("parent_id IS NULL")
	}
	query.Count(&total)

	// Get tweets for current page
	if timelineType == "home" {
		tweets, err = s.GetHomeTimeline(userID, page, limit)
	} else {
		tweets, err = s.GetUserTimeline(userID, page, limit)
	}

	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasMore := page < totalPages

	return &TimelineResponse{
		Tweets:    tweets,
		HasMore:   hasMore,
		TotalPage: totalPages,
		Page:      page,
	}, nil
}
