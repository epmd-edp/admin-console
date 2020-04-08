$(function () {

    $(document).ready(function () {
        let uri = window.location.pathname;
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#cdPipelineSuccessModal') {
                showNotification(true, null, 'Provisioning has been started.', 'Provisioning has been failed.');
            } else if (anchor === '#cdPipelineDeletedSuccessModal') {
                let name = getUrlParameter('name');
                showNotification(true, null, `CD Pipeline ${name} was marked for deletion.`);
            } else if (anchor === '#cdPipelineIsUsedAsSource') {
                let name = getUrlParameter('name'),
                    $modal = $("#delete-confirmation");
                $('.confirmation-msg').text(`Confirm Deletion of '${name}'`);
                $modal.find('.server-error').show();
                $modal.modal('show');
            } else {
                showNotification(true, null, 'The pipeline has been edited successfully.', 'Editing has been failed.');
            }
            location.hash = '';
        }

        let pipelineName = getUrlParameter('waitingforcdpipeline');
        if (pipelineName) {
            let status = getCDPipelineStatus(pipelineName);
            if (status === STATUS.IN_PROGRESS) {
                uri += "?waitingforcdpipeline=" + pipelineName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }

            window.history.replaceState({}, document.title, uri);
        }
    });

    $('.delete-cd-pipeline').click(function () {
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
});

let STATUS = {
    IN_PROGRESS: 'inactive'
};
let delayTime = 3000;

function getCDPipelineStatus(pipelineName) {
    let status;
    $.each($('.cd-pipeline-name'), function () {
        if ($(this).text().trim() === pipelineName) {
            status = $(this).parents('tr').find('.cd-pipeline-status').data('status').trim();
            return false;
        }
    });
    return status;
}

function showNotification(ok, delay, successMsg, failMsg) {
    $.notify({
            icon: ok ? 'glyphicon glyphicon-ok-circle alert-icon' : 'glyphicon gglyphicon-warning-sign alert-icon',
            message: ok ? successMsg : failMsg
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