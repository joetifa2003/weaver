package ir

func irAnd(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpAnd,
		exprs,
	}
}

func irOr(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpOr,
		exprs,
	}
}

func irEq(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpEq,
		exprs,
	}
}

func irNeq(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpNeq,
		exprs,
	}
}

func irLt(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpLt,
		exprs,
	}
}

func irGt(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpGt,
		exprs,
	}
}

func irLte(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpLte,
		exprs,
	}
}

func irGte(exprs ...Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpGte,
		exprs,
	}
}

func irString(s string) StringExpr {
	return StringExpr{Value: s}
}

func irInt(i int) IntExpr {
	return IntExpr{Value: i}
}

func irFloat(f float64) FloatExpr {
	return FloatExpr{Value: f}
}

func irBuiltIn(name string) BuiltInExpr {
	return BuiltInExpr{name}
}

func irIfExpr(cond Expr, trueExpr Expr, falseExpr Expr) IfExpr {
	return IfExpr{
		Condition: cond,
		TrueExpr:  trueExpr,
		FalseExpr: falseExpr,
	}
}

func irReturnExpr(expr Expr) ReturnExpr {
	return ReturnExpr{Expr: expr}
}

func irCall(callee Expr, args ...Expr) CallExpr {
	return CallExpr{
		Expr: callee,
		Args: args,
	}
}

func irHasType(expr Expr, t string) BinaryExpr {
	return irEq(
		irCall(irBuiltIn("type"), expr),
		StringExpr{t},
	)
}

func irOrTrue(value Expr) BinaryExpr {
	return BinaryExpr{
		BinaryOpOr,
		[]Expr{
			value,
			BoolExpr{Value: true},
		},
	}
}

func irIndex(expr Expr, index Expr) IndexExpr {
	return IndexExpr{
		Expr:  expr,
		Index: index,
	}
}
