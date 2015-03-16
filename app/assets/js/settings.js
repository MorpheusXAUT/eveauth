$(document).ready(function(e) {
	$('a.settings-edit').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: $('#settingsEdit').serialize(),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings"
		});
	});
});