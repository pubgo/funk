package log

import "fmt"

var _ StdLogger = (*stdLogImpl)(nil)

type stdLogImpl struct {
	log Logger
}

func (s *stdLogImpl) Printf(format string, v ...interface{}) {
	s.log.Info().Msgf(format, v...)
}

func (s *stdLogImpl) Print(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprint(v...))
}

func (s *stdLogImpl) Println(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprint(v...))
}
