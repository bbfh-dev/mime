scoreboard players set %[target] %[objective] 0
function ./_for_each_%[target.to_file_name]
	%[...]
	scoreboard players add %[target] %[objective] 1
	execute if score %[target] %[objective] matches %[range] run function ./_for_each_%[target.to_file_name]
