function displayError(error) {
	$('div.col-md').prepend('<div class="alert alert-danger alert-dismissible fade in" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><strong>Ooops!</strong> '+error+'</div>');
	$('html, body').animate({ scrollTop: '0px' });
}

function displayAjaxError(jqXHR, textStatus, errorThrown) {
	switch (textStatus) {
		case null:
			displayError('Received unknown error while performing AJAX request');
			break;
		case "timeout":
			displayError('Received timeout while performing AJAX request');
			break;
		case "error":
			displayError('Received error while performing AJAX request: ' + errorThrown);
			break;
		case "abort":
			displayError('AJAX request was aborted');
			break;
		case "parsererror":
			displayError('Failed to parse AJAX request');
			break;
		default:
			displayError('Received unknown error while performing AJAX request');
	}
}

jQuery.fn.filterByText = function(textbox, selectSingleMatch) {
	return this.each(function() {
		var select = this;
		var options = [];
		$(select).find('option').each(function() {
			options.push({value: $(this).val(), text: $(this).text()});
		});
		$(select).data('options', options);
		$(textbox).bind('change keyup', function() {
			var options = $(select).empty().scrollTop(0).data('options');
			var search = $.trim($(this).val());
			var regex = new RegExp(search, "gi");
			$.each(options, function(i) {
				var option = options[i];
				if (option.text.match(regex) !== null) {
					$(select).append(
						$('<option>').text(option.text).val(option.value)
					);
				}
			});
			if (selectSingleMatch === true && $(select).children().length() === 1) {
				$(select).children().get(0).selected = true;
			}
		});
	});
};