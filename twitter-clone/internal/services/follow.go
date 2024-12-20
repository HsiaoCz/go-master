package services

import (
	"errors"
	"twitter-clone/internal/database"
	"twitter-clone/internal/models"
)

type FollowService struct{}

func NewFollowService() *FollowService {
	return &FollowService{}
}

func (s *FollowService) FollowUser(followerID, followingID uint) error {
	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, followingID).Error; err != nil {
		return errors.New("user to follow not found")
	}

	// Prevent self-following
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Check if already following
	var existingFollow models.Follow
	result := database.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollow)
	if result.Error == nil {
		return errors.New("already following this user")
	}

	// Create new follow relationship
	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return database.DB.Create(follow).Error
}

func (s *FollowService) UnfollowUser(followerID, followingID uint) error {
	result := database.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{})
	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}
	return nil
}

func (s *FollowService) GetFollowers(userID uint) ([]models.User, error) {
	var followers []models.User
	err := database.DB.Model(&models.User{}).
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Find(&followers).Error
	return followers, err
}

func (s *FollowService) GetFollowing(userID uint) ([]models.User, error) {
	var following []models.User
	err := database.DB.Model(&models.User{}).
		Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Find(&following).Error
	return following, err
}

func (s *FollowService) IsFollowing(followerID, followingID uint) bool {
	var follow models.Follow
	result := database.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow)
	return result.Error == nil
}
