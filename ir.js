__$b0v0 := nil
__$b0v0 = [{name: "joe", age: 30}, {name: "jake", age: 20}, {name: "paule", age: 10}]
{
	__$b1v0 := nil
	__$b1v0 = 0
	loop {
		{
			if !((__$b1v0 < len(__$b0v0))) {
				break
			}
			{
				{
					__$b3v0 := nil
					__$b3v0 = __$b0v0[__$b1v0]
					{
						__$b4v2 := nil
						__$b4v1 := nil
						__$b4v0 := nil
						if (((type(__$b3v0) == "object") and (len(__$b3v0) >= 2) and (__$b4v0 = __$b3v0["name"] or true) and (__$b4v1 = __$b4v0 or true) and (__$b4v0 = __$b3v0["age"] or true) and (__$b4v2 = __$b4v0 or true)) and ((__$b4v2 >= 10) and (__$b4v2 <= 20))) {
							{
								echo((__$b4v1 + " is between 10 and 20"))
							}
						} else {
							if (__$b4v0 = __$b3v0 or true) {
								{

								}
							}
						}
					}
				}
			}
			__$b1v0 = (__$b1v0 + 1)
		}
	}
}
