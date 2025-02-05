__$b0v0 := []
{
	__$b1v0 := 0
	loop 		{
			if !((__$b1v0 < 100000)) {
				break
			}
			{
				push(__$b0v0, {name: string(__$b1v0), age: 10})
				push(__$b0v0, {name: string(__$b1v0), age: 20})
				push(__$b0v0, {age: 30, name: string(__$b1v0)})
			}
			__$b1v0 = (__$b1v0 + 1)
		}
}
__$b0v1 := []
__$b0v2 := 0
{
	__$b1v0 := 0
	loop 		{
			if !((__$b1v0 < len(__$b0v0))) {
				break
			}
			{
				{
					__$b3v0 := __$b0v0[__$b1v0]
					{
						__$b4v2 := nil
						__$b4v1 := nil
						__$b4v0 := nil
						if (((type(__$b3v0) == "object") and (len(__$b3v0) >= 2) and (__$b4v0 = __$b3v0["age"] or true) and (__$b4v1 = __$b4v0 or true) and (__$b4v0 = __$b3v0["name"] or true) and (__$b4v2 = __$b4v0 or true)) and ((__$b4v1 >= 10) and (__$b4v1 <= 20))) {
							{
								push(__$b0v1, __$b4v1)
							}
						} else {
							if (__$b4v0 = __$b3v0 or true) {
								{
									__$b0v2 = (__$b0v2 + 1)
								}
							}
						}
					}
				}
			}
			__$b1v0 = (__$b1v0 + 1)
		}
}
echo(len(__$b0v1))
echo(__$b0v2)
