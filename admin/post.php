<?php

if($_SERVER['REQUEST_METHOD'] === 'GET') :

require_once('admin.php');

function post_widget_tag($p=null) {
    $tag = $p ? join(',', $p->tag_names) : '';

    $title = '标签';
    $classname = 'tags';
    $types = 'post';
    $content = <<<EOD
<input type="text" name="tags" value="$tag" placeholder="中英文逗号分隔" />
EOD;

    return compact('title', 'classname', 'types', 'content');
}

add_hook('post_widget', 'post_widget_tag');

function post_widget_files($p=null) {
    $title = '文件';
    $classname = 'files';
    $types = 'page,post';
    $content = <<<EOD
<label>文件列表：</label>
<ul class="list" style="max-height:200px;overflow:auto;">
</ul>
<label>文件上传：</label>
<span class="count"></span>
<div>
    <input type="file" multiple class="files" style="display:none;"/>
    <button class="refresh">刷新</button>
    <button class="browse">浏览</button>
    <button class="submit">上传</button>
    <progress class="progress clearfix" value="0"></progress>
    <textarea class="copy_area" style="opacity:0;height:0;position:absolute;left:-10000px;"></textarea>
</div>
<script>
    function refresh_files() {
        var pid = $('#form-post input[name="id"]').val();
        $.get('/v1/posts/' + pid + '/files')
        .done(function(data) {
            var files = $('.widget-files .list');
            files.empty();

            data.forEach(function(file) {
                var li = $('<li/>')
                    .css('overflow', 'hidden')
                    .append($('<span />').text(file));
                var btns = $('<span style="float:right;" />');
                if(/\.(jpg|gif|png|bmp)$/i.test(file)) {
                    btns.append('<button class="copy_as_md" title="复制为Markdown">复制</button>');
                }
                btns.append('<button class="delete">删除</button>');
                li.append(btns);
                files.append(li);
            });

            bind_copy_as_md();
            bind_delete();
        })
        .fail(function(x){
            alert(x.responseText);
        });
    }

    $('.widget-files .refresh').click(function(){
        refresh_files();
        return false;
    });

    $('.widget-files .files').on('change', function(e) {
        $('.widget-files .count').text(e.target.files.length + ' 个文件');
    });

    $('.widget-files .browse').click(function(){
        $('.widget-files .files').click();
        return false;
    });

    function bind_delete() {
        $('.widget-files .list .delete').click(function(){
            var li = $(this).parent().parent();
            var name = $(this).parent().prev().text();
            var pid = $('#form-post input[name="id"]').val();
            $.ajax({
                url: '/v1/posts/' + pid + '/files/' + encodeURI(name),
                type: 'DELETE',
                success: function() {
                    li.remove();
                },
                error: function() {
                    alert('删除失败。');
                }
            });
            return false;
        });
    }

    function bind_copy_as_md() {
        var ta = $('.widget-files .copy_area')[0];
        $('.widget-files .list .copy_as_md').click(function() {
            var name = $(this).parent().prev().text();
            var text = '![' + name + '](' + name + ')';
            console.log('Markdown: ' + text);
            ta.value = text;
            ta.focus();
            ta.select();
            try {
                if(!document.execCommand('copy')) {
                    throw -1;
                }
            } catch (e) {
                alert('复制失败。'+e);
            }
            return false;
        });
    }

    $('.widget-files .submit').click(function(){
        var files = $('.widget-files .files')[0].files;

        if(files.length <= 0) {
            alert('请先选择文件再上传。');
            return false;
        }
        
        var data = new FormData();

        // 待上传的文件列表
        for(var i = 0, n = files.length; i < n; i++) {
            var file = files[i];
            data.append('files[]', file);
        }

        // 进度条
        var progress = $('.widget-files .progress');
        progress.attr('value', 0);

        // 当前文章ID（新文章并没有ID，这里先临时使用下一篇文章ID）
        // 所以，不能同时编辑并发表新文章
        var pid = $('#form-post input[name="id"]').val();

        // https://stackoverflow.com/a/8758614/3628322
        $.ajax({
            // Your server script to process the upload
            url: '/v1/posts/' + pid + '/files',
            type: 'POST',

            // Form data
            data: data,

            // Tell jQuery not to process data or worry about content-type
            // You *must* include these options!
            cache: false,
            contentType: false,
            processData: false,

            // Custom XMLHttpRequest
            xhr: function() {
                var myXhr = $.ajaxSettings.xhr();
                if (myXhr.upload) {
                    // For handling the progress of the upload
                    myXhr.upload.addEventListener('progress', function(e) {
                        if (e.lengthComputable) {
                            progress.attr({
                                value: e.loaded,
                                max: e.total,
                            });
                        }
                    } , false);
                }
                return myXhr;
            },

            error: function(xhr, except) {
                console.warn(xhr,except);
                alert('ajax error:'+xhr.statusText);
                progress.attr('value', 0);
            },

            success: function() {
                $('.widget-files .files').val("");
                $('.widget-files .count').text("0 个文件");
                refresh_files();
                progress.attr('value', 0);
            },
        });

        return false;
    });
</script>
EOD;

    return compact('title', 'classname', 'types', 'content');
}

