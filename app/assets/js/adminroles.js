$(document).ready(function(e) {
	$('a.admin-role-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminRolesDelete&roleID="+$(this).attr('roleID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/roles"
		});
	});
});