$(function () {

    let REGEX = {
        PIPELINE_NAME: /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/
    };


    function validatePipelineInfo(event) {
        let $pipelineBlockEl = $('.pipeline-block');

        resetErrors($pipelineBlockEl);

        let isValid = isPipelineInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($pipelineBlockEl);
            return isValid;
        }
        blockIsValid($pipelineBlockEl);

        return isValid;
    }

    function isPipelineInfoValid() {
        let $pipelineNameInputEl = $('#pipelineName'),
            isPipelineNameValid = isFieldValid($pipelineNameInputEl, REGEX.PIPELINE_NAME);

        if (!isPipelineNameValid) {
            $('.invalid-feedback.pipeline-name-validation').show();
            $pipelineNameInputEl.addClass('is-invalid');
        }

        return isPipelineNameValid;
    }

    function validateApplicationInfo(event) {
        let $applicationBlockEl = $('.application-block');

        resetErrors($applicationBlockEl);

        let isValid = isApplicationInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($applicationBlockEl);
            return isValid;
        }
        blockIsValid($applicationBlockEl);

        return isValid;
    }

    function isApplicationInfoValid() {
        let isApplicationBlockValid = $('.app-checkbox').is(':checked');

        if (!isApplicationBlockValid) {
            $('.app-checkbox-error').show();
        }

        return isApplicationBlockValid;
    }

    function validateStageInfo() {
        let $stageBlockEl = $('.stage-block');

        resetErrors($stageBlockEl);

        let isValid = isStageInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($stageBlockEl);
            return isValid;
        }
        blockIsValid($stageBlockEl);

        return isValid;

    }

    function isStageInfoValid() {
        let isStageValid = $('.stages-list .stage-info').length > 0;

        if (!isStageValid) {
            $('.invalid-feedback.pipeline-name-validation').show();
            $('.stage-error').show();
        }

        return isStageValid;
    }

    function resetErrors($el) {
        $el.find('input.is-invalid').removeClass('is-invalid');
        $el.find('.invalid-feedback').hide();
    }

    $('.application-checkbox :checkbox').change(function () {
        let $selectEl = $('.select-' + $(this).attr('id')),
            $checkboxEl = $('.checkbox-' + $(this).attr('id'));
        if ($(this).is(':checked')) {
            $selectEl.prop('disabled', false);
            $checkboxEl.prop('disabled', false);
            $('.app-checkbox-error').hide();
            blockIsValid($('.application-block'));
        } else {
            $selectEl.prop('disabled', true);
            $checkboxEl.prop('disabled', true);
        }
    });

    $('#pipelineName').focusout(function () {
        let isPipelineValueValid = isPipelineValid();
        if (!isPipelineValueValid) {
            $('#pipelineName').addClass('non-valid-input');
            $('.invalid-feedback.pipeline-name').show();
        } else {
            $('#pipelineName').removeClass('non-valid-input');
            $('.invalid-feedback.pipeline-name').hide();
        }
    });

    function isPipelineValid() {
        let $pipelineNameEl = $('#pipelineName');
        return !(!$pipelineNameEl.val() || !checkPipelineName($pipelineNameEl.val()));
    }

    function checkPipelineName(pipelineName) {
        return /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/.test(pipelineName);
    }

    $('.pipeline-info-button').click(function (event) {
        validatePipelineInfo(event);
    });

    $('.application-info-button').click(function () {
        validateApplicationInfo(event);
    });

    $('.stage-info-button').click(function () {
        validateStageInfo(event);
    });

    $('.create-cd-pipeline').click(function (event) {
        event.preventDefault();

        let canCreateCDPipeline = validatePipelineInfo(event) &
            validateApplicationInfo(event) & validateStageInfo(event);

        if (canCreateCDPipeline) {
            $('#createCDCR').submit();
        }
    });

    $('.add-stage-modal').click(function () {
        $('#stage-creation').modal('show');
        disableSelectElems();
    });

});

function disableSelectElems() {
    $.each($('.autotests-checkbox'), function () {
        let $selectEl = $('select[name="' + $(this).attr('value') + '-autotestBranch"]');
        if (!$(this).is(':checked')) {
            $selectEl.attr('disabled', 'disabled');
        } else {
            $selectEl.removeAttr('disabled');
        }
    });
}