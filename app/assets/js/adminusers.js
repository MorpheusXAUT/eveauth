$(document).ready(function(e) {
	$('a.admin-user-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminUserDelete&userID="+$(this).attr('userID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/users"
		});
	});
});