package pargo

import (
	"errors"
	"strconv"

	"github.com/joetifa2003/weaver/internal/pargo/lexer"
)

type State struct {
	source string
	tokens []lexer.Token
	pos    int
}

func (s *State) done() bool {
	return s.pos >= len(s.tokens)
}

func (s *State) consume() (lexer.Token, error) {
	if s.pos >= len(s.tokens) {
		loc := lexer.Location{}
		if len(s.tokens) > 0 {
			loc = s.tokens[len(s.tokens)-1].Location()
		}
		return lexer.RegexToken{
			Loc:   loc,
			Ttype: -1,
			Lit:   "EOF",
		}, nil
	}

	val := s.tokens[s.pos]
	s.pos++

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

func TokenType(ttype int) Parser[string] {
	return func(state State) (string, State, error) {
		old := state

		tok, err := state.consume()
		if err != nil {
			return "", old, err
		}

		if tok.Type() != ttype {
			return "", old, NewParseError(state.source, strconv.Itoa(ttype), tok)
		}

		return tok.String(), state, nil
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
			return "", old, NewParseError(state.source, "not "+s, tok)
		}

		return tok.String(), state, nil
	}
}

func OneOf[T any](parsers ...Parser[T]) Parser[T] {
	return func(initState State) (T, State, error) {
		var err error
		var res T
		var state State
		var loc lexer.Location
		var ferr error

		for _, p := range parsers {
			res, state, err = p(initState)
			if err == nil {
				return res, state, nil
			}

			var parseError ParseError
			if errors.As(err, &parseError) {
				if parseError.Location().IsAfter(loc) {
					loc = parseError.Location()
					ferr = err
				}
			}
		}

		return zero[T](), initState, ferr
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

func SomeSep[T any, S any](p Parser[T], separator Parser[S]) Parser[[]T] {
	return Sequence2(
		p,
		Many(
			Sequence2(
				separator,
				p,
				func(_ S, a T) T {
					return a
				},
			),
		),
		func(b T, a []T) []T {
			return append([]T{b}, a...)
		},
	)
}

func Optional[T any](p Parser[T]) Parser[*T] {
	return func(state State) (*T, State, error) {
		oldState := state

		res, state, err := p(state)
		if err != nil {
			return nil, oldState, nil
		}

		return &res, state, nil
	}
}

func ManySep[T any, S any](p Parser[T], separator Parser[S]) Parser[[]T] {
	return Sequence3(
		Optional(p),
		Many(
			Sequence2(
				separator,
				p,
				func(_ S, a T) T {
					return a
				},
			),
		),
		Optional(separator),
		func(b *T, a []T, _ *S) []T {
			if b == nil {
				return a
			}
			return append([]T{*b}, a...)
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

func ManyUntil[T any, E any](p Parser[T], end Parser[E]) Parser[[]T] {
	return func(state State) ([]T, State, error) {
		var res []T
		var err error
		var r T

		for {
			_, _, err = end(state)
			if err == nil {
				return res, state, nil
			}

			r, state, err = p(state)
			if err != nil {
				return res, state, err
			}
			res = append(res, r)
		}
	}
}

func ManyAll[T any](p Parser[T]) Parser[[]T] {
	return func(state State) ([]T, State, error) {
		var res []T
		var err error
		var r T

		for {
			r, state, err = p(state)
			if err != nil {
				return res, state, err
			}
			res = append(res, r)
			if state.done() {
				return res, state, nil
			}
		}
	}
}

func Some[T any](p Parser[T]) Parser[[]T] {
	return func(state State) ([]T, State, error) {
		first, state, err := p(state)
		if err != nil {
			return zero[[]T](), state, err
		}

		res := []T{first}

		var other T

		for {
			other, state, err = p(state)
			if err != nil {
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
