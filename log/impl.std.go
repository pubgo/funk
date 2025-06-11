package log

import "fmt"

var _ StdLogger = (*stdLogImpl)(nil)

type stdLogImpl struct {
	log Logger
}

func (s *stdLogImpl) Printf(format string, v ...any) {
	s.log.Info().Msgf(format, v...)
}

func (s *stdLogImpl) Print(v ...any) {
	s.log.Info().Msg(fmt.Sprint(v...))
}

func (s *stdLogImpl) Log(v ...any) {
	s.log.Info().Msg(fmt.Sprint(v...))
}

func (s *stdLogImpl) Logf(format string, v ...any) {
	s.log.Info().Msgf(format, v...)
}

func (s *stdLogImpl) Println(v ...any) {
	s.log.Info().Msg(fmt.Sprint(v...))
}
