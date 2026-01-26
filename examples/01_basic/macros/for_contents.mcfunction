scoreboard players set %[score] 0
function ./_for_each_%[score.name]
	%[...]
	scoreboard players add %[score] 1
	execute if score %[score] matches %[range] run function ./_for_each_%[score.name]
