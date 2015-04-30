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

	$('a.settings-application-edit-toggle').click(function() {
		$('#settingsApplicationsEditApplicationName').val($(this).attr('applicationName'));
		$('#settingsApplicationsEditApplicationCallback').val($(this).attr('applicationCallback'));
		$('#settingsApplicationsEditApplicationID').val($(this).attr('applicationID'));
		$('#settingsApplicationsEditApplication').collapse("show");
	});

	$('a.settings-application-edit-cancel').click(function() {
		$('#settingsApplicationsEditApplication').collapse("hide");
	});

	$('a.settings-application-edit-submit').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: $('#settingsApplicationsEditApplicationForm').serialize(),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings/applications"
		});
	});

	$('a.settings-application-edit-secret').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=settingsApplicationsEditApplicationResetSecret&applicationID="+$('#settingsApplicationsEditApplicationID').val()+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/settings/applications"
		});
	});
});