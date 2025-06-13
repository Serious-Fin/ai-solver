package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type problem struct {
	Id             int
	Title          string
	Description    string
	TestScenarions []testScenario
}

type testScenario struct {
	Input  string
	Output string
}

var problems = []problem{
	{
		Id:          1,
		Title:       "Largest Sum",
		Description: "Find the largest sum foo bar",
		TestScenarions: []testScenario{
			{
				Input:  "[1, 2, 3]",
				Output: "6",
			},
			{
				Input:  "[]",
				Output: "0",
			},
		},
	},
	{
		Id:          2,
		Title:       "Maximum sub array",
		Description: "Find the largest sub-array in the provided array",
		TestScenarions: []testScenario{
			{
				Input:  "[-8, -9, 10, 80, 30, -8]",
				Output: "[10, 80, 30]",
			},
			{
				Input:  "[]",
				Output: "0",
			},
		},
	},
}

func main() {
	router := gin.Default()
	router.GET("/problems", getProblems)
	router.GET("/problems/:id", getProblemById)

	router.Run("localhost:8080")
}

func getProblems(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, problems)
}

func getProblemById(c *gin.Context) {
	id := c.Param("id")

	// string to int
	stringifiedId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid problem id provided. Make sure it is a positive integer"})
		return
	}

	for _, problem := range problems {
		if problem.Id == stringifiedId {
			c.IndentedJSON(http.StatusOK, problem)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Problem with id not found"})
}
