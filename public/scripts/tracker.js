!(function(){
  var referer = document.referrer;
  var page = window.location.href;
	var clientId = document.location.host;
	var trackerUrl = "https://statistic.li/client/" + encodeURIComponent(clientId) + "/tracker.gif?referer=" + encodeURIComponent(referer) + "&page=" + encodeURIComponent(page);
	var img     = new Image();
	img.src = trackerUrl;
})();
