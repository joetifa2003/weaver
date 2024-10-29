package parser

import "github.com/joetifa2003/weaver/internal/pargo"

func digit() pargo.Parser[string] {
	return pargo.OneOf(
		pargo.Exactly("0"),
		pargo.Exactly("1"),
		pargo.Exactly("2"),
		pargo.Exactly("3"),
		pargo.Exactly("4"),
		pargo.Exactly("5"),
		pargo.Exactly("6"),
		pargo.Exactly("7"),
		pargo.Exactly("8"),
		pargo.Exactly("9"),
	)
}

func alpha() pargo.Parser[string] {
	return pargo.Concat(
		pargo.Some(
			pargo.OneOf(
				pargo.Exactly("a"),
				pargo.Exactly("b"),
				pargo.Exactly("c"),
				pargo.Exactly("d"),
				pargo.Exactly("e"),
				pargo.Exactly("f"),
				pargo.Exactly("g"),
				pargo.Exactly("h"),
				pargo.Exactly("i"),
				pargo.Exactly("j"),
				pargo.Exactly("k"),
				pargo.Exactly("l"),
				pargo.Exactly("m"),
				pargo.Exactly("n"),
				pargo.Exactly("o"),
				pargo.Exactly("p"),
				pargo.Exactly("q"),
				pargo.Exactly("r"),
				pargo.Exactly("s"),
				pargo.Exactly("t"),
				pargo.Exactly("u"),
				pargo.Exactly("v"),
				pargo.Exactly("w"),
				pargo.Exactly("x"),
				pargo.Exactly("y"),
				pargo.Exactly("z"),
			),
		),
	)
}

func identifier() pargo.Parser[string] {
	return pargo.Sequence2(
		alpha(),
		pargo.Concat(
			pargo.Many(
				pargo.OneOf(
					alpha(),
					digit(),
					pargo.Exactly("_"),
				),
			),
		),
		func(first string, rest string) string {
			return first + rest
		},
	)
}
