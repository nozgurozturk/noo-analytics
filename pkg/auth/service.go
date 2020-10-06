package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nozgurozturk/noo-analytics/entities"
	"github.com/nozgurozturk/noo-analytics/internal/config"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"strconv"
	"time"
)

type Service interface {
	CreateToken(u *entities.User, config *config.ServerConfig) (*entities.TokenResponse, *errors.ApplicationError)
	ValidateToken(t string, tType string, config *config.ServerConfig) (*entities.Token, *errors.ApplicationError)
	DeleteToken(token *entities.Token) *errors.ApplicationError
	ExtractTokenClaims(t string, tType string, config *config.ServerConfig) (*entities.Token, *errors.ApplicationError)
}
type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) CreateToken(u *entities.User, config *config.ServerConfig) (*entities.TokenResponse, *errors.ApplicationError) {
	var err error
	at := &entities.Token{}

	at.Expires = time.Now().Add(time.Minute * time.Duration(config.AccessExpire)).Unix()
	at.Uuid = uuid.New().String()
	at.UserId = u.ID.Hex()
	at.Role = "Member"
	at.Type = "Access"
	if u.Email == config.Admin {
		at.Role = "Admin"
	}

	atClaims := jwt.MapClaims{
		"exp":    at.Expires,
		"uuid":   at.Uuid,
		"userId": at.UserId,
		"role":   at.Role,
		"type":   at.Type,
	}
	atToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	at.Token, err = atToken.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return nil, errors.InternalServer("Access Token could not created")
	}

	sErr := s.repository.SetToken(at)

	if sErr != nil {
		return nil, sErr
	}

	rt := &entities.Token{}

	rt.Expires = time.Now().Add(time.Minute * time.Duration(config.RefreshExpire)).Unix()
	rt.Uuid = uuid.New().String()
	rt.UserId = u.ID.Hex()
	rt.Role = "Member"
	rt.Type = "Refresh"
	if u.Email == config.Admin {
		rt.Role = "Admin"
	}

	rtClaims := jwt.MapClaims{
		"exp":    rt.Expires,
		"uuid":   rt.Uuid,
		"userId": rt.UserId,
		"role":   rt.Role,
		"type":   rt.Type,
	}
	rtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	rt.Token, err = rtToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return nil, errors.InternalServer("Access Token could not created")
	}

	return &entities.TokenResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *service) DeleteToken(token *entities.Token) *errors.ApplicationError {
	err := s.repository.DeleteToken(token)
	if err != nil {
		return err
	}
	return nil
}

func verifyToken(t string, tokenType string, config *config.ServerConfig) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if tokenType == "Access" {
			return []byte(config.AccessSecret), nil
		}

		if tokenType == "Refresh" {
			return []byte(config.RefreshSecret), nil
		}
		return nil, fmt.Errorf("anexpected token type: %s", tokenType)
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *service) ValidateToken(t string, tType string, config *config.ServerConfig) (*entities.Token, *errors.ApplicationError) {
	token, err := verifyToken(t, tType, config)
	if err != nil {
		return nil, errors.Unauthorized(err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		_uuid, ok := claims["uuid"].(string)
		if !ok {
			return nil, errors.Unauthorized("can't access uuid in token claims")
		}
		expires, parseErr := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		if parseErr != nil {
			return nil, errors.Unauthorized("can't access expire time in token claims")
		}
		userId, ok := claims["userId"].(string)
		if !ok {
			return nil, errors.Unauthorized("can't access user id in token claims")
		}
		role, ok := claims["role"].(string)
		if !ok {
			return nil, errors.Unauthorized("can't access user role in token claims")
		}
		if role != "Admin" {
			return nil, errors.Unauthorized("You can't access this data")
		}
		tokenType, ok := claims["type"].(string)
		if !ok {
			return nil, errors.Unauthorized("can't access token type in token claims")
		}
		return &entities.Token{
			Uuid:    _uuid,
			Expires: expires,
			UserId:  userId,
			Role:    role,
			Type:    tokenType,
		}, nil
	}
	return nil, errors.InternalServer("Unhandled token validation error")
}

func (s *service) ExtractTokenClaims(t string, tType string, config *config.ServerConfig) (*entities.Token, *errors.ApplicationError) {
	token, err := s.ValidateToken(t, tType, config)
	if err != nil {
		return nil, err
	}
	sErr := s.repository.DeleteToken(token)
	if sErr != nil {
		return nil, sErr
	}
	return nil, errors.Unauthorized("Unauthorized user")
}
