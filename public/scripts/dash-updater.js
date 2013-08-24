$(document).ready(function() {
  var clientId = $("#clientId").val();
  setInterval(function() {
    $.getJSON("/client/" + clientId + "/uniques", function(data) {
      $("h1").text(data.uniques);
    });
  }, 5000);
});
