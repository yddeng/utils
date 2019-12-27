package transaction

import "sync/atomic"

type Transaction struct {
	trans    []*transaction
	callback func(susses bool)
	started  int32
}

type transaction struct {
	do       func() bool
	rollback func() bool
	susses   bool
}

func NewTransaction() *Transaction {
	return &Transaction{
		trans:   []*transaction{},
		started: 0,
	}
}

func (this *Transaction) Push(do, rollback func() bool) {
	this.trans = append(this.trans, &transaction{
		do:       do,
		rollback: rollback,
		susses:   false,
	})
}

func (this *Transaction) rollback() {
	for _, t := range this.trans {
		if t.susses {
			t.rollback()
		}
	}
}

func (this *Transaction) Do(callback func(susses bool)) {
	if !atomic.CompareAndSwapInt32(&this.started, 0, 1) {
		return
	}

	for _, t := range this.trans {
		t.susses = t.do()
		if !t.susses {
			this.rollback()
			callback(false)
			return
		}
	}

	callback(true)
}
