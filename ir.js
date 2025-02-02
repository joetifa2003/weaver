__$b0v0 := [["go", "run", "main.go"]]
{
	__$b1v0 := __$b0v0[0]
	{
		__$b2v3 := nil
		__$b2v2 := nil
		__$b2v1 := nil
		__$b2v0 := nil
		if (((type(__$b1v0) == "array") and (len(__$b1v0) >= 3) and (__$b2v0 = __$b1v0[0] or true) and ((type(__$b2v0) == "string") and (__$b2v0 == "go")) and (__$b2v0 = __$b1v0[1] or true) and (__$b2v1 = __$b2v0 or true) and (__$b2v0 = __$b1v0[2] or true) and (__$b2v2 = __$b2v0 or true)) and (__$b2v1 == "run")) {
			{
				echo("running")
			}
		} else {
			if (((type(__$b1v0) == "array") and (len(__$b1v0) >= 4) and (__$b2v0 = __$b1v0[0] or true) and ((type(__$b2v0) == "string") and (__$b2v0 == "go")) and (__$b2v0 = __$b1v0[1] or true) and (__$b2v1 = __$b2v0 or true) and (__$b2v0 = __$b1v0[2] or true) and (__$b2v2 = __$b2v0 or true) and (__$b2v0 = __$b1v0[3] or true) and ((type(__$b2v0) == "object") and (len(__$b2v0) >= 1) and (__$b2v3 = __$b2v0["a"] or true) and ((type(__$b2v3) == "string") and (__$b2v3 == "something")))) and (__$b2v1 == "build")) {
				{
					echo("running")
				}
			}
		}
	}
}
