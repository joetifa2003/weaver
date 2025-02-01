v1 := [0, 1, 2]
v2 := [0, 1, 3]
{
	__$e := [v1, v2]
	__$v0 := 0
	__$v1 := 0
	x1 := 0
	y1 := 0
	z1 := 0
	x2 := 0
	y2 := 0
	z2 := 0
	if (((type(__$e) == "array") and (len(__$e) >= 2) and (__$v0 = __$e[0] or true) and ((type(__$v0) == "array") and (len(__$v0) >= 3) and (__$v1 = __$v0[0] or true) and (x1 = __$v1 or true) and (__$v1 = __$v0[1] or true) and (y1 = __$v1 or true) and (__$v1 = __$v0[2] or true) and (z1 = __$v1 or true)) and (__$v0 = __$e[1] or true) and ((type(__$v0) == "array") and (len(__$v0) >= 3) and (__$v1 = __$v0[0] or true) and (x2 = __$v1 or true) and (__$v1 = __$v0[1] or true) and (y2 = __$v1 or true) and (__$v1 = __$v0[2] or true) and (z2 = __$v1 or true))) and ((x1 == x2) and (y1 == y2) and (z1 == z2))) {
		{
			echo("equal")
		}
	} else {
		{
			echo("not equal")
		}
	}
}
