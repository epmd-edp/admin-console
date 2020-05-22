$(function () {

    function isStageAdded() {
        let currentStageName = $('#stageName').val();
        let stageNames = $('input[name="stageName"][data-mode="add"]').map(function () {
            return $(this).val();
        }).get();
        return !!($.inArray(currentStageName, stageNames) !== -1);
    }

    $('.add-stage').click(function () {
        let $stageAddedMsgEl = $('.invalid-feedback.stage-added-msg');
        if (validateFields()) {
            if (isStageAdded()) {
                $stageAddedMsgEl.show();
                return;
            }
            $stageAddedMsgEl.hide();
            appendStage(collectStageData());
            resetFields();

            $('#stage-creation').modal('hide');
        }
    });

    $('.confirm-edit-stage').click(function () {
        function replaceOldValues(stageToEdit, stageData) {
            stageToEdit.removeClass($('#stage-creation').attr('old-name'));
            stageToEdit.find('.stage-name a').text(stageData.stageName);
            stageToEdit.find('.edit-stage').attr('name', stageData.stageName);
            stageToEdit.find('.remove-stage').attr('name', stageData.stageName);
            stageToEdit.addClass(stageData.stageName);
            stageToEdit.find('#stageNameForm').val(stageData.stageName);
            stageToEdit.find('#stageDescForm').val(stageData.stageDesc).attr('name', stageData.stageName + '-stageDesc');
            stageToEdit.find('#triggerTypeForm').val(stageData.triggerType).attr('name', stageData.stageName + '-triggerType');
            stageToEdit.find('#jobProvisioningForm').val(stageData.jobProvisioning).attr('name', stageData.stageName + '-jobProvisioning');
            stageToEdit.find('#pipelineLibraryNameForm').val(stageData.pipelineLibraryName);
            stageToEdit.find('#pipelineLibraryBranchForm').val(stageData.pipelineLibraryBranch);

            let $stageBlockEl = $('.stage-info.' + stageData.stageName);
            $stageBlockEl.find('.qualityGateType, .stepName, .autotestsName, .branchName').remove();

            $.each(stageData.qualityGates, function () {
                let qualityGateTypeInputName = stageData.stageName + '-' + this.stepName + '-stageQualityGateType',
                    stepNameInputName = stageData.stageName + '-stageStepName',
                    autotestsInputName = stageData.stageName + '-' + this.stepName + '-stageAutotests',
                    branchInputName = stageData.stageName + '-' + this.stepName + '-stageBranch';

                $('<input class="qualityGateType" type="hidden" name="' + qualityGateTypeInputName + '" value="' + this.qualityGateType + '">').appendTo($stageBlockEl);
                $('<input class="stepName" type="hidden" name="' + stepNameInputName + '" value="' + this.stepName + '">').appendTo($stageBlockEl);
                $('<input class="autotestsName" type="hidden" name="' + autotestsInputName + '" value="' + this.autotestName + '">').appendTo($stageBlockEl);
                $('<input class="branchName" type="hidden" name="' + branchInputName + '" value="' + this.branchName + '">').appendTo($stageBlockEl)
            });
        }

        let $stageAddedMsgEl = $('.invalid-feedback.stage-added-msg');
        if (validateFields()) {
            if (isStageAdded()) {
                $stageAddedMsgEl.show();
                return;
            }
            $stageAddedMsgEl.hide();
            let stageData = collectStageData();
            let $stageCreationModal = $('#stage-creation');
            let stageToEdit = $('.stage-info.' + $stageCreationModal.attr('old-name'));
            replaceOldValues(stageToEdit, stageData);
            resetFields();
            toggleAdding();

            $stageCreationModal.modal('hide');
        }
    });

    $('.stage-modal-close, .cancel-edit-stage').click(function () {
        resetFields();
        toggleAdding();
    });

    $('.tooltip-icon').tooltip();

    $('#qualityGateType').change(function () {
        let $autotestsEl = $(this).parents('.quality-gate-row').find('.autotest-block-el');
        if (this.value === 'autotests') {
            $autotestsEl.removeClass('hide-element');
        } else {
            $autotestsEl.addClass('hide-element');
        }
    });

    $('.autotests-checkbox').change(function () {
        if (this.checked) {
            $('#qualityGateType').removeClass('non-valid-input');
            $('.autotests-validation-msg').hide();
        }
        disableSelectElems();
    });

    $('.add-quality-gate-row').click(function () {
        let $qualityGateTypeEl = $('.quality-gate-row:first').clone(true);
        $qualityGateTypeEl.find('.qualityGateTypeLabel, .nameOfStepLabel, .autotestLabel, .branchLabel').remove();
        $qualityGateTypeEl.find('#qualityGateType').val('manual');
        $qualityGateTypeEl.find('.remove-quality-gate-type').removeClass('hide-element');
        $qualityGateTypeEl.find('input.non-valid-input').removeClass('non-valid-input');
        $qualityGateTypeEl.find('.invalid-feedback.step-name-validation-msg').hide();
        $qualityGateTypeEl.find('.qualityGateType').val('manual');
        $qualityGateTypeEl.find('.autotest-block-el').addClass('hide-element');
        $qualityGateTypeEl.find('.nameOfStep').val('');
        $qualityGateTypeEl.insertBefore($('.step-name-validation-msg'));
    });

    $('.remove-quality-gate-type').click(function () {
        $(this).parents('.quality-gate-row').remove();
    });

    $('.autotest-projects').change(function () {
        let selectedAutotest = $(this).val();
        $.each($(this).parents('.quality-gate-row').find('.autotest-branches'), function () {
            $(this).data('selected-autotest') === selectedAutotest ? $(this).show() : $(this).hide();
        })
    });

    $('.pipeline-library').change(function () {
        let $branchEl = $(this).parents('.pipeline-library-row').find('.branch-block-el');
        if (this.value === 'default') {
            $branchEl.addClass('hide-element');
        } else {
            $branchEl.removeClass('hide-element');
        }

        let selectedPipelineLibrary = $(this).val();
        $.each($(this).parents('.pipeline-library-row').find('.pipeline-library-branches'), function () {
            $(this).data('selected-pipeline-library') === selectedPipelineLibrary ? $(this).show() : $(this).hide();
            $(this).prop('selectedIndex', 0);
        })
    });

    !function () {
        $('.quality-gate-row .autotest-branches').hide();

        $.each($('.quality-gate-row .autotest-projects'), function () {
            $('select[data-selected-autotest="' + $(this).val() + '"]').show();
        });

        $.each($('.pipeline-library-row .pipeline-library-branches'), function () {
            $('select[data-selected-pipeline-library="' + $(this).val() + '"]').show();
        });

    }();
});

