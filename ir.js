v0_local = 0
v1_local = 0
{
	v2_local = 0
	loop 		{
			if !((v2_local < 10000000)) {
				break
			}
			{
				if ((v2_local % 2) == 0) {
					{
						v0_local = (v0_local + 1)
					}
				} else {
					{
						v1_local = (v1_local + 1)
					}
				}
			}
			v2_local = (v2_local + 1)
		}
}
@echo(v0_local)
@echo(v1_local)
