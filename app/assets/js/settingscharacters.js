$(document).ready(function(e) {
	$('a.character-set-default').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=characterSetDefault&characterID="+$(this).attr('characterID'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings/characters"
		});
	});
});