add_hook('post_widget', 'post_widget_files');

function post_widget_metas($p=null) {
    $metas = str_replace(['\\','\''], ['\\\\','\\\''], $p ? $p->metas_raw : '{}');
    $title = '自定义';
    $classname = 'metas';
    $content = <<< DOM
<label>类型：</label>
<select class="keys">
    <option>&lt;新建&gt;</option>
</select>
<span class="new">
    <input class="key" type="text" style="width: 100px;" />
    <button class="ok">添加</button>
</span>
<textarea class="content" style="margin-top: 10px; height: 200px; display: block;"></textarea>

<input type="hidden" name="metas" value="" />

<script>
(function() {
    var keys = $('.widget-metas .keys');
    var metas = JSON.parse('$metas');
    var newf = $('.widget-metas .new');
    var content = $('.widget-metas .content');

    $('.widget-metas input[name=metas]').val('$metas');

    for(var key in metas) {
        keys.append($('<option>', {value: key, text: key}));
    }

    var prev_key = '';

    function save_prev() {
        if(prev_key) {
            metas[prev_key] = content.val();
        }
    }

    content.on('blur', function() {
        save_prev();
        $('.widget-metas input[name=metas]').val(JSON.stringify(metas));
    });

    keys.on('change', function() {
        var i = this.selectedIndex;


        if(i == 0) {
            newf.css('visibility', 'visible');
            prev_key = '';
            content.val('');
        }
        else {
            newf.css('visibility', 'hidden');
            prev_key = keys[0].options[i].value;
            content.val(metas[prev_key]);
        }

    });

    $('.widget-metas .new .ok').on('click', function() {
        var key = $('.widget-metas .new .key').val().trim();
        if(key == '' || metas.hasOwnProperty(key)) {
            alert('为空或已经存在。');
            return false;
        }

        keys.append($('<option>', {value: key, text: key}));
        keys.val(key);
        prev_key = key;
        content.focus();
        newf.css('visibility', 'hidden');

        return false;
    });
})();
</script>
DOM;

    return compact('title', 'content', 'classname');
}

add_hook('post_widget', 'post_widget_metas');

function post_widget_tax_add(&$taxes, $tax=1) {
    $s = '';
    foreach($taxes as $t) {
        $s .= '<li style="margin-bottom: 4px;"><label><input type="radio" style="margin-right: 6px" name="taxonomy" value="'.$t->id.'"'.
            ($tax==$t->id?' checked="checked"':'').'/>'.htmlspecialchars($t->name)."</label>\n";
        if(isset($t->children) && count($t->children)) {
            $s .= '<ul class="children" style="margin-left: 14px;">';
            $s .= post_widget_tax_add($t->children, $tax);
            $s .= "</ul>\n";
        }
        $s .= '</li>'."\n";
    }
    return $s;
}

