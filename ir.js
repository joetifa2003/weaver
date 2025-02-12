{
	v0_local = ""
	{
		if ((@type(v0_local) == "object") && (@len(v0_local) >= 2) && (v1_local = v0_local["name"] || true) && (v2_local = v1_local || true) && (v3_local = v0_local["age"] || true) && (v4_local = v3_local || true)) {
			{
				@echo(v2_local)
			}
		} else {
			if (v1_local = v0_local || true) {
				{
					@echo(v1_local)
				}
			}
		}
	}
}
