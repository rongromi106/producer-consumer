package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Day 3 content for pipelining
	// PipeliningDemo()

	// Day 4 content
	// Non blocking Worker
	// NonBlockWorkerDemo()

	// Day 5 context
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	res := RunPipeline(ctx, 10, 2, 3)
	fmt.Printf("Results from pipeline run %d\n", res)
}
