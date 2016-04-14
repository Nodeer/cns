$(document).on('click', 'tr.clickable', function() {
    window.location.href = this.getAttribute('data-url');
});
