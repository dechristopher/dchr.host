<?php
$use_sts = true;
// Force HTTPS (Strict Transport Security)
if ($use_sts && isset($_SERVER['HTTPS']) && $_SERVER['HTTPS'] != 'off') {
	header('Strict-Transport-Security: max-age=31536000');
}
elseif ($use_sts) {
	header('Location: https://'.$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'], true, 301);
	die();
}
header("Cache-Control: max-age=2592000");
?>
<!DOCTYPE html>
<html lang="en" aria-label="dchr.host">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name=viewport content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
  <title>dchr.host</title>
  <link rel="apple-touch-icon-precomposed" sizes="57x57" href="res/apple-touch-icon-57x57.png" />
  <link rel="apple-touch-icon-precomposed" sizes="114x114" href="res/apple-touch-icon-114x114.png" />
  <link rel="apple-touch-icon-precomposed" sizes="72x72" href="res/apple-touch-icon-72x72.png" />
  <link rel="apple-touch-icon-precomposed" sizes="144x144" href="res/apple-touch-icon-144x144.png" />
  <link rel="apple-touch-icon-precomposed" sizes="60x60" href="res/apple-touch-icon-60x60.png" />
  <link rel="apple-touch-icon-precomposed" sizes="120x120" href="res/apple-touch-icon-120x120.png" />
  <link rel="apple-touch-icon-precomposed" sizes="76x76" href="res/apple-touch-icon-76x76.png" />
  <link rel="apple-touch-icon-precomposed" sizes="152x152" href="res/apple-touch-icon-152x152.png" />
  <link rel="icon" type="image/png" href="res/favicon-196x196.png" sizes="196x196" />
  <link rel="icon" type="image/png" href="res/favicon-96x96.png" sizes="96x96" />
  <link rel="icon" type="image/png" href="res/favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="res/favicon-16x16.png" sizes="16x16" />
  <link rel="icon" type="image/png" href="res/favicon-128.png" sizes="128x128" />
  <meta name="application-name" content="dchr.host"/>
  <meta name="msapplication-TileColor" content="#FFFFFF" />
  <meta name="msapplication-TileImage" content="res/mstile-144x144.png" />
  <meta name="msapplication-square70x70logo" content="res/mstile-70x70.png" />
  <meta name="msapplication-square150x150logo" content="res/mstile-150x150.png" />
  <meta name="msapplication-wide310x150logo" content="res/mstile-310x150.png" />
  <meta name="msapplication-square310x310logo" content="res/mstile-310x310.png" />
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
  <!--<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/3.5.2/animate.min.css">-->
  <link rel="stylesheet" type="text/css" href="">
</head>
<body>
  <pre id="me"><b class="to">&#109;&#101;&#064;&#100;&#99;&#104;&#114;&#46;&#104;&#111;&#115;&#116;</b>
github.com/dechristopher
keybase.io/dechristopher
dev.to/dechristopher
</pre>
<script src="https://code.jquery.com/jquery-3.3.1.min.js"
  integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8="
  crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/lettering.js/0.7.0/jquery.lettering.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/textillate/0.4.0/jquery.textillate.min.js"></script>
<script src="me.js"></script>
<noscript id="deferred-styles">
      <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/3.5.2/animate.min.css"/>
    </noscript>
    <script>
      var loadDeferredStyles = function() {
        var addStylesNode = document.getElementById("deferred-styles");
        var replacement = document.createElement("div");
        replacement.innerHTML = addStylesNode.textContent;
        document.body.appendChild(replacement)
        addStylesNode.parentElement.removeChild(addStylesNode);
      };
      var raf = window.requestAnimationFrame || window.mozRequestAnimationFrame ||
          window.webkitRequestAnimationFrame || window.msRequestAnimationFrame;
      if (raf) raf(function() { window.setTimeout(loadDeferredStyles, 0); });
      else window.addEventListener('load', loadDeferredStyles);
    </script>
</body>