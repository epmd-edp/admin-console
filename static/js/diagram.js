function buildCodebaseData() {
    let elements = [];
    $.each(JSON.parse($('#diagram').attr('codebase-attr')), function (ci, cv) {
        elements.push(
            {
                data: {id: cv.id + '_codebase', type: 'codebase', name: cv.name}
            }
        );
        $.each(cv.codebase_branch, function (bi, bv) {
            elements.push(
                {
                    data: {id: bv.id + '_branch', type: 'branch', name: bv.branchName}
                }
            );

            elements.push(
                {
                    data: {
                        id: 'from' + cv.id + '_codebase' + 'to' + bv.id + '_branch',
                        source: cv.id + '_codebase',
                        target: bv.id + '_branch'
                    }
                },
            );

            $.each(bv.codebaseDockerStream, function (cdsi, cdsv) {
                elements.push(
                    {
                        data: {
                            id: 'from' + bv.id + '_branch' + 'to' + cdsv.ocImageStreamName + '_codebase_docker_stream',
                            source: bv.id + '_branch',
                            target: cdsv.ocImageStreamName,
                            type: 'codebase_docker_stream'
                        }
                    },
                );
            });

        });
    });
    return elements
}

function buildDataDockerStream() {
    let elements = [];
    $.each(JSON.parse($('#diagram').attr('codebase-docker-stream-attr')), function (cdsi, cdsv) {
        elements.push(
            {
                data: {id: cdsv, type: 'codebase_docker_stream', name: cdsv}
            }
        );
    });
    return elements;
}

function buildCdPipeline() {
    let elements = [];
    $.each(JSON.parse($('#diagram').attr('pipeline-attr')), function (pi, pv) {
        elements.push(
            {
                data: {id: pv.id + '_pipeline', type: 'pipeline', name: pv.name}
            }
        );

        $.each(pv.cd_stage, function (si, sv) {
            elements.push(
                {
                    data: {id: sv.id + '_stage', type: 'stage', name: sv.name}
                }
            );

            elements.push(
                {
                    data: {
                        id: 'from' + pv.id + '_pipeline' + 'to' + sv.id + '_stage',
                        source: pv.id + '_pipeline',
                        target: sv.id + '_stage'
                    }
                },
            );


            $.each(sv.qualityGates, function (qgi, qgv) {
                if (qgv.qualityGateType === 'autotests') {
                    elements.push(
                        {
                            data: {
                                id: 'from' + sv.id + '_stage' + 'to' + qgv.branchId + '_branch',
                                source: sv.id + '_stage',
                                target: qgv.branchId + '_branch'
                            }
                        },
                    );
                }
            });

            $.each(sv.stageCodebaseDockerStream, function (cdsi, cdsv) {
                elements.push(
                    {
                        data: {
                            id: 'from' + sv.id + '_stage' + 'to' + cdsv.inputCodebaseDockerStreamId,
                            source: sv.id + '_stage',
                            target: cdsv.inputCodebaseDockerStreamId
                        }
                    },
                    {
                        data: {
                            id: 'from' + sv.id + '_stage' + 'to' + cdsv.outputCodebaseDockerStreamId,
                            source: sv.id + '_stage',
                            target: cdsv.outputCodebaseDockerStreamId
                        }
                    }
                );
            });
        });
    });
    return elements;
}

function initDiagram() {
    let diagram = cytoscape({
        container: $('#diagram'),
        elements: buildCodebaseData().concat(buildDataDockerStream()).concat(buildCdPipeline()),
        style: [
            {
                selector: 'node',
                style: {
                    'label': 'data(name)',
                    'font-size': '7px',
                    'font-weight': 'bold',
                    'text-valign': 'bottom',
                }
            },
            //codebase style
            {
                selector: 'node[type="codebase"]',
                style: {
                    'background-fit': 'cover cover',
                    'background-image': '/static/img/codebase.png',
                    'background-color': '#ffffff',
                    'background-width': '0.1px',
                    'background-height': '0.1px',
                    'shape': 'square',
                    'width': '40px',
                    'height': '40px'
                }
            },
            //codebase branch style
            {
                selector: 'node[type="branch"]',
                style: {
                    'background-fit': 'cover cover',
                    'background-image': '/static/img/branch.png',
                    'background-color': '#ffffff',
                    'background-width': '0.1px',
                    'background-height': '0.1px',
                    'shape': 'square',
                    'width': '40px',
                    'height': '40px'
                }
            },
            //codebase docker stream style
            {
                selector: 'node[type="codebase_docker_stream"]',
                style: {
                    'background-fit': 'cover cover',
                    'background-image': '/static/img/registry.png',
                    'background-color': '#ffffff',
                }
            },
            //codebase docker stream line style
            {
                selector: 'edge[type="codebase_docker_stream"]',
                style: {
                    'line-style': 'dashed'
                }
            },
            //cd pipeline style
            {
                selector: 'node[type="pipeline"]',
                style: {
                    'background-fit': 'cover cover',
                    'background-image': '/static/img/cd-pipeline.png',
                    'background-color': '#ffffff',
                    'background-width': '0.1px',
                    'background-height': '0.1px',
                    'shape': 'square',
                    'height': '40px',
                    'width': '40px',
                }
            },
            //stage style
            {
                selector: 'node[type="stage"]',
                style: {
                    'background-fit': 'cover cover',
                    'background-image': '/static/img/stage.png',
                    'background-color': '#ffffff',
                    'background-width': '0.1px',
                    'background-height': '0.1px',
                    'shape': 'square',
                    'height': '40px',
                    'width': '40px',
                }
            },
            //edge style
            {
                selector: 'edge',
                style: {
                    'width': 1,
                    'line-color': '#ccc',
                    'target-arrow-color': '#ccc',
                    'target-arrow-shape': 'triangle',
                    'curve-style': 'bezier'
                }
            },
        ],
        layout: {
            name: 'dagre',
            rankDir: 'LR',
            fit: true,
            padding: 30,
            avoidOverlap: true,
            avoidOverlapPadding: 10,
            nodeDimensionsIncludeLabels: false,
            condense: true,
            rows: 5,
            cols: 5,
        }
    });

    //display nodes related to parent
    diagram.on('tap', 'node', function () {
        diagram.filter('node').style("display", "none");
        this.style("display", "element");
        this.successors().targets().style("display", "element");
    });

    //display whole diagram
    diagram.on('tap', function (event) {
        if (event.target === diagram) {
            diagram.filter('node').style("display", "element");
        }
    });
}

initDiagram();


