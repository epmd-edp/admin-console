function handleBranchVersionValidation(branchVersion) {
    let isValid = isBranchVersionValid(branchVersion),
        elClass = branchVersion.attr("id");
    if (!isValid) {
        branchVersion.addClass('non-valid-input');
        $('.invalid-feedback.' + elClass).show();
    } else {
        branchVersion.removeClass('non-valid-input');
        $('.invalid-feedback.' + elClass).hide();
    }
    return isValid;
}

function isBranchVersionValid(branchVersion) {
    if (branchVersion.val().length === 0) {
        return false;
    } else {
        return !(!branchVersion.val() || !checkBranchVersion(branchVersion.val()));
    }
}

function checkBranchVersion(branchVersion) {
    return /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|(?!.*RC|.*GA|.*SNAPSHOT)\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$/i.test(branchVersion)
}