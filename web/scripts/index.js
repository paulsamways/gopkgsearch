$(function() {
  var submitTimer = null;
  
  var $results = $("#results");
  var $input = $("#query"); 

  var renderResult = function(r) {
    return "<li class=\"" + r.Kind + "\"><h2>" + r.Name + "</h2></li>";
  };

  var submitQuery = function(q) {
   
    var callback = function(data, status, xhr) {

      $("li", $results).remove();

      if (data && data.length) {
        for (var i = 0; i < data.length; i++) {
          var r = data[i];

          $(renderResult(r)).appendTo($results);
        }
      }
    };

    if (q == "") {
      callback();
      return;
    }
    
    $.post("/query?q=" + q, { }, callback, "json");
  };

  $input.keyup(function(e) {
    var self = this;

    if (submitTimer != null) {
      clearTimeout(submitTimer);
    }

    submitTimer = setTimeout(function() {
      submitQuery($(self).val());
    }, 250);

  });
});
