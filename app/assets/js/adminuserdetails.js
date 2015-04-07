$(document).ready(function(e) {
	$('a.admin-userdetails-group-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminUserDetailsGroupDelete&userID="+$(this).attr('userID')+"&groupID="+$(this).attr('groupID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/users"
		});
	});

	$('a.admin-userdetails-group-add').click(function() {
		
	});

	$('a.admin-userdetails-role-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminUserDetailsRoleDelete&userID="+$(this).attr('userID')+"&roleID="+$(this).attr('roleID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/users"
		});
	});

	$('a.admin-userdetails-role-add').click(function() {
		
	});

	$('a.admin-userdetails-account-delete').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=adminUserDetailsAccountDelete&userID="+$(this).attr('userID')+"&accountID="+$(this).attr('accountID')+"&csrfToken="+$(this).attr('csrfToken'),
			dataType: "json",
			error: displayAjaxError,
			success: displayResponse,
			timeout: 10000,
			type: "PUT",
			url: "/admin/users"
		});
	});
});