// Note: the documentation / examples for error handling in the go-flags
// library is currently a combination of lacking and wrong:
// * https://github.com/jessevdk/go-flags/issues/306
// * https://github.com/jessevdk/go-flags/issues/361
// * https://github.com/jessevdk/go-flags/issues/377
//
// This package attempts to make working with CLI errors easier (and
// easier to maintain).
package flagswrap

import (
	"errors"
	"fmt"

	"github.com/jessevdk/go-flags"
)

var (
	ErrGoFlagsErrorWrapper = errors.New("couldn't wrap go-flags error")
	// These are the errors we've noticed are handled nicely in the go-flags
	// library:
	goFlagsVerboseErrors = map[flags.ErrorType]bool{
		flags.ErrHelp:              true,
		flags.ErrCommandRequired:   true,
		flags.ErrUnknownCommand:    true,
		flags.ErrUnknownFlag:       true,
		flags.ErrNoArgumentForBool: true,
		flags.ErrInvalidChoice:     true,
	}
	// These are the go-flags errors that are swallowed silently and require some
	// logging / feedback to users:
	goFlagsSilentErrors = map[flags.ErrorType]bool{
		flags.ErrUnknown:        true,
		flags.ErrDuplicatedFlag: true,
		// The remaining errors haven't had their verbosity levels confirmed;
		// if you see double-logging happening in the CLI, just move the error
		// in question from below up to the goFlagsVerboseErrors map.
		flags.ErrExpectedArgument: true,
		flags.ErrUnknownGroup:     true,
		flags.ErrMarshal:          true,
		flags.ErrRequired:         true,
		flags.ErrShortNameTooLong: true,
		flags.ErrTag:              true,
		flags.ErrInvalidTag:       true,
	}
)

type Error struct {
	wrapped     error
	convertErr  error
	flagErr     *flags.Error
	flagErrType flags.ErrorType
}

func WrapError(err error) *Error {
	if err == nil {
		return nil
	}
	cliErr := &Error{
		wrapped: err,
	}
	flagsErr, ok := err.(*flags.Error)
	if !ok {
		cliErr.convertErr = ErrGoFlagsErrorWrapper
		return cliErr
	}
	cliErr.flagErr = flagsErr
	cliErr.flagErrType = flagsErr.Type
	return cliErr
}

func (e *Error) IsHelp() bool {
	return e.flagErrType == flags.ErrHelp
}

func (e *Error) IsSilent() bool {
	exists, ok := goFlagsSilentErrors[e.flagErrType]
	if e.convertErr == nil && ok && exists {
		return true
	}
	return false
}

func (e *Error) IsVerbose() bool {
	exists, ok := goFlagsVerboseErrors[e.flagErrType]
	if e.convertErr == nil && ok && exists {
		return true
	}
	return false
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v (%v) - %v", e.wrapped, e.convertErr, e.flagErrType)
}

func (e *Error) Unwrap() error {
	return e.wrapped
}
