$(document).ready(function(e) {
	$('a.admin-groupdetails-role-toggle-granted').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminGroupDetailsRoleToggleGranted&groupID="+$(this).attr('groupID')+"&roleID="+$(this).attr('roleID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/groups"
		});
	});
	
	$('a.admin-groupdetails-role-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminGroupDetailsRoleDelete&groupID="+$(this).attr('groupID')+"&roleID="+$(this).attr('roleID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/groups"
		});
	});
});