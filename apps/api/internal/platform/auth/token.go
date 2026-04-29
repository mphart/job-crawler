package authx

import (
    "encoding/base64"
    "fmt"
    "strings"
)

func IssueToken(userID string) string { return base64.StdEncoding.EncodeToString([]byte("uid:" + userID)) }
func ParseToken(token string) (string, error) {
    b, err := base64.StdEncoding.DecodeString(token)
    if err != nil { return "", err }
    s := string(b)
    if !strings.HasPrefix(s, "uid:") { return "", fmt.Errorf("invalid token") }
    return strings.TrimPrefix(s, "uid:"), nil
}
