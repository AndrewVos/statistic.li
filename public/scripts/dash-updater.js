$(document).ready(function() {
  var clientId = $("#clientId").val();
  setInterval(function() {
    $.getJSON("/client/" + clientId + "/views", function(data) {
      $("h1").text(data.views);
    });
  }, 5000);
});
