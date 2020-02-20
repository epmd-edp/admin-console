function deleteConfirmation() {
    let $modal = $("#delete-confirmation"),
        targetName = $modal.data('name'),
        confirmationName = $modal.find('#entity-name').val(),
        $errName = $modal.find('.invalid-feedback.different-name'),
        $errUsed = $modal.find('.server-error');
    if (targetName !== confirmationName) {
        $errName.show();
        return
    }
    $errName.hide();
    $errUsed.hide();
    $("#delete-action").submit();
}

function closeConfirmation() {
    let $modal = $("#delete-confirmation");
    $modal.find('.invalid-feedback.different-name').hide();
    $modal.find('.server-error').hide();
    $modal.find('#entity-name').val('');
}