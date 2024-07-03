package xerrors

import "io"

func (e *Errors) seek(i *int, offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		*i = int(offset)
	case io.SeekCurrent:
		*i += int(offset)
	case io.SeekEnd:
		*i = (len(e.errs)) + int(offset)
	}
	if *i < 0 || *i >= len(e.errs) {
		return 0, io.EOF
	}
	return int64(*i), nil
}

// Seek implements an [io.Seeker] for the purposes of controlling [errstack.Errors.Next] output.
func (e *Errors) Seek(offset int64, whence int) (int64, error) {
	e.mu.RLock()
	n, err := e.seek(&e.i, offset, whence)
	e.mu.RUnlock()
	return n, err
}

func (e *Errors) next(i *int) error {
	if len(e.errs) == 0 || *i >= len(e.errs) {
		return nil
	}
	err := e.errs[len(e.errs)-1-*i]
	*i++
	return err
}
