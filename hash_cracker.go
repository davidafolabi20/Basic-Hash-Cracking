package main

import (
    "bufio"
    "fmt"
    "os"
    "runtime"
    "sync"
    "time"
    "golang.org/x/crypto/bcrypt"
)

// Result represents a successful hash match
type Result struct {
    sequence string
    found    bool
}

// Worker tests sequences against the hash
func worker(id int, sequences <-chan string, results chan<- Result, targetHash []byte, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for sequence := range sequences {
        // Check if this sequence matches the hash
        err := bcrypt.CompareHashAndPassword(targetHash, []byte(sequence))
        if err == nil {
            results <- Result{sequence: sequence, found: true}
            return
        }
        
        if id == 0 && time.Now().Unix()%5 == 0 {
            fmt.Printf("Worker %d: Still checking sequences...\n", id)
        }
    }
}

func main() {
    // Use all available CPU cores
    numCPU := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPU)
    
    targetHash := []byte("$2y$05$tJ5qkcBGrjiRfZZAlkSsP.kcVStH7oCzsery3nN1sgXk02xThNck6")
    
    // Open sequences file
    file, err := os.Open("sequences.txt")
    if err != nil {
        fmt.Printf("Error opening sequences file: %v\n", err)
        return
    }
    defer file.Close()

    // Channel setup
    sequences := make(chan string, numCPU*100)
    results := make(chan Result, numCPU)
    
    // Start worker goroutines
    var wg sync.WaitGroup
    fmt.Printf("Starting %d workers...\n", numCPU)
    for i := 0; i < numCPU; i++ {
        wg.Add(1)
        go worker(i, sequences, results, targetHash, &wg)
    }

    // Start sequence feeder goroutine
    go func() {
        scanner := bufio.NewScanner(file)
        count := 0
        for scanner.Scan() {
            sequences <- scanner.Text()
            count++
            if count%10000 == 0 {
                fmt.Printf("Loaded %d sequences...\n", count)
            }
        }
        close(sequences)
    }()

    // Start timer
    start := time.Now()

    // Wait for either a match or completion
    go func() {
        wg.Wait()
        close(results)
    }()

    // Check results
    found := false
    var password string
    
    for result := range results {
        if result.found {
            found = true
            password = result.sequence
            break
        }
    }

    elapsed := time.Since(start)
    
    if found {
        fmt.Println("\n==================================================")
        fmt.Println("MATCH FOUND!")
        fmt.Printf("Password: %s\n", password)
        fmt.Printf("Flag format: 02h{%s}\n", password)
        fmt.Println("==================================================")
    } else {
        fmt.Println("\nNo matches found.")
    }
    
    fmt.Printf("\nTime elapsed: %s\n", elapsed)
}