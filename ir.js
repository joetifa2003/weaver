__$b0v0 := 10000000
__$b0v1 := 0
__$b0v2 := 0
{
	__$b1v0 := 0
	loop 		{
			if !((__$b1v0 < __$b0v0)) {
				break
			}
			{
				{
					__$b3v0 := (__$b1v0 % 2)
					{
						if ((type(__$b3v0) == "int") and (__$b3v0 == 0)) {
							{
								__$b0v1 = (__$b0v1 + 1)
							}
						} else {
							if ((type(__$b3v0) == "int") and (__$b3v0 == 1)) {
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
echo(__$b0v1)
echo(__$b0v2)
