$(function () {
    $('button.data-source').click(function (event) {
        validatePerfIntegrationBlock(event);
    });
});

function validatePerfIntegrationBlock(event) {
    let $dataSourceEl = $('.data-source-block');

    resetErrors($dataSourceEl);

    if (!$('.data-source-block').hasClass('hide-element') && !$dataSourceEl.find('.dataSources input').is(':checked')) {
        event.stopPropagation();
        blockIsNotValid($dataSourceEl);
        $('.data-source-checkbox-error').show();
        return false;
    }

    blockIsValid($dataSourceEl);
    $('.data-source-checkbox-error').hide();

    return true;
}

function resetErrors($el) {
    $el.find('input.is-invalid').removeClass('is-invalid');
    $el.find('.invalid-feedback').hide();
}