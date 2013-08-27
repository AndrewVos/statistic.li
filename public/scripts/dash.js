var uniques = new Array();

function updateUniques() {
  var clientId = $("#clientId").val();
  $.getJSON("/client/" + clientId + "/uniques", function(data) {
    uniques.push(data.Count)
    if (uniques.length == 301) { uniques.splice(0, 1); }
    $('#uniques').text(data.Count + " users online")
    $('#uniques-graph').sparkline(uniques, {width: "100%",height: "100%", type:"line", tooltipSuffix: " users online"});
  });
}

function updateTopPages() {
  var clientId = $("#clientId").val();
  $.getJSON("/client/" + clientId + "/pages", function(data) {
    var container = $("<div id='pages-container'></div>");
    $.each(data, function(i, d) {
      var row = $("<div class='row'></div>");
      row.append("<div class='col-md-10'>" + d.Page + "</div>");
      row.append("<div class='col-md-2'>" + d.Count + "</div>");
      container.append(row);
    });
    $("#pages-container").replaceWith(container);
  });
}

function updateTopReferers() {
  var clientId = $("#clientId").val();
  $.getJSON("/client/" + clientId + "/referers", function(data) {
    var container = $("<div id='referers-container'></div>");
    $.each(data, function(i, d) {
      var row = $("<div class='row'></div>");
      row.append("<div class='col-md-10'>" + d.Referer + "</div>");
      row.append("<div class='col-md-2'>" + d.Count + "</div>");
      container.append(row);
    });
    $("#referers-container").replaceWith(container);
  });
}

function update() {
  updateUniques();
  updateTopPages();
  updateTopReferers();
}

$(document).ready(function() {
  update();
  setInterval(updateUniques, 1000);
  setInterval(function() {
    updateTopPages();
    updateTopReferers();
  }, 5000);

  $("#generate-traffic").click(function() {
    update();
    var clientId = $("#clientId").val();
    $.post("/client/" + clientId + "/generate");
  });
});
