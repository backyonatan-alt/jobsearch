package mail

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// UnsubToken signs a user id for one-click unsubscribe links. Keyed off a
// long-lived server secret so links in old emails keep working.
func UnsubToken(secret string, userID int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "unsub:%d", userID)
	return hex.EncodeToString(mac.Sum(nil))[:32]
}

func ValidUnsubToken(secret string, userID int64, token string) bool {
	return hmac.Equal([]byte(UnsubToken(secret, userID)), []byte(token))
}
