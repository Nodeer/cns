function readyData() {
	var submitData = {
		id: id,
		data: JSON.stringify(inputTools.getJSON())
	};
	return submitData;
}

function send(url) {
	$.ajax({
		url: url,
		data: readyData(),
		method: 'POST',
		success: function(data) {
			console.log(data);
			window.location.href = data;
		},
		error: function(data) {
			console.log(data);
		}
	});
}

$(document).ready(function() {

	$('button#save').click(function() {
		send(url + 'save');
	});

	$('button#complete').click(function() {
		$('div#invalidMsg').addClass('hide');
		if (inputTools.validate()) {
			send(url + 'complete');
		} else {
			$('div#invalidMsg').removeClass('hide');
			$('html, body').animate({ scrollTop: 0 }, 'fast');
		}
	});

	if (data !== '') {
		inputTools.fill(JSON.parse(data));
	}
});
