$(document).ready(function(e) {
	$('a.api-key-add').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: $('#apiKeyAdd').serialize(),
			dataType: "json",
			error: displayAjaxError,
			success: function(response) {
				console.log(response)
				if (response.success === true && response.error === null) {
					location.reload(true);
				} else {
					displayError(response.error);
				}
			},
			timeout: 10000,
			type: "PUT",
			url: "/settings/accounts"
		});
	});
});