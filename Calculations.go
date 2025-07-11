package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Picker struct {
	items       []Item
	usedCombos  map[string]struct{}
	totalCombos int
	mu          sync.Mutex
}

var totalCount uint64

func NewPicker(items []Item) *Picker {
	n := len(items)
	totalCombos := combinationsCount(n, 0)
	return &Picker{
		items:       items,
		usedCombos:  make(map[string]struct{}),
		totalCombos: totalCombos,
	}
}

func combinationsCount(n, k int) int {
	if k > n {
		return 0
	}
	num := 1
	den := 1
	for i := 0; i < k; i++ {
		num *= n - i
		den *= i + 1
	}
	return num / den
}

func (p *Picker) PickUnique6Parallel(workers int) ([]Item, error) {
	n := len(p.items)
	if n < 6 {
		return nil, fmt.Errorf("not enough items")
	}

	p.mu.Lock()
	if len(p.usedCombos) >= p.totalCombos {
		p.mu.Unlock()
		return nil, fmt.Errorf("all combinations have been used")
	}
	p.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type result struct {
		selection []Item
		key       string
	}

	resultCh := make(chan result)
	var wg sync.WaitGroup

	var totalCount uint64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Randomly pick 6 unique indices
					indices := r.Perm(n)[:6]

					ids := make([]int, 6)
					selection := make([]Item, 6)
					for j, idx := range indices {
						idInt, err := strconv.Atoi(p.items[idx].ID)
						if err != nil {
							// If invalid ID, skip
							continue
						}
						ids[j] = idInt
						selection[j] = p.items[idx]
					}
					sort.Ints(ids)

					keyParts := make([]string, 6)
					for j, id := range ids {
						keyParts[j] = strconv.Itoa(id)
					}
					key := strings.Join(keyParts, "-")

					// Increment attempt counter
					newCount := atomic.AddUint64(&totalCount, 1)

					// Print progress every 1000 tries
					if workerID == 1 && newCount%1000 == 0 {
						// Create a string of item names for display
						itemNames := make([]string, len(selection))
						for j, item := range selection {
							itemNames[j] = item.Name
						}
						itemStr := strings.Join(itemNames, ", ")

						line := fmt.Sprintf("[Worker %d] Attempts: %d | Last combo: %s", workerID, newCount, itemStr)

						const maxLineLength = 160
						if len(line) < maxLineLength {
							line += strings.Repeat(" ", maxLineLength-len(line))
						}
						fmt.Printf("\r%s", line)
					}

					// Check if this combination was already used
					p.mu.Lock()
					_, exists := p.usedCombos[key]
					p.mu.Unlock()

					if !exists {
						resultCh <- result{selection: selection, key: key}
						return
					}
				}
			}
		}(i + 1)
	}

	// Wait for the first successful combination
	res := <-resultCh

	// Mark the combination as used
	p.mu.Lock()
	p.usedCombos[res.key] = struct{}{}
	p.mu.Unlock()

	// Cancel all workers
	cancel()

	wg.Wait()

	fmt.Println() // Clean up the progress line

	return res.selection, nil
}

func CalculateSingleChamp(champion *Champion, items []Item, workers int) {
	if champion == nil {
		fmt.Println("No champion provided for calculation.")
		return
	}

	fmt.Printf("Calculating stats for champion: %s\n", champion.Name)
	fmt.Printf("Title: %s\n", champion.Title)

	baseItems := items

	if len(baseItems) < 6 {
		fmt.Println("Not enough items")
		return
	}

	picker := NewPicker(baseItems)

	selection, err := picker.PickUnique6Parallel(workers)
	if err != nil {
		fmt.Println("Error picking unique combinations: ", err)
		return
	}

	fmt.Println("Selected Items: ")
	for _, item := range selection {
		fmt.Printf(" %s (Total Gold: %d)\n", item.Name, item.Gold.Total)
	}

}

