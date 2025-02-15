v0_local = @frame(vars:2, freeVars:0,) 	{
		v1_local = {name: "joe"}
		return {getName: @frame(vars:0, freeVars:1,) 			{
				return f0_free["name"]
			}}
	}
@echo(v0_local("joe")["getName"]())
