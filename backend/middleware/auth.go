package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストヘッダーからアクセストークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// アクセストークンの形式をチェック（例: "Bearer token"）
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := headerParts[1]

		// アクセストークンの検証
		if !isValidToken(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			return
		}

		// アクセストークンが有効な場合、次のハンドラーに進む
		c.Next()
	}
}

func isValidToken(token string) bool {
	// アクセストークンの検証ロジックをここに実装する
	// 例えば、データベースやトークン発行サービスに問い合わせるなど
	// デモンストレーションのため、固定のトークンを使用
	validToken := "sksksk"
	return token == validToken
}
