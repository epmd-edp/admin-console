$('#perfServerToggle').change(function () {
    let $perfBlockEl = $('div.perfServerBlock'),
        $dataSourceBlockEl = $('div.data-source-block');
    if ($(this).is(':checked')) {
        $perfBlockEl.removeClass('hide-element')
            .find('select[name="perfServer"]')
            .prop('disabled', false);
        $dataSourceBlockEl.removeClass('hide-element');
        displayButtonWhenPerfIsSelected();
        return
    }
    $perfBlockEl.addClass('hide-element')
        .find('select[name="perfServer"]')
        .prop('disabled', true);
    $dataSourceBlockEl.addClass('hide-element');
    displayButtonWhenPerfIsNotSelected();
});

function displayButtonWhenPerfIsSelected() {
    let $advBlockEl = $('div.card.advanced-settings-block'),
        $codebaseType = $advBlockEl.data('codebase-type'),
        $advancedButtonEl = $('button.advanced-settings-submit');
    if ($codebaseType === 'application') {
        $advancedButtonEl.attr('data-target', '#collapseDataSource');
    } else if ($codebaseType === 'autotests') {
        $('button.adv-setting-autotest-submit')
            .attr('disabled', false)
            .removeClass('hide-element')
            .attr('data-target', '#collapseDataSource');
        $('button.adv-setting-create-autotest')
            .attr('disabled', true)
            .addClass('hide-element');
    } else {
        $('button.adv-setting-library-submit')
            .attr('disabled', false)
            .removeClass('hide-element')
            .attr('data-target', '#collapseDataSource');
        $('button.adv-setting-create-library')
            .attr('disabled', true)
            .addClass('hide-element');
    }
}

function displayButtonWhenPerfIsNotSelected() {
    let $advBlockEl = $('div.card.advanced-settings-block'),
        vcsEnabled = $advBlockEl.data('vcs-enabled') === true,
        $codebaseType = $advBlockEl.data('codebase-type'),
        $advancedButtonEl = $('button.advanced-settings-submit');
    if ($codebaseType === 'application') {
        $advancedButtonEl.attr('data-target', vcsEnabled ? '#collapseVCS' : '#collapseFour');
    } else if ($codebaseType === 'autotests') {
        if (vcsEnabled) {
            $('button.adv-setting-autotest-submit').attr('data-target', '#collapseVCS');
        } else {
            $('button.adv-setting-autotest-submit')
                .attr('disabled', true)
                .addClass('hide-element');
            $('button.adv-setting-create-autotest')
                .attr('disabled', false)
                .removeClass('hide-element');
        }
    } else {
        if (vcsEnabled) {
            $('button.adv-setting-library-submit').attr('data-target', '#collapseVCS');
        } else {
            $('button.adv-setting-library-submit')
                .attr('disabled', true)
                .addClass('hide-element');
            $('button.adv-setting-create-library')
                .attr('disabled', false)
                .removeClass('hide-element');
        }
    }
}

!function () {
    togglePerfDataSourcesView();
}();

$('#strategy').change(function () {
    togglePerfDataSourcesView();
});

function togglePerfDataSourcesView() {
    let strategyType = $('#strategy').val(),
        $dsDivEl = $('.dataSources').find('div input[value="GitLab"]').parent('div'),
        importStrategy = 'import';

    if (importStrategy === strategyType.toLowerCase()) {
        $dsDivEl.removeClass('hide-element');
        return
    }
    $dsDivEl.addClass('hide-element');
}