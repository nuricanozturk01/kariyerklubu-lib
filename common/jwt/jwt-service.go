package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nuricanozturk01/kariyerklubu-lib/common/config"
	"strconv"
	"time"
)

const AccessTokenDuration = time.Hour * 24
const RefreshTokenDuration = time.Hour * 24 * 30
const ResetPasswordTokenDuration = time.Minute * 15
const (
	ClaimUserID string = "user_id"
)

type Jwt struct {
	configuration *config.Config
}

func NewJwt(configuration *config.Config) *Jwt {
	return &Jwt{
		configuration: configuration,
	}
}

func (j *Jwt) GenerateAccessToken(roles []string, userID, email string) (string, error) {
	jwtSecret := j.configuration.GetEnvironment(config.JwtSecret)

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
		"roles":   roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

func (j *Jwt) GenerateAccessAndRefreshToken(roles []string, userID, email string) (string, string, error) {
	accessToken, err := j.GenerateAccessToken(roles, userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *Jwt) GenerateRefreshToken(userID string) (string, error) {
	jwtSecret := j.configuration.GetEnvironment(config.JwtSecret)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(RefreshTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

func (j *Jwt) ValidateToken(tokenString string) error {
	expirationStr, err := j.GetClaim(tokenString, "exp")
	if err != nil {
		return err
	}

	exp, err := strconv.ParseInt(expirationStr, 10, 64)
	if err != nil {
		return err
	}

	if time.Unix(exp, 0).Before(time.Now()) {
		println("token is expired")
		return fmt.Errorf("token is expired")
	} else {
		println("token is not expired")
	}

	return nil
}

func (j *Jwt) getToken(tokenStr string) (*jwt.Token, error) {
	jwtSecret := j.configuration.GetEnvironment(config.JwtSecret)

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j *Jwt) GetClaim(tokenString, key string) (string, error) {
	jwtSecret := j.configuration.GetEnvironment(config.JwtSecret)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if value, found := claims[key]; found {
			if strValue, ok := value.(string); ok {
				return strValue, nil
			}
			return "", fmt.Errorf("claim %s is not a string", key)
		}
		return "", fmt.Errorf("claim %s not found", key)
	}

	return "", fmt.Errorf("invalid token or claims")
}

func (j *Jwt) GetRoles(tokenString, key string) ([]string, error) {
	jwtSecret := j.configuration.GetEnvironment(config.JwtSecret)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if value, found := claims[key]; found {
			if rolesMap, ok := value.([]interface{}); ok {
				var result []string
				for _, role := range rolesMap {
					if roleStr, ok := role.(string); ok {
						result = append(result, roleStr)
					} else {
						return nil, fmt.Errorf("unexpected role type: %T", role)
					}
				}
				return result, nil
			}
			return nil, fmt.Errorf("unexpected value type: %T", value)
		}
		return nil, fmt.Errorf("key '%s' not found in claims", key)
	}

	return nil, fmt.Errorf("invalid token or claims")
}

func (j *Jwt) GenerateRefreshPasswordToken(userId, email string) string {
	jwtSecret := j.configuration.GetEnvironment(config.JwtSecret)

	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userId,
		"exp":     time.Now().Add(ResetPasswordTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshPasswordToken, _ := token.SignedString([]byte(jwtSecret))

	return refreshPasswordToken
}