function validateFields() {
    return handleStageNameValidation() & handleStageDescriptionValidation() & handleStepNameValidation();
}

function appendStage(stageData) {
    $('<div class="d-flex stage-info ' + stageData.stageName + '">\n' +
        '        <div class="form-group w-50 mb-2 stage-name">\n' +
        '<a class="edit-stage" onclick="editStage(this.name)" name="' + stageData.stageName + '" href="#">\n' + stageData.stageName + '</a>\n' +
        '        </div>\n' +
        '        <div class="d-flex flex-column justify-content-end mb-2">\n' +
        '            <button type="button" onclick="removeStage(this.name)" class="delete remove-stage" name="' + stageData.stageName + '" data-toggle="modal" data-target="#exampleModal">\n' +
        '                <i class="icon-trashcan"></i>\n' +
        '            </button>\n' +
        '        </div>\n' +
        '<input data-mode="add" id="stageNameForm" name="stageName" type="hidden" value="' + stageData.stageName + '">' +
        '<input id="stageDescForm" name="' + stageData.stageName + '-stageDesc" type="hidden" value="' + stageData.stageDesc + '">' +
        '<input id="triggerTypeForm" name="' + stageData.stageName + '-triggerType" type="hidden" value="' + stageData.triggerType + '">' +
        '<input id="jobProvisioningForm" name="' + stageData.stageName + '-jobProvisioning" type="hidden" value="' + stageData.jobProvisioning + '">' +
        '<input id="pipelineLibraryNameForm" name="' + stageData.stageName + '-pipelineLibraryName" type="hidden" value="' + stageData.pipelineLibraryName + '">' +
        '<input id="pipelineLibraryBranchForm" name="' + stageData.stageName + '-pipelineLibraryBranch" type="hidden" value="' + stageData.pipelineLibraryBranch + '">' +
        '    </div>').appendTo($('.stages-list'));

    let $stageBlockEl = $('.stage-info.' + stageData.stageName);
    $.each(stageData.qualityGates, function () {
        let qualityGateTypeInputName = stageData.stageName + '-' + this.stepName + '-stageQualityGateType',
            stepNameInputName = stageData.stageName + '-stageStepName',
            autotestsInputName = stageData.stageName + '-' + this.stepName + '-stageAutotests',
            branchInputName = stageData.stageName + '-' + this.stepName + '-stageBranch';

        $('<input class="qualityGateType" type="hidden" name="' + qualityGateTypeInputName + '" value="' + this.qualityGateType + '">').appendTo($stageBlockEl);
        $('<input class="stepName" type="hidden" name="' + stepNameInputName + '" value="' + this.stepName + '">').appendTo($stageBlockEl);
        $('<input class="autotestsName" type="hidden" name="' + autotestsInputName + '" value="' + this.autotestName + '">').appendTo($stageBlockEl);
        $('<input class="branchName" type="hidden" name="' + branchInputName + '" value="' + this.branchName + '">').appendTo($stageBlockEl)
    });
}

