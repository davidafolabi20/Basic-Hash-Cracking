package main

import (
    "bufio"
    "fmt"
    "os"
    "runtime"
    "sync"
)

// Chunk represents a portion of the prime number to process
type Chunk struct {
    start int
    data  string
}

// Worker processes chunks of the prime number
func worker(id int, chunks <-chan Chunk, results chan<- []string, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for chunk := range chunks {
        sequences := make([]string, 0)
        data := chunk.data
        
        // Process sequences in this chunk
        for i := 0; i < len(data)-9; i++ { // -9 for 10-digit sequences
            sequence := data[i : i+10]
            // Check if sequence contains only digits
            if isDigitSequence(sequence) {
                sequences = append(sequences, sequence)
            }
        }
        
        if len(sequences) > 0 {
            results <- sequences
        }
    }
}

// isDigitSequence checks if a string contains only digits
func isDigitSequence(s string) bool {
    for _, r := range s {
        if r < '0' || r > '9' {
            return false
        }
    }
    return true
}

func main() {
    // Use all available CPU cores
    numCPU := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPU)
    
    // Open input file
    inputFile := "prime.txt" // Replace with your prime number file
    file, err := os.Open(inputFile)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()

    // Open output file
    outFile, err := os.Create("sequences.txt")
    if err != nil {
        fmt.Printf("Error creating output file: %v\n", err)
        return
    }
    defer outFile.Close()
    
    // Create buffered writer for better performance
    writer := bufio.NewWriter(outFile)
    defer writer.Flush()

    // Channel setup
    chunkSize := 1024 * 1024 // 1MB chunks
    chunks := make(chan Chunk, numCPU)
    results := make(chan []string, numCPU)
    
    // Start worker goroutines
    var wg sync.WaitGroup
    for i := 0; i < numCPU; i++ {
        wg.Add(1)
        go worker(i, chunks, results, &wg)
    }

    // Start result processor goroutine
    done := make(chan bool)
    go func() {
        totalSequences := 0
        for sequences := range results {
            for _, seq := range sequences {
                writer.WriteString(seq + "\n")
                totalSequences++
            }
        }
        fmt.Printf("Total sequences found: %d\n", totalSequences)
        done <- true
    }()

    // Read and process file in chunks
    reader := bufio.NewReader(file)
    buffer := make([]byte, chunkSize)
    leftover := ""
    
    fmt.Println("Processing prime number file...")
    
    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            data := leftover + string(buffer[:n])
            
            // Keep the last 9 characters for the next chunk
            if len(data) > 9 {
                chunks <- Chunk{
                    data: data[:len(data)-9],
                }
                leftover = data[len(data)-9:]
            } else {
                leftover = data
            }
        }
        
        if err != nil {
            break
        }
    }

    // Process any remaining data
    if len(leftover) >= 10 {
        chunks <- Chunk{data: leftover}
    }

    // Close channels and wait for completion
    close(chunks)
    wg.Wait()
    close(results)
    <-done

    fmt.Println("Sequence extraction complete!")
}