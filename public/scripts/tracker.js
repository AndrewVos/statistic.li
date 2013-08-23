!(function(){
  var referer = document.referrer;
  // var page = window.location.href;
	var trackerUrl = "https://statistic.li/client/" + clientId + "/tracker.gif?referer=" + encodeURIComponent(referer);
	var img     = new Image();
	img.src = trackerUrl;
})();