function post_widget_tax($p=null) {
    $taxes = get_tax_tree();
    $content = '<ul>'.post_widget_tax_add($taxes, ($p ? $p->taxonomy : 1)).'</ul>';

    return [
        'title'		=> '分类',
        'content'	=> $content,
        'classname'	=> 'category',
        'types' => 'post',
        ];
}

add_hook('post_widget', 'post_widget_tax');

function post_widget_page_parents($p=null) {
    global $tbpost;
    if($p) {
        $v = $tbpost->get_the_parents_string($p->id);
        if($v) {
            $v = substr($v, 1);
            $v = implode(',', explode('/',$v));
        }
    }
    $content = '<input type="text" name="page_parents" value="'.($p ? $v : '').'" />';

    return [
        'title' => '父页面',
        'content' => $content,
        'types' => 'page',
    ];
}

add_hook('post_widget', 'post_widget_page_parents');

function post_widget_slug($p=null) {
    return [
        'title' => '别名',
        'types' => 'page,post',
        'content' => '<input type="text" name="slug" value="'.($p ? htmlspecialchars($p->slug) : '').'" />',
        ];
}

add_hook('post_widget', 'post_widget_slug');

function post_widget_date($p=null) {
    global $tbdate;

    $title = '日期';
    $content = '<input type="text" name="date" value="'.($p ? $p->date : '').'"/><br>'
        .'<input type="text" name="modified" value="'.($p ? $p->modified : '').'" />';

    return compact('title', 'content');
}

add_hook('post_widget', 'post_widget_date');

function post_admin_head() {
    $post = $GLOBALS['__p__'] ?? null;

    echo '<title>', $post ? '【编辑文章】'.htmlspecialchars($post->title) : '新文章','</title>';

?>

<script src="scripts/marked.js"></script>

<script>
    marked.setOptions({
        sanitize: false,
    });
</script>

<style>
    .sidebar {

    }

    .sidebar input[type="text"] {
        padding: 4px;
    }

    .sidebar .widget {
        background-color: white;
        border: 1px solid #ccc;
        margin-bottom: 20px;
    }

    .sidebar .widget h3 {
        padding: 4px 6px;
        border-bottom: 1px solid #ccc;
    }

    .sidebar .widget-content {
        padding: 10px;
    }

    .sidebar .widget ul {
        list-style: none;
    }

    .post-area {
        margin-bottom: 3em;
    }

    .widget-category .widget-content {
        max-height: 200px;
        overflow: auto;
    }

    .widget-content input[type=text], .widget-content textarea {
        padding: 4px;
        width: 100%;
        box-sizing: border-box;
    }

    #source {
        max-height: 2000px;
        height: 70vh;
        min-height: 300px;
        width: 100%;
        padding: 4px;
        box-sizing: border-box;
    }

#form-post {
    display: flex;
}

.post {
    flex: 1;
}
.sidebar-right {
    flex: 1;
}

/* TODO 根据主题修改 */
@media screen and (min-width: 851px) {
    .sidebar-right {
        width: 280px;
        max-width: 280px;
        min-width: 280px;
    }
    .post {
        margin-right: 1em;
    }
}

/* TODO 根据主题修改 */
@media screen and (max-width: 850px) {
    #form-post {
        flex-direction: column;
    }
}
</style>
<?php
    apply_hooks('admin:post:head');
}

add_hook('admin_head', 'post_admin_head');

function post_admin_footer() { ?>
    <script type="text/javascript">
        $('.widget h3').click(function(e) {
            var div = e.target.nextElementSibling;
            $(div).toggle();
        });
    </script>
<?php
    apply_hooks('admin:post:foot');
}

