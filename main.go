package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Relation struct {
	InfluencedBy []string `json:"influenced-by"`
	Influences   []string `json:"influences"`
}

type Language struct {
	Language       string   `json:"language"`
	Appeared       int      `json:"appeared"`
	Created        []string `json:"created"`
	Functional     bool     `json:"functional"`
	ObjectOriented bool     `json:"object-oriented"`
	Relation       Relation `json:"relation"`
}

func handleOtherRoutes(ctx *gin.Context) {
	ctx.String(405,"Method not allowed")
}

func reverser(text string) string {
	var sb strings.Builder
	runes := []rune(text)
	for i:= len(runes)-1; i >= 0; i-- {
		sb.WriteRune(runes[i])
	}
	fmt.Println(text, sb.String())
	return sb.String()
}

func checkPalindrome(text string) bool {
	return strings.EqualFold(text,reverser(text))
}

func main() {
	router := gin.Default()
	newLanguage := Language{
		Language:       "C",
		Appeared:       1972,
		Created:        []string{"Dennis Ritchie"},
		Functional:     true,
		ObjectOriented: false,
		Relation: Relation{
			[]string{"B",
				"ALGOL 68",
				"Assembly",
				"FORTRAN"},
			[]string{"C++",
				"Objective-C",
				"C#",
				"Java",
				"JavaScript",
				"PHP",
				"Go"},
		},
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello Go Developers")
	})

	router.GET("/language", func(ctx *gin.Context) {
		fmt.Println(newLanguage)
		b, err := json.MarshalIndent(newLanguage,""," ")
		fmt.Println(string(b))
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		ctx.Data(http.StatusOK, gin.MIMEJSON, b)
	})

	router.GET("/palindrome",func(ctx *gin.Context) {
		text, _ := ctx.GetQuery("text")

		result := "Not Palindrome"
		code := 400
		if checkPalindrome(text) {
			result = "Palindrome"
			code = 200
		}

		ctx.String(code, result)
	})

	router.NoMethod(handleOtherRoutes)
	router.NoRoute(handleOtherRoutes)

	router.HandleMethodNotAllowed = true
	router.Run(":8888")
}
