package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Optimized batch generator - same as your original but with better buffering
func combinationsBatch(n, k, batchSize int) <-chan [][]int {
	ch := make(chan [][]int, 100) // Increased buffer size

	go func() {
		defer close(ch)
		indices := make([]int, k)
		for i := 0; i < k; i++ {
			indices[i] = i
		}

		batch := make([][]int, 0, batchSize)

		for {
			comb := make([]int, k)
			copy(comb, indices)
			batch = append(batch, comb)

			// Send the batch when it's full
			if len(batch) == batchSize {
				ch <- batch
				batch = make([][]int, 0, batchSize)
			}

			// Increment indices
			i := k - 1
			for i >= 0 && indices[i] == i+n-k {
				i--
			}
			if i < 0 {
				break
			}
			indices[i]++
			for j := i + 1; j < k; j++ {
				indices[j] = indices[j-1] + 1
			}
		}

		// Send remaining combinations
		if len(batch) > 0 {
			ch <- batch
		}
	}()

	return ch
}

// Optimized version that removes bottlenecks
func ProcessAllCombinationsOptimized(items []Item, k int, workerCount, batchSize int) {
	n := len(items)
	if n < k {
		fmt.Println("Not enough items to make combinations.")
		return
	}

	combCh := combinationsBatch(n, k, batchSize)
	var wg sync.WaitGroup
	var count int64

	startTime := time.Now()

	// Start a separate goroutine for progress reporting
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentCount := atomic.LoadInt64(&count)
				elapsed := time.Since(startTime).Seconds()
				rate := float64(currentCount) / elapsed
				fmt.Printf("\rGenerated %d combinations | Rate: %.0f/sec", currentCount, rate)
			case <-done:
				return
			}
		}
	}()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batch := range combCh {
				for _, indices := range batch {
					// Do your actual processing here
					// For now, just counting
					atomic.AddInt64(&count, 1)

					// Remove the expensive string operations and printing
					_ = indices
				}
			}
		}(i + 1)
	}

	wg.Wait()
	close(done)

	// Print final stats
	totalTime := time.Since(startTime).Seconds()
	finalRate := float64(count) / totalTime
	fmt.Printf("\nCompleted %d combinations in %.2f seconds (%.0f/sec)\n",
		count, totalTime, finalRate)
}

// If you need to do actual work with the items, here's a more realistic version
func ProcessAllCombinationsWithWork(items []Item, k int, workerCount, batchSize int) {
	n := len(items)
	if n < k {
		fmt.Println("Not enough items to make combinations.")
		return
	}

	combCh := combinationsBatch(n, k, batchSize)
	var wg sync.WaitGroup
	var count int64

	startTime := time.Now()

	// Progress reporting
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentCount := atomic.LoadInt64(&count)
				elapsed := time.Since(startTime).Seconds()
				rate := float64(currentCount) / elapsed
				fmt.Printf("\rGenerated %d combinations | Rate: %.0f/sec", currentCount, rate)
			case <-done:
				return
			}
		}
	}()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batch := range combCh {
				for _, indices := range batch {
					// Do your actual processing here
					// Example: calculate total gold for this combination
					var totalGold int
					for _, idx := range indices {
						totalGold += items[idx].Gold.Total
					}

					// You can store results, check conditions, etc.
					// But avoid expensive operations like string formatting and printing

					atomic.AddInt64(&count, 1)
				}
			}
		}(i + 1)
	}

	wg.Wait()
	close(done)

	// Print final stats
	totalTime := time.Since(startTime).Seconds()
	finalRate := float64(count) / totalTime
	fmt.Printf("\nCompleted %d combinations in %.2f seconds (%.0f/sec)\n",
		count, totalTime, finalRate)
}

// Ultra-fast version for counting only
func CountAllCombinations(n, k int) {
	startTime := time.Now()

	var count int64
	batchSize := 10000

	combCh := combinationsBatch(n, k, batchSize)
	var wg sync.WaitGroup

	// Progress reporting
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentCount := atomic.LoadInt64(&count)
				elapsed := time.Since(startTime).Seconds()
				rate := float64(currentCount) / elapsed
				fmt.Printf("\rGenerated %d combinations | Rate: %.0f/sec", currentCount, rate)
			case <-done:
				return
			}
		}
	}()

	// Use more workers for pure counting
	workerCount := 8
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range combCh {
				// Just add the batch size to count
				atomic.AddInt64(&count, int64(len(batch)))
			}
		}()
	}

	wg.Wait()
	close(done)

	totalTime := time.Since(startTime).Seconds()
	finalRate := float64(count) / totalTime
	fmt.Printf("\nCompleted %d combinations in %.2f seconds (%.0f/sec)\n",
		count, totalTime, finalRate)
}

// Replace your original function with this optimized version
func ProcessAllCombinationsBatchedOptimized(items []Item, k int, workerCount, batchSize int) {
	ProcessAllCombinationsOptimized(items, k, workerCount, batchSize)
}