add_hook('admin_footer', 'post_admin_footer');

function new_post_html($p=null){
    global $tbpost;

    // 先生成所有的挂件对象
    // 因为分布在不同地方（hook对象无法保存这些分布）
    $widgets = [];

    $widget_objs = get_hooks('post_widget');
    foreach($widget_objs as $wo) {
        $fn = $wo->func;
        $w = (object)$fn($p);
        $w->classname = $w->classname ?? 'widget';

        $dom = <<< DOM
<div class="widget widget-$w->classname">
    <h3>$w->title</h3>
    <div class="widget-content">
        $w->content
    </div>
</div> 
DOM;
        $widget = new stdClass;
        $widget->dom = $dom;
        $widget->pos = $w->position ?? 'right';
        $widget->types = $w->types ?? 'post,page';

        $widgets[] = $widget;
    }

    $type = $p ? $p->type : ($_GET['type'] ?? '');
    if(!in_array($type, ['post','page']))
        $type = 'post';

?><div id="admin-post">
    <form method="POST" id="form-post">
        <div class="post">
            <div class="post-area">
                <div style="margin-bottom: 1em;">
                    <h2>标题</h2>
                    <div>
                    <input style="padding: 8px; width: 100%; box-sizing: border-box;" type="text" name="title" value="<?php
                        if($p) {
                            echo htmlspecialchars($p->title);
                        }
                    ?>" />
                    </div>
                </div>
                <?php if($p) {
                    $link = the_link($p);
                ?>
                <div class="permanlink" style="margin-bottom: 1em;">
                    <span>固定链接：</span>
                    <a id="permalink" href="<?php echo $link; ?>"><?php echo $link; ?></a>
                    <script type="text/javascript">
                        var new_window = null;
                        $('#permalink').click(function() {
                            if(!new_window || new_window.closed) {
                                new_window = window.open($('#permalink').prop('href'));
                            } else {
                                new_window.location.href = $('#permalink').prop('href');
                            }
                            return false;
                        });
                    </script>
                </div>
                <?php } else {
                    $next_id = $tbpost->the_next_id();
                ?>
                <div class="permalink_id" style="margin-bottom: 1em;">
                    <span>文章ID：</span>
                    <span><?php echo $next_id; ?></span>
                </div>
                <?php } ?>
                <div>
                    <h2>内容</h2>
                    <div class="textarea-wrap">
                        <textarea id="source" name="source" wrap="off"><?php
                            if($p) {
                                echo htmlspecialchars($p->source ? $p->source : $p->content);
                            }
                        ?></textarea>
                    </div>
                </div>
                <div>
                    <input type="hidden" name="do" value="<?php echo $p ? 'update' : 'new'; ?>" />
                    <input type="hidden" name="type" value="<?php echo $p ? $p->type : $type; ?>" />
                    <input type="hidden" name="id" value="<?php echo $p ? $p->id : $next_id; ?>" />
                </div>
            </div><!-- post-area -->
            <div class="sidebar sidebar-left">
                <?php foreach($widgets as &$widget) {
                    if(in_array($type, explode(',',$widget->types)) && $widget->pos == 'left') {
                        echo $widget->dom;
                    }
                } ?>
            </div>
        </div><!-- post -->
        <div class="sidebar sidebar-right">
            <div class="widget widget-post">
                <h3>发表</h3>
                <div class="widget-content">
                    <input type="submit" value="发表" />
                    <select name="status">
                        <option value="public"<?php if($p && $p->status == 'public') echo ' selected'; ?>>公开</option>
                        <option value="draft"<?php if($p && $p->status == 'draft') echo ' selected'; ?>>草稿</option>
                    </select>
                    <select name="source_type">
                        <?php if($p && $p->source_type === '') $p->source_type = 'html'; ?>
                        <option value="markdown"<?php if($p && $p->source_type == 'markdown') echo ' selected'; ?>>Markdown</option>
                        <option value="html"<?php if($p && $p->source_type == 'html') echo ' selected'; ?>>HTML</option>
                    </select>
                </div>
            </div>
            <?php foreach($widgets as &$widget) {
                if(in_array($type, explode(',',$widget->types)) && $widget->pos == 'right') {
                    echo $widget->dom;
                }
            } ?>
        </div><!-- sidebar right -->
        <script>
            // TODO 临时代码，在切换源的类型时切换编辑器语法
            $('select[name="source_type"]').change(function() {
                console.log('源类型切换为：', this.value);
                if(typeof codemirror == 'object') {
                    var mode = '';
                    var value = this.value;

                    if(value == 'markdown')
                        mode = 'markdown';
                    else if(value == 'html')
                        mode = 'htmlmixed';

                    codemirror.setOption('mode', mode);
                }
                else {
                    console.warn('codemirror != object, cannot apply syntax.');
                }
            });

            setTimeout(function(){
                $('select[name="source_type"]').change();
            },0);

            $('#form-post').submit(function() {
                var form = document.getElementById('form-post');
                var value = function(name) {
                    return form.elements[name].value;
                };
                var post = {
                    id: value('do') === 'new' ? 0 : +value('id'),
                    title: value('title'),
                    source: value('source'),
                    source_type: value('source_type'),
                    type: value('type'),
                    tags: value('tags').replace('，',',').split(','),
                    metas: value('metas'),
                    category: +value('taxonomy'),
                    slug: value('slug'),
                    date: value('date'),
                    modified: value('modified'),
                };
                if (post.tags.length === 1 && post.tags[0] === "") {
                    post.tags = [];
                }
                console.log(post);
                var url = post.id === 0 ? '/v1/posts' : '/v1/posts/' + post.id;
                $.ajax(url,{
                    type: 'POST',
                    data: JSON.stringify(post),
                    contentType: 'application/json',
                    success: function(data) {
                        console.log(data);
                        location.href = '/admin/post.php?do=edit&id='+data;
                    },
                    error: function(xhr) {
                        alert(xhr.responseText);
                    },
                });
                return false;
            });
        </script>
    </form>
</div><!-- admin-post -->
<?php } 

