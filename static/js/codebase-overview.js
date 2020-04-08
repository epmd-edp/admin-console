$(function () {

    let STATUS = {
        IN_PROGRESS: 'inactive'
    };
    let delayTime = 10000;

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#codebaseSuccessModal') {
                showNotification(true);
            } else if (anchor === '#codebaseIsUsed') {
                let codebase = getUrlParameter('codebase'),
                    $modal = $("#delete-confirmation");
                $('.confirmation-msg').text(`Confirm Deletion of '${codebase}'`);
                $modal.find('.server-error').show();
                $modal.modal('show');
            } else if (anchor === '#codebaseIsDeleted') {
                let codebase = getUrlParameter('codebase');
                showNotification(true, `Codebase ${codebase} was marked for deletion.`);
            }
            location.hash = '';
        }

        let uri = window.location.pathname;
        let codebaseName = getUrlParameter('waitingforcodebase');
        if (codebaseName) {
            let status = $("tr[data-codebase-name='" + codebaseName + "']").attr("data-codebase-status");
            if (status === STATUS.IN_PROGRESS) {
                uri += "?waitingforcodebase=" + codebaseName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });

    $('.delete-codebase').click(function () {
        let codebase = $(this).data('codebase'),
            $modal = $("#delete-confirmation");
        $('.confirmation-msg').text(`Confirm Deletion of '${codebase}'`);
        $modal.data('name', codebase).modal('show');
    });

    $('.delete-confirmation').click(function () {
        deleteConfirmation();
    });

    $('.close,.cancel-delete').click(function () {
        closeConfirmation();
    });

});

function showNotification(ok, msg, delay) {
    $.notify({
            icon: ok ? 'glyphicon glyphicon-ok-circle alert-icon' : 'glyphicon gglyphicon-warning-sign alert-icon',
            message: msg ? msg : (ok ? 'Provisioning has been started.' : 'Provisioning has been failed.')
        },
        {
            type: ok ? 'success' : 'error',
            delay: delay ? delay : 5000,
            animate: {
                enter: 'animated fadeInRight',
                exit: 'animated fadeOutRight'
            },
            onShow: function() {
                this.css({'width':'auto', 'display': 'flex'});
            },
        });
}
