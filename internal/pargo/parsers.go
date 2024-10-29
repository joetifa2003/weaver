package pargo

import (
	"fmt"
	"io"
	"strings"

	"github.com/joetifa2003/weaver/internal/pargo/lexer"
)

type State struct {
	source string
	tokens []lexer.Token
	pos    int
}

func (s *State) consume() (lexer.Token, error) {
	if s.pos >= len(s.tokens) {
		return nil, io.EOF
	}

	val := s.tokens[s.pos]
	s.pos++

	return val, nil
}

func (s *State) peek() (lexer.Token, error) {
	if s.pos >= len(s.tokens) {
		return nil, io.EOF
	}

	val := s.tokens[s.pos]

	return val, nil
}

type Parser[T any] func(state State) (T, State, error)

func Exactly(s string) Parser[string] {
	return func(state State) (string, State, error) {
		old := state

		tok, err := state.consume()
		if err != nil {
			return "", old, err
		}

		if tok.String() != s {
			return "", old, NewParseError(state.source, s, tok)
		}

		return s, state, nil
	}
}

func Except(s string) Parser[string] {
	return func(state State) (string, State, error) {
		old := state

		tok, err := state.consume()
		if err != nil {
			return "", old, err
		}

		if tok.String() == s {
			return "", old, NewParseError(state.source, fmt.Sprintf("not %s", s), tok)
		}

		return tok.String(), state, nil
	}
}

func OneOf[T any](parsers ...Parser[T]) Parser[T] {
	return func(state State) (T, State, error) {
		var res T
		var err error

		for _, p := range parsers {
			res, state, err = p(state)
			if err == nil {
				return res, state, nil
			}
		}

		return zero[T](), state, err
	}
}

func Map[T, U any](p Parser[T], f func(T) (U, error)) Parser[U] {
	return func(state State) (U, State, error) {
		res, state, err := p(state)
		if err != nil {
			return zero[U](), state, err
		}

		mapped, err := f(res)
		if err != nil {
			return zero[U](), state, err
		}

		return mapped, state, err
	}
}

func ManySep[T any, S any](p Parser[T], separator Parser[S]) Parser[[]T] {
	return Sequence2(
		Many(
			Sequence2(
				p,
				separator,
				func(a T, _ S) T {
					return a
				},
			),
		),
		p,
		func(a []T, b T) []T {
			return append(a, b)
		},
	)
}

func Sequence[T any, O any](mapper func(T) O, psT Parser[T]) Parser[O] {
	return func(state State) (O, State, error) {
		resT, newState, err := psT(state)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT)

		return res, newState, nil
	}
}

func Many[T any](p Parser[T]) Parser[[]T] {
	return func(state State) ([]T, State, error) {
		var res []T
		var err error
		var r T

		for {
			r, state, err = p(state)
			if err != nil {
				return res, state, nil
			}
			res = append(res, r)
		}
	}
}

func Some[T any](p Parser[T]) Parser[[]T] {
	return func(state State) ([]T, State, error) {
		first, state, error := p(state)
		if error != nil {
			return zero[[]T](), state, error
		}

		res := []T{first}

		var other T

		for {
			other, state, error = p(state)
			if error != nil {
				return res, state, nil
			}
			res = append(res, other)
		}
	}
}

func Lazy[T any](f func() Parser[T]) Parser[T] {
	var p Parser[T]

	return func(state State) (T, State, error) {
		if p == nil {
			p = f()
		}

		return p(state)
	}
}

func Concat(ps Parser[[]string]) Parser[string] {
	return Map(ps, func(ss []string) (string, error) {
		return strings.Join(ss, ""), nil
	})
}

func zero[T any]() T {
	var t T
	return t
}

func Parse[T any](p Parser[T], l lexer.Lexer, input string) (T, error) {
	tokens, err := l.Lex(input)
	if err != nil {
		return zero[T](), err
	}

	initialState := State{tokens: tokens, pos: 0, source: input}
	res, _, err := p(initialState)
	return res, err
}
