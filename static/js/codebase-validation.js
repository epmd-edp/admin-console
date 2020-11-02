$(function () {

    $(document).ready(function () {
        $('#versioningType').change(function () {
            checkVersioningType($(this).val())
        });
    });

    let REGEX = {
        CAPACITY: /\w/,
        SERVICE_PATH: /^\/.*$/,
        SERVICE_NAME: /^$|^[a-z][a-z0-9-]*[a-z0-9]$/,
        VCS_LOGIN: /\w/,
        VCS_PASSWORD: /\w/,
        DESCRIPTION: /^[a-zA-Z0-9]/,
        CODEBASE_NAME: /^[a-z][a-z0-9-]*[a-z0-9]$/,
        CODEBASE_DEFAULT_BRANCH: /^[a-z0-9][a-z0-9]*[\/-]?[a-z0-9]*[a-z0-9]$/,
        REPO_LOGIN: /\w/,
        REPO_PASSWORD: /\w/,
        REPO_URL: /(?:^git|^ssh|^https?|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/,
        RELATIVE_PATH: /^\/.*$/
    };

    let DEPLOYMENT_SCRIPT = {
        OPENSHIFT_TEMPLATE: "openshift-template",
        HELM_CHART: "helm-chart"
    };

    let INTEGRATION_STRATEGIES = {
        CREATE: 'create',
        CLONE: 'clone',
        IMPORT: 'import'
    };

    $('.tooltip-icon').add('[data-toggle="tooltip"]').tooltip();

    !function () {
        let $deployScriptEl = $('.deploymentScript');
        if ($('.advanced-settings-block').data('openshift')) {
            $deployScriptEl.val(DEPLOYMENT_SCRIPT.OPENSHIFT_TEMPLATE);
            return;
        }
        $deployScriptEl.val(DEPLOYMENT_SCRIPT.HELM_CHART);
    }();

    !function () {
        $('.form-group .js-form-subsection select:not(.jenkinsSlave)').attr('disabled', true);

        $('.multi-module').addClass('hide-element');
        $('#multiModule').attr("disabled", true);
    }();

    !function () {
        $.each($('.route, .dataBase .card-body').find('input, select'), function () {
            if ($(this).is('input')) {
                $(this).attr('readonly', true);
            } else if ($(this).is('select')) {
                $(this).attr('disabled', true);
            }
        });
    }();

    !function () {
        if ($('.db-block').length !== 0) {
            _sendGetRequest(true, `${$('input[id="basepath"]').val()}/api/v1/storage-class`,
                function (storageClasses) {
                    var $select = $('#dbPersistentStorage');

                    $.each(storageClasses, function () {
                        $select.append('<option value="' + this.toString() + '">' + this.toString() + '</option>');
                    });
                }, function (resp) {
                    console.log(resp);
                })
        }
    }();

    $('#jiraServerToggle').change(function () {
        let $jiraEl = $('.jiraServerBlock'),
            $commitMessagePatternBlockEl = $('.commitMessagePatternBlock'),
            $ticketNamePatternBlockEl = $('.ticketNamePatternBlock');
        if ($(this).is(':checked')) {
            $jiraEl.removeClass('hide-element')
                .find('select[name="jiraServer"]')
                .prop('disabled', false);

            $commitMessagePatternBlockEl.removeClass('hide-element')
                .find('input[id="commitMessagePattern"]')
                .prop('disabled', false);

            $ticketNamePatternBlockEl.removeClass('hide-element')
                .find('input[id="ticketNamePattern"]')
                .prop('disabled', false);
            return;
        }
        $jiraEl.addClass('hide-element')
            .find('select[name="jiraServer"]')
            .prop('disabled', true);

        $commitMessagePatternBlockEl.addClass('hide-element')
            .find('input[id="commitMessagePattern"]')
            .prop('disabled', true);

        $ticketNamePatternBlockEl.addClass('hide-element')
            .find('input[id="ticketNamePattern"]')
            .prop('disabled', true);
    });

    function activateCloneBlock() {
        $('.other-language').removeClass('button-disable');

        $('.main-block').data('import-strategy', false);
        $('.app-name').removeClass('hide-element');

        $('.gitServerEl').addClass('hide-element');
        $('.gitRelativePathEl').addClass('hide-element');

        $('.repo-url').add($('.private-repo')).removeClass('hide-element');

        if ($('#isRepoPrivate').is(':checked')) {
            $('.repoLogin').add($('.repoPassword')).removeClass('hide-element');
        }
    }

    function activateCreateBlock() {
        $('.other-language').addClass('button-disable');

        $('.main-block').data('import-strategy', false);
        $('.app-name').removeClass('hide-element');

        $('.gitServerEl').addClass('hide-element');
        $('.gitRelativePathEl').addClass('hide-element');

        $('.repo-url').add($('.private-repo')).addClass('hide-element');
        $('.repoLogin').add($('.repoPassword')).addClass('hide-element');
    }

    function activateImportBlock() {
        $('.other-language').removeClass('button-disable');

        $('.main-block').data('import-strategy', true);
        $('.gitServerEl').removeClass('hide-element');
        $('.gitRelativePathEl').removeClass('hide-element');

        $('.repo-url').add($('.private-repo')).addClass('hide-element');
        $('.repoLogin').add($('.repoPassword')).addClass('hide-element');
    }

    function toggleStrategy(strategy) {
        if (strategy === INTEGRATION_STRATEGIES.CLONE) {
            activateCloneBlock();
        } else if (strategy === INTEGRATION_STRATEGIES.CREATE) {
            activateCreateBlock();
        } else {
            activateImportBlock();
        }
    }

    !function () {
        let strategy = $('#strategy').val().toLowerCase();
        toggleStrategy(strategy);
        toggleCiToolView(strategy);
    }();

    function checkVersioningType(value) {
        let $startVersioningFromEl = $('.start-versioning-from');
        if (value === 'default') {
            $('.form-group.startVersioningFrom').addClass('hide-element');
            $('#startVersioningFrom').attr("disabled", true).removeAttr("value", "0.0.0");
            resetErrors($startVersioningFromEl);
        } else {
            $('.form-group.startVersioningFrom').removeClass('hide-element');
            $('#startVersioningFrom').attr("disabled", false).attr("value", "0.0.0");
            resetErrors($startVersioningFromEl);
        }
    }

    $('#versioningType').change(function () {
        checkVersioningType($(this).val())
    });

    $('#languageSelection').on('change', function (e) {
        $('.frameworkError').hide();
        if ($(this).find('input:checked').val() === "Java") {
            $('#framework-java8').prop('checked', true);
        }
        $('#framework-other').prop('disabled', !($(this).find('input:checked').val() === "other"));

        $('.js-form-subsection, .appLangError').hide();
        let langDivEl = $($(e.target).data('target'));
        langDivEl.show();

        $('.js-form-subsection input[name="framework"]').prop('checked', false);

        $('.multi-module').addClass('hide-element');
        $('#multiModule').attr("disabled", true);

        $('.main-block').data('code-language', $(e.target).data('target').replace('.formSubsection-', ''));

        $.each($('.build-tool .js-form-subsection, .jenkinsSlave .js-form-subsection'), function () {
            if ($(this).hasClass($(e.target).data('target').substring(1))) {
                $(this).show();
                $(this).find('select').attr('disabled', false);
                $(this).find('select.buildTool').val($(this).find('select.buildTool option:first').val());
            } else {
                $(this).find('select').attr('disabled', true);
            }
        });

        let codebaseVal = $('.card.main-block').data('codebase-type');
        if (codebaseVal === 'application' || codebaseVal === 'library') {
            $('.java-build-tools').val('Gradle');
        } else {
            $('.java-build-tools').val('Maven');
        }

        $('.test-report-framework').val('allure');

        setJenkinsSlave($('.buildTool:enabled'));
    });

    $('.formSubsection-java .java-frameworks').change(function () {
        setJenkinsSlave($('.buildTool:enabled'));
    });

    $('.formSubsection-dotnet .form__input-wrapper').change(function () {
        setJenkinsSlave($('.buildTool:enabled'));
    });

    $('#isRepoPrivate').change(function () {
        let $login = $('.repoLogin'),
            $pass = $('.repoPassword');
        if ($(this).is(':checked')) {
            $login.add($pass).removeClass('hide-element');
        } else {
            $login.add($pass).addClass('hide-element');
            $login.add($pass).find('.invalid-feedback').hide();
            $login.add($pass).find('input').removeClass('is-invalid');
        }
    });

    $('#strategy').change(function () {
        toggleStrategy(this.value.toLowerCase());
        toggleCiToolView(this.value.toLowerCase());

        $('div.jenkins-slave').removeClass('hide-element')
            .find('select[name="jenkinsSlave"]');

        $('div.ci-provision').removeClass('hide-element')
            .find('select[name="jobProvisioning"]');

        let $jiraEl = $('#jiraServerToggle'),
            $jiraDivEl = $('div.jiraServerToggle');
        if ($jiraEl.is(':checked')) {
            $jiraEl.click();
        }

        $jiraDivEl.removeClass('hide-element');

        $('div.jenkins-slave')
            .find('select[name="jenkinsSlave"]')
            .attr('disabled', false);

        $('div.ci-provision')
            .find('select[name="jobProvisioning"]')
            .attr('disabled', false);

        $jiraEl.attr('disabled', false);
    });

    function toggleCiToolView(strategy) {
        let $ciEl = $('div.ciTools');
        if (INTEGRATION_STRATEGIES.IMPORT === strategy.toLowerCase()) {
            $ciEl.removeClass('hide-element')
                .find('select[name="ciTool"]');
            return
        }
        $ciEl.addClass('hide-element')
            .find('select[name="ciTool"]')
            .val('Jenkins');

    }

    $('select.ciTool').change(function () {
        toggleJenkinsSlaveView($(this).val());
        toggleCiProvisionView($(this).val());
        toggleJiraIntegrationView($(this).val());
    });

    function toggleJenkinsSlaveView(ciTool) {
        let $jsEl = $('div.jenkins-slave');
        if (ciTool === "GitLab CI") {
            $jsEl.addClass('hide-element')
                .find('select[name="jenkinsSlave"]')
                .attr('disabled', true);
            return
        }
        $jsEl.removeClass('hide-element')
            .find('select[name="jenkinsSlave"]')
            .attr('disabled', false);
    }

    function toggleCiProvisionView(ciTool) {
        let $pEl = $('div.ci-provision');
        if (ciTool === "GitLab CI") {
            $pEl.addClass('hide-element')
                .find('select[name="jobProvisioning"]')
                .attr('disabled', true);
            return
        }
        $pEl.removeClass('hide-element')
            .find('select[name="jobProvisioning"]')
            .attr('disabled', false);
    }

    function toggleJiraIntegrationView(ciTool) {
        let $jiraEl = $('#jiraServerToggle'),
            $jiraDivEl = $('div.jiraServerToggle');
        if (ciTool === "GitLab CI") {
            if ($jiraEl.is(':checked')) {
                $jiraEl.click();
            }
            $jiraDivEl.addClass('hide-element');
            $jiraEl.attr('disabled', true);
            return
        }
        $jiraDivEl.removeClass('hide-element');
        $jiraEl.attr('disabled', false);
    }

    $('#btn-modal-continue').click(function () {
        $('form.edp-form').submit();
        $('#confirmationPopup').modal('hide');
        $(".window-table-body").remove();
    });

    $("#btn-cross-close, #btn-modal-close").click(function () {
        $(".window-table-body").remove();
    });

    $('#needRoute').change(function () {
        let $exposeServiceBlockEl = $('.route-block'),
            $inputsEl = $exposeServiceBlockEl.find('input');

        if ($(this).is(":checked")) {
            $inputsEl.attr('readonly', false);
        } else {
            $inputsEl.attr('readonly', true);
        }

        $inputsEl.removeClass('is-invalid').next('.invalid-feedback').hide();
    });

    $('#needDb').change(function () {
        let $dbBlockEl = $('.db-block'),
            $inputsEl = $dbBlockEl.find('input'),
            $selectsEl = $dbBlockEl.find('select');

        if ($(this).is(":checked")) {
            $inputsEl.attr('readonly', false);
            $selectsEl.attr('disabled', false);
        } else {
            $inputsEl.attr('readonly', true);
            $selectsEl.attr('disabled', true);
        }

        $('.capacity-error.invalid-feedback').hide();
        $inputsEl.removeClass('is-invalid');
    });

    $('.codebase-info-button').click(function (event) {
        validateCodebaseInfo(event);
    });

    $('.application-submit,.autotest-submit,.library-submit').click(function (event) {
        validateMainInfo(event);
        validateAdvancedInfo(event);
    });

    $('.advanced-settings-submit').click(function (event) {
        validateMainInfo(event);
        validateAdvancedInfo(event);
    });

    $('.vcs-submit,.create-library,.create-autotest').click(function (event) {
        if ($(this).hasClass('create-autotest') || $(this).hasClass('create-library')) {
            event.preventDefault();

            let canCreateAutotest = validateCodebaseInfo(event) & validateMainInfo(event) & validateVCSInfo(event) & validateAdvancedInfo(event);
            if (canCreateAutotest) {
                createConfirmTable($(this).hasClass('create-autotest') ? '#createAutotest' : '#createLibrary');
                $('#confirmationPopup').modal('show');
            }
        } else {
            validateVCSInfo(event);
        }
    });

    $('.route-submit').click(function (event) {
        validateRouteInfo(event);
    });

    $('.db-submit').click(function (event) {
        let canCreateApplication = validateCodebaseInfo(event) &
            validateMainInfo(event) & validateVCSInfo(event) &
            validateRouteInfo(event) & validateDbInfo(event) & validateAdvancedInfo(event);
        if (canCreateApplication) {
            createConfirmTable('#createAppForm');
            $('#confirmationPopup').modal('show');
        }
    });

    $('.java-build-tools,.js-build-tools,.dotnet-build-tools,.groovy-pipeline-build-tools,.other-build-tools').change(function () {
        if (this.value === 'Maven') {
            $('#multiModule').attr("disabled", false);
            $('.multi-module').removeClass('hide-element');
        } else {
            $('.multi-module').addClass('hide-element');
            $('#multiModule').attr("disabled", true);
        }

        setJenkinsSlave($(this));
    });

    $('#startVersioningFrom').focusout(function () {
        let branchVersion = $('#startVersioningFrom');
        handleBranchVersionValidation(branchVersion);
    });

    $('#gitRelativePath').focusout(function () {
        if (!isFieldValid($(this), REGEX.RELATIVE_PATH)) {
            return;
        }
        $('#appName').val($(this).val().match(/([^\/]*)\/*$/)[1]);
    });

    function setJenkinsSlave(el) {
        let $slave = getSlaveElement(el);
        if ($slave.length) {
            $slave.prop({selected: true});
            return;
        }
        $('.jenkinsSlave').val($('.jenkinsSlave option:first').val());
    }

    function getSlaveElement(el) {
        let $frameworkVersion = $('input[name="framework"]:checked').val();
        if (!!$frameworkVersion) {
            return $(`.jenkinsSlave option:contains(${el.find(':selected').data('build-tool') + "-" + $frameworkVersion})`);
        }
        return $(`.jenkinsSlave option:contains(${el.find(':selected').data('build-tool')})`);
    }

    function validateCodebaseInfo(event) {
        let $codebaseBlockEl = $('.codebase-block');

        resetErrors($codebaseBlockEl);

        let isValid = isCodebaseInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($codebaseBlockEl);
            return isValid;
        }
        blockIsValid($codebaseBlockEl);

        return isValid;
    }

    function validateMainInfo(event) {
        let $mainBlockEl = $('.main-block');

        resetErrors($mainBlockEl);

        let isValid = isMainInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($mainBlockEl);
            return isValid;
        }
        blockIsValid($mainBlockEl);

        return isValid;
    }

    function validateAdvancedInfo(event) {
        let $advancedBlockEl = $('.advanced-settings-block');

        resetErrors($advancedBlockEl);

        let isValid = isAdvancedInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($advancedBlockEl);
            return isValid;
        }
        blockIsValid($advancedBlockEl);

        return isValid;
    }

    function validateVCSInfo(event) {
        let $vcsBlockEl = $('.vcs-block');

        resetErrors($vcsBlockEl);

        let isValid = $vcsBlockEl.length === 0 ? true : isVCSValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($vcsBlockEl);
            return isValid;
        }
        blockIsValid($vcsBlockEl);

        return isValid;
    }

    function validateRouteInfo(event) {
        let $exposeServiceBlockEl = $('.route-block');

        resetErrors($exposeServiceBlockEl);

        let isValid = isExposingServiceInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($exposeServiceBlockEl);
            return isValid;
        }
        blockIsValid($exposeServiceBlockEl);

        return isValid;
    }

    function validateDbInfo(event) {
        let $dbBlockEl = $('.db-block');

        resetErrors($dbBlockEl);

        let isValid = isDatabaseValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($dbBlockEl);
            return isValid;
        }
        blockIsValid($dbBlockEl);

        return isValid;
    }

    function resetErrors($el) {
        $el.find('input.is-invalid').removeClass('is-invalid');
        $el.find('.invalid-feedback').hide();
    }

    function isCodebaseInfoValid() {
        let isValid = true;
        let $codebaseBlockEl = $('.codebase-block'),
            $strategyEl = $codebaseBlockEl.find('#strategy');

        if ($strategyEl.length === 0 || $strategyEl.val().toLowerCase() === 'clone') {
            let $repoUrl = $('#gitRepoUrl'),
                isRepoUrlValid = isFieldValid($repoUrl, REGEX.REPO_URL),
                $repoMsg = $repoUrl.next('.invalid-feedback');

            if (isRepoUrlValid) {
                let $gitRepoMsg = $('.git-repo-error');

                let creds = {
                    url: $repoUrl.val()
                };

                let isRepoPrivate = $('#isRepoPrivate').is(':checked'),
                    $repoLogin = $('#repoLogin'),
                    $repoPassword = $('#repoPassword');
                if (isRepoPrivate) {
                    let isLoginValid = isFieldValid($repoLogin, REGEX.REPO_LOGIN);
                    if (!isLoginValid) {
                        $repoLogin.next('.invalid-feedback').show();
                        $repoLogin.addClass('is-invalid');
                        isValid = false;
                    }
                    let isPasswordValid = isFieldValid($repoPassword, REGEX.REPO_PASSWORD);
                    if (!isPasswordValid) {
                        $repoPassword.next('.invalid-feedback').show();
                        $repoPassword.addClass('is-invalid');
                        isValid = false;
                    }

                    if (isLoginValid && isPasswordValid) {
                        creds.login = $repoLogin.val();
                        creds.password = $repoPassword.val();
                    }
                }

                _sendPostRequest.bind(this)(false, `${$('input[id="basepath"]').val()}/api/v1/repository/available`, creds, $('input[name="_xsrf"]').val(),
                    function (isAvailable) {
                        if (isRepoPrivate) {
                            if (isAvailable) {
                                isValid = true;
                            } else {
                                $('.git-creds').show();
                                $repoUrl.addClass('is-invalid');
                                $repoLogin.addClass('is-invalid');
                                $repoPassword.addClass('is-invalid');
                                isValid = false;
                            }
                        } else {
                            if (isAvailable) {
                                isValid = true;
                            } else {
                                $gitRepoMsg.show();
                                $repoUrl.addClass('is-invalid');
                                isValid = false;
                            }
                        }
                    }, function () {
                        console.log('an error has occurred while checking repository accessibility')
                    });

            } else {
                isValid = false;
                $repoMsg.show();
                $repoUrl.addClass('is-invalid');
            }
        } else if ($strategyEl.val().toLowerCase() === 'import') {
            let $gitRelativePath = $('#gitRelativePath'),
                isGitRelativePathValid = isFieldValid($gitRelativePath, REGEX.RELATIVE_PATH),
                $errMsg = $gitRelativePath.next('.invalid-feedback');
            if (!isGitRelativePathValid) {
                isValid = false;
                $errMsg.show();
                $gitRelativePath.addClass('is-invalid');
            }
        }

        return isValid;
    }

    function isMainInfoValid() {
        let $codebaseEl = $('.main-block'),
            $codebaseInputEl = $('.codebase-name'),
            isCodebaseNameValid = true,
            importStrategy = !!$codebaseEl.data('import-strategy');

        if (!importStrategy) {
            isCodebaseNameValid = isFieldValid($codebaseInputEl, REGEX.CODEBASE_NAME);
            if (!isCodebaseNameValid) {
                $('.codebase-name-validation.regex-error').show();
                $codebaseInputEl.addClass('is-invalid');
            }
        }

        let $defaultBranchInputEl = $('.default-branch-name'),
        isDefaultBranchNameValid = isFieldValid($defaultBranchInputEl, REGEX.CODEBASE_DEFAULT_BRANCH);
        if (!isDefaultBranchNameValid) {
            $('.default-branch-name-validation.regex-error').show();
            $defaultBranchInputEl.addClass('is-invalid');
        }

        let $descriptionInputEl = $('#description'),
            $descriptionErrEl = $('.description-validation.regex-error'),
            isDescriptionValid = $descriptionInputEl.length === 0 ? true : isFieldValid($descriptionInputEl, REGEX.DESCRIPTION);

        if (!isDescriptionValid) {
            $descriptionErrEl.show();
            $descriptionInputEl.addClass('is-invalid');
        }

        let isLanguageChosen = $codebaseEl.find('.language input').is(':checked'),
            isFrameworkChosen = true;
        if (isLanguageChosen) {
            let language = $('.main-block').data('code-language');
            if (language !== "other") {
                let $frameworksEl = $codebaseEl.find(`.form__input-wrapper .formSubsection-${language} input`);
                isFrameworkChosen = $frameworksEl.length === 0 ? true : $frameworksEl.is(":checked");
                if (!isFrameworkChosen) {
                    $('.frameworkError').show();
                }
            }
        } else {
            $('.appLangError').show();
        }

        return isCodebaseNameValid && isDefaultBranchNameValid && isDescriptionValid && isLanguageChosen && isFrameworkChosen;
    }

    function isAdvancedInfoValid() {
        let $advancedSettingsEl = $('.advanced-settings-block'),
            $versioningInputEl = $('.start-versioning-from'),
            isStartVersioningFromValid = true,
            jiraIntegration = $('#jiraServerToggle').is(':checked'),
            isCommitMessageRegexValid = jiraIntegration ? $('#commitMessagePattern').val().length !== 0 : true,
            isTicketNameRegexValid = jiraIntegration ? $('#ticketNamePattern').val().length !== 0 : true;

        if ($('#versioningType').val() === "edp") {
            isStartVersioningFromValid = isBranchVersionValid($versioningInputEl)
        }

        if (!isStartVersioningFromValid) {
            $('.invalid-feedback.startVersioningFrom').show();
            $advancedSettingsEl.addClass('is-invalid');
        }

        if (!isCommitMessageRegexValid) {
            $('.invalid-feedback.commitMessagePattern').show();
            $('#commitMessagePattern').addClass('is-invalid');
        }

        if (!isTicketNameRegexValid) {
            $('.invalid-feedback.ticketNamePattern').show();
            $('#ticketNamePattern').addClass('is-invalid');
        }

        return isStartVersioningFromValid && isCommitMessageRegexValid && isTicketNameRegexValid
    }

    function isVCSValid() {
        let $vcsLoginInputEl = $('#vcsLogin'),
            isVcsLoginValid = isFieldValid($vcsLoginInputEl, REGEX.VCS_LOGIN);

        if (!isVcsLoginValid) {
            $('.invalid-feedback.vcs-login-validation').show();
            $vcsLoginInputEl.addClass('is-invalid');
        }

        let $vcsPasswordInputEl = $('#vcsPassword'),
            isVcsPasswordValid = isFieldValid($vcsPasswordInputEl, REGEX.VCS_PASSWORD);

        if (!isVcsPasswordValid) {
            $('.invalid-feedback.vcs-password-validation').show();
            $vcsPasswordInputEl.addClass('is-invalid');
        }

        return isVcsLoginValid && isVcsPasswordValid;
    }

    function isExposingServiceInfoValid() {
        let needRoute = $('#needRoute').is(':checked');

        if (needRoute) {
            let $serviceNameInputEl = $('#routeSite'),
                isServiceNameValid = isCodebaseSiteFieldValid($serviceNameInputEl, REGEX.SERVICE_NAME);

            if (!isServiceNameValid) {
                $('.route-site.invalid-feedback').show();
                $serviceNameInputEl.addClass('is-invalid');
            }

            let $servicePathInputEl = $('#routePath'),
                isServicePathValid = isFieldValid($servicePathInputEl, REGEX.SERVICE_PATH);

            if (!isServicePathValid) {
                $('.route-path.invalid-feedback').show();
                $servicePathInputEl.addClass('is-invalid');
            }

            return isServiceNameValid && isServicePathValid;
        }

        return true;
    }

    function isDatabaseValid() {
        let needDb = $('#needDb').is(':checked');

        if (needDb) {
            let $capacityInputEl = $('#dbCapacity'),
                isCapacityValid = isFieldValid($capacityInputEl, REGEX.CAPACITY);

            if (!isCapacityValid) {
                $('.capacity-error.invalid-feedback').show();
                $capacityInputEl.addClass('is-invalid');
            }

            return isCapacityValid;
        }

        return true;
    }
})
;
