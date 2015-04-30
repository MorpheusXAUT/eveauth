$(document).ready(function(e) {
	$('a.settings-application-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=settingsApplicationsDelete&applicationID="+$(this).attr('applicationID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings/applications"
		});
	});
});