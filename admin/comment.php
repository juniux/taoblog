<?php

if($_SERVER['REQUEST_METHOD'] == 'GET') :




else :

require_once('load.php');

function  cmt_header_json() {
	header('HTTP/1.1 200 Not OK');
	header('Content-Type: application/json');
}

function error($msg) {
	cmt_header_json();
	echo json_encode([
		'errno' => 'error',
		'error' => $msg
		]);
	die(-1);
}

function cmt_filter_cmt(&$cmts) {
	$flts = ['email', 'url', 'ip', 'agent', 'status'];
	for($i=0; $i<count($cmts); $i++) {
		foreach($flts as $f) {
			unset($cmts[$i]->$f);
		}

		if(isset($cmts[$i]->children)) {
			for($x=0; $x<count($cmts[$i]->children); $x++) {
				foreach($flts as $f) {
					unset($cmts[$i]->children[$x]->$f);
				}
			}
		}
	}

	return $cmts;
}

function cmt_get_cmt() {
	global $tbcmts;
	cmt_header_json();

	$cmts = cmt_filter_cmt($tbcmts->get($_POST));

	echo json_encode([
		'errno'		=> 'success',
		'cmts'		=> $cmts,
		]);
	die(0);
}

function cmt_post_cmt() {
	global $tbcmts;
	global $tbdb;

	$ret_cmt = (int)(isset($_POST['return_cmt']) ? $_POST['return_cmt'] : '');
	
	$r = $tbcmts->insert($_POST);
	if(!$r) {
		header('HTTP/1.1 200 Not OK');
		header('Content-Type: application/json');

		echo json_encode([
			'errno'	=> 'error',
			'error' => $tbcmts->error
			]);
		die(-1);
	}

	if($ret_cmt) {
		header('HTTP/1.1 200 OK');
		header('Content-Type: application/json');

		$c = ['id'=>$r];

		echo json_encode([
			'errno'	=> 'success',
			'cmt'	=> cmt_filter_cmt($tbcmts->get($c))[0],
			]);
		die(0);
	} else {
		header('HTTP/1.1 200 OK');
		header('Content-Type: application/json');

		echo json_encode([
			'errno'	=> 'success',
			'id'	=> $r,
			]);
		die(0);
	}
}

function cmt_get_count() {
	global $tbcmts;

	echo $tbcmts->get_count((int)$_POST['post_id']);

	die(0);
}

$do = isset($_POST['do']) ? $_POST['do'] : '';

if($do == 'get-cmt') cmt_get_cmt();
if($do == 'post-cmt') cmt_post_cmt();
if($do == 'get-count') cmt_get_count();

endif;

