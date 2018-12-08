<?php

namespace codemirror;

function is_mobile() {
    $subject = $_SERVER['HTTP_USER_AGENT'];
    $pattern = '/(Android|iPhone|iPad)/i';
    return preg_match($pattern, $subject);
}


function post_head() {
    if(!is_mobile()) { ?>
        <link rel="stylesheet" href="/plugins/codemirror/codemirror.css" />
        <script src="/plugins/codemirror/codemirror.js"></script>

        <script src="/plugins/codemirror/xml.js"></script>
        <script src="/plugins/codemirror/css.js"></script>
        <script src="/plugins/codemirror/javascript.js"></script>
        <script src="/plugins/codemirror/htmlmixed.js"></script>
        <script src="/plugins/codemirror/markdown.js"></script>

        <script src="/plugins/codemirror/vim.js"></script>
        <link rel="stylesheet" href="/plugins/codemirror/dialog.css" />
        <script src="/plugins/codemirror/dialog.js"></script>
        <link rel="stylesheet" href="/plugins/codemirror/fullscreen.css" />
        <script src="/plugins/codemirror/fullscreen.js"></script>
        <style>
            .CodeMirror {
                border: 1px solid gray;
                font-family: Microsoft YaHei Mono, monospace;
                font-size: 13px;
                height: 70vh;
            }
        </style>
    <?php }
}

function post_foot() {
    if(!is_mobile()) { ?>
        <script>
            var textarea = document.getElementById('source');
            var codemirror = CodeMirror.fromTextArea(textarea, {
                lineNumbers: true,
                keyMap: 'vim',
                lineWrapping: true,
                extraKeys: {
                    'Alt-F': function(cm) {
                        cm.setOption('fullScreen', !cm.getOption('fullScreen'));
                    },
                },
            });
            codemirror.on('blur', function(cm) {
                textarea.value = cm.getValue();
            });
        </script>
    <?php }
}
