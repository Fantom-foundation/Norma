package app

import "sync"

// transactionTracker is a utility to track the number of transactions sent
// per account being part of an application.
type transactionTracker struct {
	mu      sync.Mutex
	sentTxs []uint64
}

// OnSentTransaction should be used to inform the tracker about a transaction
// being send by an account.
func (t *transactionTracker) OnSentTransaction(accountId int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for len(t.sentTxs) <= accountId {
		t.sentTxs = append(t.sentTxs, 0)
	}
	t.sentTxs[accountId] += 1
}

// GetSentTransactionCounts can be used to retrieve the number of transactions
// sent so far per account associated to an application.
func (t *transactionTracker) GetSentTransactionCounts() []uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	res := make([]uint64, len(t.sentTxs))
	copy(res, t.sentTxs)
	return res
}
