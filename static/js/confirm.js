/*********************** Required Message Form **********************************\
 *																				*
 *		<div id="confirmAlert" class="alert clearfix hide">						*
 *			<form id="confirmForm" action="" method="post" class="col-lg-2">	*
 *				<button id="confirmButton" class="btn btn-sm">Yes</button>		*
 *				<a id="confirmCancel" class="btn btn-default btn-sm">No</a>		*
 *			</form>																*
 *			<span id="message"></span>											*
 *		</div>																	*
 *																				*
\********************************************************************************/

/*	confirm prompt example
	<a data-message="Are you sure you would like todo something?" data-color="danger" data-url="/endpoint" class="confirm-action">Deactivate</a>
*/

function Confirm() {
	this.registerDisplay();
}
Confirm.prototype = {
	color: '',
	registerDisplay: function() {
		$('.confirm-action').click(function(e) {
            e.stopPropagation();
            var btn = $(this);
            swal({
                title: '',
                text: btn.attr('data-message'),
                type: btn.attr('data-type'),
                showCancelButton: true,
                confirmButtonColor: btn.attr('data-color'),
                confirmButtonText: "Yes",
                closeOnConfirm: false
            }, function(){
                swal("Deleted!", "Your imaginary file has been deleted.", "success");
            });
		});
	}
};

var Confirm = new Confirm();
