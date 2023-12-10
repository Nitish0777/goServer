package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func main() {
	r := gin.Default()

	// Define your routes
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello MarkUp.ai! I like to work with you!")
	})
	r.POST("/process-single", processSingle)
	r.POST("/process-concurrent", processConcurrent)

	// Start the Gin server
	fmt.Println("Listening on port 8000...")
	r.Run(":8000")
}

func processSingle(c *gin.Context) {
	var payload RequestPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	startTime := time.Now()
	sortedArrays := make([][]int, len(payload.ToSort))

	for i, arr := range payload.ToSort {
		sortedArrays[i] = sortSequential(arr)
	}

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       time.Since(startTime).Nanoseconds(),
	}

	c.JSON(http.StatusOK, response)
}

func processConcurrent(c *gin.Context) {
	var payload RequestPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	sortedArrays := make([][]int, len(payload.ToSort))
	ch := make(chan struct {
		index int
		arr   []int
	})

	for i, arr := range payload.ToSort {
		wg.Add(1)
		go sortConcurrent(arr, ch, i, &wg, sortedArrays)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		sortedArrays[result.index] = result.arr
	}

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       time.Since(startTime).Nanoseconds(),
	}

	c.JSON(http.StatusOK, response)
}

func sortSequential(arr []int) []int {
	sortedArr := make([]int, len(arr))
	copy(sortedArr, arr)
	sort.Ints(sortedArr)
	return sortedArr
}

func sortConcurrent(arr []int, ch chan struct {
	index int
	arr   []int
}, index int, wg *sync.WaitGroup, sortedArrays [][]int) {
	defer wg.Done()
	sortedArr := sortSequential(arr)
	ch <- struct {
		index int
		arr   []int
	}{index, sortedArr}
}