$do = $_GET['do'] ?? '';
if(!$do) {
    admin_header();
    new_post_html();
    admin_footer();
    die(0);
} else if($do === 'edit') {
    $id = intval($_GET['id']);
    $post = $tbpost->query_by_id($id,'');
    if($post === false || empty($post)){
        tb_die(200, '没有这篇文章！');
    }

    // 输出编辑内容之前过滤
    if(isset($post[0]->content))
        $post[0]->content = apply_hooks('edit_the_content', $post[0]->content, $post[0]->id);

    // 罪过，使用全局变量了
    $GLOBALS['__p__'] = $post[0];

    admin_header();
    new_post_html($post[0]);
    admin_footer();
    die(0);
}


die(0);

/* GET */ else :

function post_die_json($arg) {
    header('HTTP/1.1 200 OK');
    header('Content-Type: application/json');

    echo json_encode($arg, JSON_UNESCAPED_UNICODE);
    die(0);
}

require_once('login-auth.php');

if(!login_auth()) {
    post_die_json([
        'errno' => 'unauthorized',
        'error' => '需要登录后才能进行该操作！',
        ]);
}

require_once('load.php');

function post_update() {
    global $tbpost;

    $r = $tbpost->update($_POST);
    if(!$r) {
        post_die_json([
            'errno' => 'error',
            'error' => $tbpost->error
            ]);
    }

    $id = (int)$_POST['id'];

    header('HTTP/1.1 302 Updated');
    header('Location: /admin/post.php?do=edit&id='.$id);
    die(0);
}

$do = $_POST['do'];

if($do == 'update') {
    post_update();
}

die(0);

endif; // POST
