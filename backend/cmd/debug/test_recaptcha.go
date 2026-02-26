package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// Load từ environment hoặc dùng hardcode để test
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		secretKey = "6LdvPQIsAAAAAIvC1z3UPeLA7vVwQbi6Wyf2PZd8"
	}

	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")
	if siteKey == "" {
		siteKey = "6LdvPQIsAAAAAG7glbICpFiBR9o5MhboFU4JvxAJ"
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("🔍 KIỂM TRA RECAPTCHA CONFIG")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("\n📌 Site Key (Frontend):\n   %s\n", siteKey)
	fmt.Printf("\n🔐 Secret Key (Backend):\n   %s\n", secretKey)

	// Test với token giả để xem response
	fmt.Println("\n📡 Testing với token rỗng (để xem error response)...")

	verifyURL := "https://www.google.com/recaptcha/api/siteverify"
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", "") // Token rỗng để test

	resp, err := http.PostForm(verifyURL, data)
	if err != nil {
		fmt.Printf("\n❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	fmt.Println("\n📋 Response từ Google:")
	prettyJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(prettyJSON))

	if success, ok := result["success"].(bool); ok && !success {
		if errorCodes, ok := result["error-codes"].([]interface{}); ok {
			fmt.Println("\n⚠️  Error codes:")
			for _, code := range errorCodes {
				fmt.Printf("   - %v\n", code)

				// Giải thích lỗi
				switch code.(string) {
				case "missing-input-secret":
					fmt.Println("     → Secret key thiếu")
				case "invalid-input-secret":
					fmt.Println("     → Secret key không hợp lệ (SAI KEY!)")
				case "missing-input-response":
					fmt.Println("     → Token thiếu (expected - test case)")
				case "invalid-input-response":
					fmt.Println("     → Token không hợp lệ")
				case "timeout-or-duplicate":
					fmt.Println("     → Token đã hết hạn hoặc đã dùng")
				}
			}
		}

		// Nếu chỉ có lỗi missing-input-response thì secret key đúng
		if errorCodes, ok := result["error-codes"].([]interface{}); ok {
			if len(errorCodes) == 1 && errorCodes[0].(string) == "missing-input-response" {
				fmt.Println("\n✅ SECRET KEY HỢP LỆ!")
				fmt.Println("   (Lỗi 'missing-input-response' là expected vì ta test với token rỗng)")
			} else {
				fmt.Println("\n❌ SECRET KEY KHÔNG HỢP LỆ hoặc có vấn đề khác!")
			}
		}
	} else {
		fmt.Println("\n✅ Response success!")
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("💡 HƯỚNG DẪN:")
	fmt.Println("   1. Nếu có lỗi 'invalid-input-secret' → Key sai")
	fmt.Println("   2. Nếu chỉ có 'missing-input-response' → Key đúng!")
	fmt.Println("   3. Kiểm tra keys tại:")
	fmt.Println("      https://www.google.com/recaptcha/admin")
	fmt.Println("═══════════════════════════════════════════════════════════\n")
}
