$(function () {

    let STATUS = {
        IN_PROGRESS: 'inactive'
    };
    let delayTime = 10000;

    $(document).ready(function () {

        let uri = window.location.pathname;
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#branchExistsModal') {
                let errorMessage = 'The release branch with the ' + getUrlParameter('errorExistingBranch') + ' name already exists. To proceed, use another branch name.';
                $('.branch-exists-modal').text(errorMessage).show();
                $('#releaseBranchModal').modal('show');
            } else if (anchor === '#branchSuccessModal') {
                showNotification(true);
            } else if (anchor === "#branchDeletedSuccessModal") {
                let name = getUrlParameter('name');
                showNotification(true, null, `Codebase Branch ${name} was marked for deletion.`);
            } else if (anchor === "#branchIsUsedSuccessModal") {
                let $modal = $("#delete-confirmation"),
                    name = getUrlParameter('name');
                $('.confirmation-msg').text(`Confirm Deletion of '${name}'`);
                $modal.find('.server-error').show();
                $modal.modal('show');
            }
        }

        let branchName = getUrlParameter('waitingforbranch');
        if (branchName) {
            let branchStatus = $("tr[data-branch-name='" + branchName + "']").attr("data-branch-status");
            if (branchStatus === STATUS.IN_PROGRESS) {
                uri += "?waitingforbranch=" + branchName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });

    $('.tooltip-icon').tooltip();

    $('#btn-modal-close, #btn-cross-close').click(function () {
        $('#releaseBranch').prop('checked', false);
        $('#commitNumber').val("");
        showBranchModalControls();
        $('.branch-exists-modal').hide();
        if ($('#versioningPostfix').length) {
            $('#branchName,#commitNumber,#branch-version,#master-branch-version').removeClass('non-valid-input');
            $('.invalid-feedback.master-branch-version').hide();
            restoreBranchModalWindowValues()
        } else {
            $('#branchName,#commitNumber,#branch-version').val('').removeClass('non-valid-input');
        }
        $('.invalid-feedback.branch-name').hide();
        $('.invalid-feedback.commit-message').hide();
        $('.invalid-feedback.branch-version').hide();
    });

    $('.modal-release-branch').click(function () {
        $('#releaseBranchModal').modal('show');
        if ($('#versioningPostfix').length) {
            let branchName = $('#branchName').val(),
                branchVersion = $('#branch-version').val(),
                masterBranchVersion = $('#master-branch-version').val();
            saveBranchModalWindowValues(branchName, branchVersion, masterBranchVersion)
        }
    });

    $('#create-release-branch').click(function () {
        $('.branch-exists-modal').hide();
        let isBranchValid = true;
        if ($('#releaseBranch').is(':not(:checked)')) {
            isBranchValid = handleBranchNameValidation();
        }
        let isCommitValid = handleCommitHashValidation();

        if ($("#branch-version").length === 0) {
            if (isBranchValid && isCommitValid) {
                $('#create-branch-action').submit();
            }
            return
        }

        if ($('#releaseBranch').is(':checked')) {
            let branchVersion = $('#branch-version'),
                masterBranchVersion = $('#master-branch-version'),
                isVersionValid = handleBranchVersionValidation(branchVersion),
                isMasterVersionValid = handleBranchVersionValidation(masterBranchVersion);
            if (isCommitValid && isVersionValid && isMasterVersionValid) {
                $('#create-branch-action').submit();
            }
        } else {
            let branchVersion = $('#branch-version'),
                isVersionValid = handleBranchVersionValidation(branchVersion);
            if (isBranchValid && isCommitValid && isVersionValid) {
                $('#create-branch-action').submit();
            }
        }
    });

    function showBranchModalControls() {
        let $createNewBranchModalEl = $('.create-new-branch-modal'),
            $versioningPostfixEl = $createNewBranchModalEl.find('.versioning-postfix'),
            $masterBranchVersionInputEl = $createNewBranchModalEl.find('.master-branch-version'),
            $branchNameInputEl = $createNewBranchModalEl.find('.branch-name'),
            $branchVersionInputEl = $createNewBranchModalEl.find('.branch-version');

        if ($('#releaseBranch').is(":checked")) {
            $('#branchName').removeClass('non-valid-input');
            $('.invalid-feedback.branch-name').hide();
            $branchNameInputEl.attr('readonly', true).val("release/" + trimMinorVersionComponent($branchVersionInputEl.val()));
            $versioningPostfixEl.val("RC");
            $masterBranchVersionInputEl.attr('disabled', false).removeClass('hide-element');
        } else {
            $branchNameInputEl.removeAttr('readonly');
            restoreBranchModalWindowValues();
            $versioningPostfixEl.val("SNAPSHOT");
            $masterBranchVersionInputEl.attr('disabled', true).addClass('hide-element');
        }
    }

    $('#branch-version').focusout(function () {
        let branchVersion = $('#branch-version');
        handleBranchVersionValidation(branchVersion);
    });

    $('#master-branch-version').focusout(function () {
        let masterBranchVersion = $('#master-branch-version');
        handleBranchVersionValidation(masterBranchVersion);
    });

    $('#commitNumber').focusout(function () {
        handleCommitHashValidation();
    });

    $('.delete-branch').click(function () {
        let name = $(this).data('name'),
            $modal = $("#delete-confirmation");
        $('.confirmation-msg').text(`Confirm Deletion of '${name}'`);
        $modal.data('name', name).modal('show');
    });

    $('.delete-confirmation').click(function () {
        deleteConfirmation();
    });

    $('.close,.cancel-delete').click(function () {
        closeConfirmation();
    });

    $('#branchName').on('input', function () {
        if ($('#releaseBranch').is(':not(:checked)')) {
            $('#versioningPostfix').val("SNAPSHOT");
            $('#versioningPostfix').val($('#versioningPostfix').val() + processBranchName($(this).val()));
        }
    });

    $('#branch-version').on('input', function () {
        if ($('#releaseBranch').is(":checked")) {
            $('#branchName').val("release/" + trimMinorVersionComponent($(this).val()));
        } else {
            $('#branchName').val($(this).val());
        }
    });

    $('#releaseBranch').change(function () {
        showBranchModalControls()
    });
});

