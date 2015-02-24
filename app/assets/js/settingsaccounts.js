$(document).ready(function(e) {
	$('a.api-key-add').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: $('#apiKeyAdd').serialize(),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings/accounts"
		});
	});

	$('a.api-key-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=apiKeyDelete&apiKeyID="+$(this).attr('apiKeyID'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings/accounts"
		});
	});
});