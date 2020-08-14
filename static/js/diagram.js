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
                    'shape': 'round-triangle',
                    'background-color': '#39c2d7',
                    'width': '40px',
                    'height': '40px'
                }
            },
            //codebase branch style
            {
                selector: 'node[type="branch"]',
                style: {
                    'shape': 'ellipse',
                    'background-color': '#7a797b',
                    'width': '20px',
                    'height': '20px'
                }
            },
            //codebase docker stream style
            {
                selector: 'node[type="codebase_docker_stream"]',
                style: {
                    'shape': 'barrel',
                    'background-color': '#39c2d7'
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
                    'shape': 'roundrectangle',
                    'background-color': '#a3c644',
                    'height': 30,
                    'width': 50,
                }
            },
            //stage style
            {
                selector: 'node[type="stage"]',
                style: {
                    'shape': 'circle',
                    'background-color': '#b22746',
                    'height': 20,
                    'width': 30,
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


