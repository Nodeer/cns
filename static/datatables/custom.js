$(document).on('click', 'tr.clickable', function() {
    if (this.getAttribute('data-target') === '_blank') {
        window.open(this.getAttribute('data-url'));
        return
    }
    window.location.href = this.getAttribute('data-url');
});
