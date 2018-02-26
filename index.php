<?php
$use_sts = true;
// Force HTTPS (Strict Transport Security)
if ($use_sts && isset($_SERVER['HTTPS']) && $_SERVER['HTTPS'] != 'off') {
	header('Strict-Transport-Security: max-age=31536000');
} elseif ($use_sts) {
	header('Location: https://'.$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'], true, 301);
	die();
}
?>
<!DOCTYPE html>
<html lang="en" aria-label="dchr.host">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name=viewport content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
  <title>dechristopher</title>
  <style>
    html, body, pre {
      margin: 0px;
      height: 100%;
      text-align: center;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
  </style>
</head>
<body>
  <pre id="me"><b>&#109;&#101;&#064;&#100;&#99;&#104;&#114;&#46;&#104;&#111;&#115;&#116;</b>
github.com/${me}
keybase.io/${me}
dev.to/${me}
</pre>
<script async src="me.js"></script>
<script>const me = 'dechristopher';</script>
</body>