func SearchChampionsByName(champions []Champion, name string) *Champion {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter champion name: ")
	var champ *Champion
	if name == "Search" {
		for {
			input, _ := reader.ReadString('\n')
			if input == "" || input == "\n" {
				fmt.Println("No input provided: please enter a champion name.")
				continue // Prompt again for input
			}
			input = input[:len(input)-1] // Remove the newline character
			champ = FindChampByName(champions, input)
			if champ != nil {
				fmt.Printf("Champion found: %s, Title: %s\n", champ.Name, champ.Title)
				return champ // Exit the loop if a champion is found
			} else {
				fmt.Printf("Champion '%s' not found. Please try again.\n", input)
			}
		}
	} else {
		return FindChampByName(champions, name)
	}
}

// combinations generates all combinations of k indices out of n elements.
// It sends each combination over a channel.
func combinations(n, k int) <-chan []int {
	ch := make(chan []int, 1000)

	go func() {
		defer close(ch)
		indices := make([]int, k)
		for i := 0; i < k; i++ {
			indices[i] = i
		}
		for {
			// Send a copy to the channel
			comb := make([]int, k)
			copy(comb, indices)
			ch <- comb

			// Find the rightmost index to increment
			i := k - 1
			for i >= 0 && indices[i] == i+n-k {
				i--
			}
			if i < 0 {
				return
			}
			indices[i]++
			for j := i + 1; j < k; j++ {
				indices[j] = indices[j-1] + 1
			}
		}
	}()

	return ch
}

func ListAllItemCombos(items []Item) {
	n := len(items)
	k := 6

	if n < k {
		fmt.Println("Not enough items for combinations")
		return
	}

	count := 0
	for indices := range combinations(n, k) {
		count++
		fmt.Printf("Combination %d:\n", count)
		for _, idx := range indices {
			item := items[idx]
			fmt.Println(" %s (Total Gold: %d)\n", item.Name, item.Gold.Total)
		}
		fmt.Println()

		fmt.Printf("Processed %d combinations.\n", count)
	}
}

func ProcessAllCombinationsParallel(items []Item, k int, workerCount int) {
	n := len(items)
	if n < k {
		fmt.Println("Not enough items to make combinations.")
		return
	}

	combCh := combinations(n, k)

	var wg sync.WaitGroup

	// Launch workers
	var count int
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for indices := range combCh {
				// Example processing: print
				// fmt.Printf("[Worker %d] Combination:\n", workerID)
				for _, idx := range indices {
					item := items[idx]
					fmt.Printf("  %s (Total Gold: %d) ", item.Name, item.Gold.Total)
				}
				count++
				fmt.Println()
				fmt.Printf("Generated %d combinations", count)
				fmt.Println()
			}
		}(i + 1)
	}

	wg.Wait()
}

func combinationsBatch(n, k, batchSize int) <-chan [][]int {
	ch := make(chan [][]int, 10)

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

func ProcessAllCombinationsBatched(items []Item, k int, workerCount, batchSize int) {
	n := len(items)
	if n < k {
		fmt.Println("Not enough items to make combinations.")
		return
	}
	combCh := combinationsBatch(n, k, batchSize)
	var wg sync.WaitGroup
	var count int64
	var mu sync.Mutex

	startTime := time.Now()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batch := range combCh {
				for _, indices := range batch {
					// Example processing: print
					line := ""
					// fmt.Printf("[Worker %d] Combination:\n", workerID)
					for _, idx := range indices {
						item := items[idx]
						itemString := fmt.Sprintf("  %s (Total Gold: %d) ", item.Name, item.Gold.Total)
						line += itemString
					}

					mu.Lock()
					count++
					currentCount := count
					mu.Unlock()

					if currentCount%1000 == 0 {
						elapsed := time.Since(startTime).Seconds()
						rate := float64(currentCount) / elapsed
						fmt.Printf("\rGenerated %d combinations | Rate: %.0f/sec | %s",
							currentCount, rate, line)
					}

					_ = indices
				}
			}
		}(i + 1)
	}

	wg.Wait()

	// Print final stats
	totalTime := time.Since(startTime).Seconds()
	finalRate := float64(count) / totalTime
	fmt.Printf("\nCompleted %d combinations in %.2f seconds (%.0f/sec)\n",
		count, totalTime, finalRate)
}
