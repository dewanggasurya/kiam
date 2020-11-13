package kiam

type Claim interface {
	SessionID() string
	SetSessionID(id string) Claim
}

type BaseClaim struct {
	ID    string
	RefID string
}

func (bc *BaseClaim) SessionID() string {
	return bc.ID
}

func (bc *BaseClaim) SetSessionID(id string) Claim {
	bc.ID = id
	return bc
}
