function handleBranchVersionValidation(branchVersion) {
    let isValid = isBranchVersionValid(branchVersion),
        id = branchVersion.attr("id");
    if (!isValid) {
        branchVersion.addClass('non-valid-input');
        $('.invalid-feedback.' + id).show();
    } else {
        branchVersion.removeClass('non-valid-input');
        $('.invalid-feedback.' + id).hide();
    }
    return isValid;
}

function isBranchVersionValid(branchVersion) {
    if (branchVersion.val().length === 0) {
        return false;
    }

    return !(!branchVersion.val() || !checkBranchVersion(branchVersion.val()));
}

function checkBranchVersion(branchVersion) {
    return /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)?$/i.test(branchVersion)
}

function processBranchName(name) {
    if (!name.trim()) {
        return `${name}`.toUpperCase()
    }
    return `${name}-`.toUpperCase()
}

function trimMinorVersionComponent(version) {
    let components = version.split('.');

    return `${components[0]}.${components[1]}`
}