function isBranchNameValid() {
    let $branchName = $('#branchName');
    return !(!$branchName.val() || !checkBranchName($branchName.val()));
}

function isHashCommitValid() {
    let $commitNumber = $('#commitNumber');

    if ($commitNumber.val().length === 0) {
        return true;
    } else {
        return !(!$commitNumber.val() || !checkHashCommit($commitNumber.val()));
    }
}

function checkBranchName(branchName) {
    return /^[a-z0-9][a-z0-9-.]*[a-z0-9]$/.test(branchName);
}

function handleBranchNameValidation() {
    let isBranchValid = isBranchNameValid();
    if (!isBranchValid) {
        $('#branchName').addClass('non-valid-input');
        $('.invalid-feedback.branch-name').show();
    } else {
        $('#branchName').removeClass('non-valid-input');
        $('.invalid-feedback.branch-name').hide();
    }
    return isBranchValid;
}

function checkHashCommit(hashCommit) {
    return /\b([a-f0-9]{40})\b/.test(hashCommit);
}

function handleCommitHashValidation() {
    let isCommitValid = isHashCommitValid();
    if (!isCommitValid) {
        $('#commitNumber').addClass('non-valid-input');
        $('.invalid-feedback.commit-message').show();
    } else {
        $('#commitNumber').removeClass('non-valid-input');
        $('.invalid-feedback.commit-message').hide();
    }
    return isCommitValid;
}

function showNotification(ok, delay, successMsg) {
    $.notify({
            icon: ok ? 'glyphicon glyphicon-ok-circle alert-icon' : 'glyphicon gglyphicon-warning-sign alert-icon',
            message: ok && successMsg != null ? successMsg : ok ? 'Provisioning has been started.' : 'Provisioning has been failed.'
        },
        {
            type: ok ? 'success' : 'error',
            delay: delay ? delay : 5000,
            animate: {
                enter: 'animated fadeInRight',
                exit: 'animated fadeOutRight'
            }
        });
}

function saveBranchModalWindowValues(branchName, branchVersion, masterBranchVersion) {
    sessionStorage.setItem("branch", JSON.stringify({
        "branchName": branchName,
        "branchVersion": branchVersion,
        "masterBranchVersion": masterBranchVersion
    }))
}

function restoreBranchModalWindowValues() {
    let branchConf = sessionStorage.getItem("branch");
    branchConf = JSON.parse(branchConf);
    $('#branchName').val(branchConf.branchName);
    $('#branch-version').val(branchConf.branchVersion);
    $('#master-branch-version').val(branchConf.masterBranchVersion)
}