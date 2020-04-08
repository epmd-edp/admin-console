$(function () {
    $('.tooltip-icon').tooltip();

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            let stage = getUrlParameter('stage');
            if (anchor === '#stageSuccessModal') {
                showNotification(true, `Stage ${stage} was marked for deletion.`);
            } else if (anchor === '#stageIsUsedAsSource') {
                let $modal = $("#delete-confirmation");
                $('.confirmation-msg').text(`Confirm Deletion of '${stage}'`);
                $modal.find('.server-error').show();
                $modal.modal('show');
            }
            location.hash = '';
        }
    });

    !function () {
        $.each($('.applications-to-promote input'), function () {
            let appToPromote = $(this).data('app-name');
            $.each($('.applications-info .edp-table tbody tr'), function () {
                let $promoteEl = $(this).find('.promoteCDPipeline');
                if ($(this).find('.codebaseName').text().trim() === appToPromote) {
                    $promoteEl.find('.promote-checkbox-overview').removeClass('cancel').addClass('promoted');
                }
            });

        });
    }();

    !function () {
        $.each($('.platform-link a'), function () {
            if (!$(this).attr('href')) {
                $(this).addClass('hover-popup')
                    .attr('disabled', true)
                    .css('color', '#aaa');
            }
        })
    }();

    let moveLeft = 20,
        moveDown = 10,
        $link = $('a.hover-popup');

    $link.hover(function () {
        $('div#kubernetes-component').show();
    }, function () {
        $('div#kubernetes-component').hide();
    });

    $link.mousemove(function (e) {
        $("div#kubernetes-component")
            .css('top', e.pageY + moveDown).css('left', e.pageX + moveLeft);
    });

    $('.platform-link a.edp-link').click(function (e) {
        if ($(this).attr('disabled')) {
            e.preventDefault();
        }
    });

    $('.delete-stage').click(function () {
        let stage = $(this).data('name'),
            order = $(this).data('order'),
            $modal = $("#delete-confirmation");
        if (order === 0) {
            $modal.find('.popup-message.default').hide();
            $modal.find('.popup-message.last-stage').show();
        }
        $('.confirmation-msg').text(`Confirm Deletion of '${stage}'`);
        $modal.data('name', stage).modal('show');
        $('input#order').val(order);
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