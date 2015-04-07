$(document).ready(function(e) {
	$('a.admin-group-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminGroupDelete&groupID="+$(this).attr('groupID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/groups"
		});
	});
});