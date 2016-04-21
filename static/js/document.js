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
		success: function(resp) {
			if (resp.status === 'success') {
				alert(resp.msg);
			}
		},
		error: function(data) {
			alert("error")
			console.log(data);
		}
	});
}

$(document).ready(function() {

	$('button#save').click(function() {
		send(url + '/save');
	});

	$('button#complete').click(function() {
		$('div#invalidMsg').addClass('hide');
		if (inputTools.validate()) {
			send(url + '/complete');
		} else {
			$('div#invalidMsg').removeClass('hide');
			$('html, body').animate({ scrollTop: 0 }, 'fast');
		}
	});

	if (data !== '') {
		inputTools.fill(JSON.parse(data));
	}
});
