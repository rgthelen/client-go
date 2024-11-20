package rownd

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "encoding/json"
    "context"
)

var testConfig = &Config{
    APIUrl:    "https://mock-api.local",
    AppKey:    "test-app-key",
    AppSecret: "test-app-secret",
}

func TestClient(t *testing.T) {
    // Setup test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/hub/auth/validate":
            json.NewEncoder(w).Encode(map[string]interface{}{
                "user_id": "rownd-test-user-1",
            })
        case "/hub/users/rownd-test-user-1":
            json.NewEncoder(w).Encode(map[string]interface{}{
                "data": map[string]interface{}{
                    "email": "test@rownd.io",
                    "first_name": "Test",
                    "last_name": "User",
                },
            })
        }
    }))
    defer server.Close()

    // Create client with test server URL
    testConfig.APIUrl = server.URL
    client := NewClient(testConfig)

    // Create context for tests
    ctx := context.Background()

    t.Run("validate token", func(t *testing.T) {
        token, err := testing.GenerateTestToken()
        if err != nil {
            t.Fatalf("Failed to generate test token: %v", err)
        }

        tokenInfo, err := client.ValidateToken(ctx, token)
        if err != nil {
            t.Fatalf("Failed to validate token: %v", err)
        }

        if tokenInfo.UserID == "" {
            t.Error("Expected user_id in token info")
        }
    })

    t.Run("fetch user info", func(t *testing.T) {
        user, err := client.GetUser(ctx, "rownd-test-user-1")
        if err != nil {
            t.Fatalf("Failed to fetch user: %v", err)
        }

        if user.Data["email"] != "test@rownd.io" {
            t.Error("Expected user email to match")
        }
    })
}