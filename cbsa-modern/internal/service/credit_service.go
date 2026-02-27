package service

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

const numAgencies = 5

// CreditService replaces the five COBOL programs CRDTAGY1.cbl through
// CRDTAGY5.cbl. Each simulated a separate credit agency with a random delay
// (0-3s) and random score (1-999). The COBOL version required CICS RUN
// TRANSID with PUT/GET CONTAINER for each agency; here goroutines and
// errgroup replace the entire async mechanism.
type CreditService struct{}

func NewCreditService() *CreditService {
	return &CreditService{}
}

func (s *CreditService) CheckCredit(ctx context.Context, customerNumber string) (int, error) {
	g, ctx := errgroup.WithContext(ctx)
	scores := make([]int, numAgencies)

	for i := 0; i < numAgencies; i++ {
		i := i
		g.Go(func() error {
			delay := time.Duration(rand.Intn(3000)) * time.Millisecond
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
			scores[i] = rand.Intn(999) + 1
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	total := 0
	for _, score := range scores {
		total += score
	}
	return total / numAgencies, nil
}