function resetFields() {
    $('#qualityGateType option:first').prop('selected', true);
    $('#triggerType option:first').prop('selected', true);
    $('#stage-creation input[type="text"]').val("");
    $('#pipeline-library option:first').prop('selected', true);

    $('input.non-valid-input, select.non-valid-input').removeClass('non-valid-input');
    $('div.invalid-feedback').hide();

    let $qualityGateTypeElems = $('.quality-gate-row');
    $qualityGateTypeElems.not(':first').remove();
    $qualityGateTypeElems.find('.qualityGateType').val('manual');
    $qualityGateTypeElems.find('.autotest-block-el').addClass('hide-element');

    $('.branch-block-el').addClass('hide-element');
}

function removeStage(stageName) {
    $('.stage-info.' + stageName).remove();
}

function editStage(stageName) {
    toggleEditing();
    let $stageCreationModal = $('#stage-creation');
    $stageCreationModal.attr('old-name', stageName);
    $stageCreationModal.modal('show');
    $('#stageNameForm[value="' + stageName + '"]').attr('data-mode', 'edit');
    fillFields(stageName);
    disableSelectElems();
}

function fillFields(stageName) {
    let $stageEl = $('.stage-info.' + stageName);
    $('#stageName').val($stageEl.find('#stageNameForm').val());
    $('#stageDesc').val($stageEl.find('#stageDescForm').val());
    $("#triggerType").val($stageEl.find('#triggerTypeForm').val());
    $("#jobProvisioning").val($stageEl.find('#jobProvisioningForm').val());
    $('#pipeline-library').val($stageEl.find('#pipelineLibraryNameForm').val()).change();
    $('#pipeline-library-branches').val($stageEl.find('#pipelineLibraryBranchForm').val()).change();

    let qualityGateData = collectOldQualityGatesData($stageEl);
    createQualityGateRows(qualityGateData);
}

function collectOldQualityGatesData($stageEl) {
    let $qualityGateTypeEl = $stageEl.find('.qualityGateType'),
        $stepNameEl = $stageEl.find('.stepName'),
        $autotestsNameEl = $stageEl.find('.autotestsName'),
        $branchNameEl = $stageEl.find('.branchName');

    let result = [];
    $.each($qualityGateTypeEl, function (i) {
        result.push({
            qualityGateType: $(this).val(),
            stepName: $($stepNameEl[i]).val(),
            autotestName: $(this).val() === 'autotests' ? $($autotestsNameEl[i]).val() : null,
            branchName: $(this).val() === 'autotests' ? $($branchNameEl[i]).val() : null
        });
    });

    return result;
}

function createQualityGateRows(qualityGateData) {
    for (let i = 0; i < qualityGateData.length - 1; i++) {
        let $firstQualityGateRowEl = $('.quality-gate-row:first').clone(true);
        $firstQualityGateRowEl.insertBefore($('.step-name-validation-msg'));
    }

    let $readyToFillRowElems = $('.quality-gate-row');
    $.each($readyToFillRowElems, function (i) {
        if (i !== 0) {
            $(this).find('.qualityGateTypeLabel, .nameOfStepLabel, .autotestLabel, .branchLabel').remove();
            $(this).find('.remove-quality-gate-type').show();
        }

        $(this).find('.qualityGateType').val(qualityGateData[i].qualityGateType);
        $(this).find('.nameOfStep').val(qualityGateData[i].stepName);

        if (qualityGateData[i].qualityGateType === 'autotests') {
            $(this).find('.autotest-block-el').removeClass('hide-element');
            $(this).find('.autotest-projects').val(qualityGateData[i].autotestName);
            $(this).find('.autotest-branches').hide();
            $(this).find('[data-selected-autotest="' + qualityGateData[i].autotestName + '"]').show();
            $(this).find('[data-selected-autotest="' + qualityGateData[i].autotestName + '"]').val(qualityGateData[i].branchName);
        }
    });
}

