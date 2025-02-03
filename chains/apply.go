package chains

import (
	"context"
	"sync"

	"github.com/peterhellberg/llm"
)

const defaultApplyMaxNumberWorkers = 5

// Apply executes the chain for each of the inputs asynchronously.
func Apply(ctx context.Context, c llm.Chain, inputValues []map[string]any, maxWorkers int, options ...llm.ChainOption) ([]map[string]any, error) {
	var (
		inputJobs   = make(chan applyInputJob, len(inputValues))
		resultsChan = make(chan applyResult, len(inputValues))
	)

	if maxWorkers <= 0 {
		maxWorkers = defaultApplyMaxNumberWorkers
	}

	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	for w := 0; w < maxWorkers; w++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case input, ok := <-inputJobs:
					if !ok {
						return
					}

					res, err := Call(ctx, c, input.input, options...)

					resultsChan <- applyResult{
						result: res,
						err:    err,
						i:      input.i,
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	sendApplyInputJobs(inputJobs, inputValues)

	return getApplyResults(ctx, resultsChan, inputValues)
}

type applyInputJob struct {
	input map[string]any
	i     int
}

type applyResult struct {
	result map[string]any
	err    error
	i      int
}

func sendApplyInputJobs(inputJobs chan applyInputJob, inputValues []map[string]any) {
	for i, input := range inputValues {
		inputJobs <- applyInputJob{
			input: input,
			i:     i,
		}
	}

	close(inputJobs)
}

func getApplyResults(ctx context.Context, resultsChan chan applyResult, inputValues []map[string]any) ([]map[string]any, error) {
	results := make([]map[string]any, len(inputValues))

	for range results {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case r := <-resultsChan:
			if r.err != nil {
				return nil, r.err
			}

			results[r.i] = r.result
		}
	}

	return results, nil
}
