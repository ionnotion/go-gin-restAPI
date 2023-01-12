package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	ctx.JSON(405, gin.H{"message":"Method not allowed"})
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

func removeIndex(arr []Language, index int) []Language {
    ret := make([]Language, 0)
    ret = append(ret, arr[:index]...)
    return append(ret, arr[index+1:]...)
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
	languages := []Language{newLanguage}

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message":"Hello Go Developers"})
	})

	router.GET("/languages", func(ctx *gin.Context) {
		// fmt.Println(newLanguage)
		b, err := json.MarshalIndent(languages,""," ")
		// fmt.Println(string(b))
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		ctx.Data(http.StatusOK, gin.MIMEJSON, b)
	})

	router.GET("/palindrome",func(ctx *gin.Context) {
		text, _ := ctx.GetQuery("text")

		if text == "" {
			handleOtherRoutes(ctx)
			return
		}

		result := "Not Palindrome"
		code := 400
		if checkPalindrome(text) {
			result = "Palindrome"
			code = 200
		}

		ctx.JSON(code, gin.H{"message":result})
	})

	router.POST("/language",func(ctx *gin.Context) {
		var language Language
		ctx.BindJSON(&language)

		languages = append(languages, language)
		b, err := json.MarshalIndent(language,""," ")
		// fmt.Println(string(b))
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		ctx.Data(201, gin.MIMEJSON, b)
	})

	router.GET("/language/:id",func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		
		if err !=nil {
			handleOtherRoutes(ctx)
			return
		}

		language := languages[id]

		b, err := json.MarshalIndent(language,""," ")

		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		ctx.Data(http.StatusOK, gin.MIMEJSON, b)
	})

	router.PATCH("/language/:id",func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		
		if err !=nil {
			handleOtherRoutes(ctx)
			return
		}

		var language Language
		ctx.BindJSON(&language)

		languages[id] = language
		b, err := json.MarshalIndent(language,""," ")

		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		ctx.Data(http.StatusOK, gin.MIMEJSON, b)
	})

	router.DELETE("/language/:id",func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		
		if err !=nil {
			handleOtherRoutes(ctx)
			return
		}

		languages = removeIndex(languages,id)

		ctx.JSON(http.StatusOK, gin.H{"message":"deleted"})
	})

	router.NoMethod(handleOtherRoutes)
	router.NoRoute(handleOtherRoutes)

	router.HandleMethodNotAllowed = true
	router.Run(":8888")
}
