package domain

import "errors"

var (
	ErrAuctionAlreadyStarted = errors.New("auction already started")
	ErrAuctionNotStarted     = errors.New("auction is not started")
	ErrAuctionCancelled      = errors.New("auction cancelled")
	ErrAuctionCompleted      = errors.New("auction completed")
	ErrBidTooLow             = errors.New("bid amount is too low")
	ErrTimeoutTooEarly       = errors.New("timeout is too early")
	ErrSellerIDInvalid       = errors.New("seller ID is invalid")
	ErrStartTimeInvalid      = errors.New("start time is invalid")
	ErrEndTimeInvalid        = errors.New("end time is invalid")
)