function toggleAdding() {
    $('.stage-info input[data-mode="edit"]').attr('data-mode', 'add');
    $('#add-header').show();
    $('#edit-header').hide();
    $('button.add-stage').show();
    $('button.confirm-edit-stage').hide();
}

function toggleEditing() {
    $('#add-header').hide();
    $('#edit-header').show();
    $('button.add-stage').hide();
    $('button.confirm-edit-stage').show();
}

function collectStageData() {
    let pipelineLibrary = $('#pipeline-library').val();
    return {
        stageName: $('#stageName').val(),
        stageDesc: $('#stageDesc').val(),
        pipelineLibraryName: $('#pipeline-library').val(),
        pipelineLibraryBranch: pipelineLibrary === 'default' ? null : $('.pipeline-library-row').find('[data-selected-pipeline-library="' + pipelineLibrary + '"]').val(),
        triggerType: $('#triggerType').val(),
        jobProvisioning: $('#jobProvisioning').val(),
        qualityGates: collectQualityGates(),
    };
}

function collectQualityGates() {
    let result = [];
    $.each($('.quality-gate-row'), function () {
        let qualityGateType = $(this).find('.qualityGateType').val();
        result.push({
            qualityGateType: qualityGateType,
            stepName: $(this).find('.nameOfStep').val(),
            autotestName: qualityGateType === 'autotests' ? $(this).find('.autotest-projects').val() : null,
            branchName: qualityGateType === 'autotests' ? $(this).find('[data-selected-autotest="' + $(this).find('.autotest-projects').val() + '"]').val() : null
        });
    });
    return result;
}

function handleStageNameValidation() {
    let $stageNameEl = $('#stageName');
    let valid = isStageNameValid();
    if (!valid) {
        $stageNameEl.addClass('non-valid-input');
        $stageNameEl.parents('div.form-group').find('.invalid-feedback.stage-name-msg').show();
    } else {
        $stageNameEl.removeClass('non-valid-input');
        $stageNameEl.parents('div.form-group').find('.invalid-feedback.stage-name-msg').hide();
    }
    return valid;
}

function handleStageDescriptionValidation() {
    let $stageDescEl = $('#stageDesc');
    let valid = isStageDescriptionValid();
    if (!valid) {
        $stageDescEl.addClass('non-valid-input');
        $stageDescEl.parents('div.form-group').find('.invalid-feedback').show();
    } else {
        $stageDescEl.removeClass('non-valid-input');
        $stageDescEl.parents('div.form-group').find('.invalid-feedback').hide();
    }
    return valid;
}

function handleStepNameValidation() {
    let $nameOfStepElems = $('.nameOfStep'),
        $validationMsgEl = $('.step-name-validation-msg'),
        $duplicateValidationMsgEl = $('.duplicate-step-name-validation-msg'),
        isValid = true;

    $.each($nameOfStepElems, function () {
        let isStepNameValid = isFieldValid($(this), /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/);
        if (!isStepNameValid) {
            isValid = false;
            $(this).addClass('non-valid-input');
            $validationMsgEl.show();
        } else {
            $(this).removeClass('non-valid-input');
        }
    });

    if (!isValid) {
        return false
    }

    let nameOfStepArr = getValues($nameOfStepElems),
        duplicateErr = doesArrayContainDuplicates(nameOfStepArr);
    if (duplicateErr) {
        $nameOfStepElems.addClass('non-valid-input');
        $duplicateValidationMsgEl.show();
    } else {
        $nameOfStepElems.removeClass('non-valid-input');
        $duplicateValidationMsgEl.hide();
    }

    if (isValid) {
        $validationMsgEl.hide();
    }

    return isValid && !duplicateErr;
}

function getValues($elemsArray) {
    let result = [];
    $.each($elemsArray, function () {
        result.push($(this).val());
    });
    return result
}

function doesArrayContainDuplicates(arr) {
    let sortedArr = arr.sort();

    let duplicates = [];
    for (let i = 0; i < sortedArr.length - 1; i++) {
        if (sortedArr[i + 1] === sortedArr[i]) {
            duplicates.push(sortedArr[i]);
        }
    }

    return duplicates.length > 0;
}

function isStageNameValid() {
    let checkStageName = function (stageName) {
        return /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/.test(stageName);
    };

    let $stageNameEl = $('#stageName');
    return !(!$stageNameEl.val() || !checkStageName($stageNameEl.val()));
}

function isStageDescriptionValid() {
    let $stageDescriptionEl = $('#stageDesc');
    return $stageDescriptionEl.val().length !== 0;
}