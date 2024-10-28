package log_internal

import "fmt"

var _ StdLogger = (*stdLogImpl)(nil)

func NewStd(log Logger) StdLogger {
	return &stdLogImpl{log: log.WithCallerSkip(1)}
}

type stdLogImpl struct {
	log Logger
}

func (s *stdLogImpl) Printf(format string, v ...interface{}) {
	s.log.Info().Msgf(format, v...)
}

func (s *stdLogImpl) Print(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprint(v...))
}

func (s *stdLogImpl) Log(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprint(v...))
}

func (s *stdLogImpl) Logf(format string, v ...interface{}) {
	s.log.Info().Msgf(format, v...)
}

func (s *stdLogImpl) Println(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprint(v...))
}
