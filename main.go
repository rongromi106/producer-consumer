package main

import (
	"context"
	"fmt"
)

func main() {
	// Day 3 content for pipelining
	// PipeliningDemo()

	// Day 4 content
	// Non blocking Worker
	// NonBlockWorkerDemo()

	// Day 5 context
	res := RunPipeline(context.Background(), 10, 2, 3)
	fmt.Printf("Results from pipeline run %d\n", res)
}
