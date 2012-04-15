$(function() {
  var submitTimer = null;
  
  var $results = $("#results");
  var $input = $("#query");

  var escapeHtml = function(raw) {
    if (raw === undefined || raw === null)
      return ""

    var r = raw.replace(/&/g,'&amp;')
              .replace(/</g,'&lt;')
              .replace(/>/g,'&gt;');

    return r;
  };

  var renderResult = function(r) {
    return  "<li class=\"" + r.Kind + "\">" +
              "<h2 title=\"" + r.FilePath + "\"><span class=\"package\">" + r.Package + "</span>." + r.Name + "</h2>" +
              "<div style=\"display: none;\" class=\"doc\">" + escapeHtml(r.Doc) + "</div>" +
              "<pre style=\"display: none;\" class=\"source\">" + escapeHtml(r.Source) + "</pre>" +
            "</li>";
  };

  var submitQuery = function(q) {
   
    var callback = function(data, status, xhr) {

      $("li", $results).remove();

      if (data && data.length) {
        for (var i = 0; i < data.length; i++) {
          var r = data[i];

          var $d = $(renderResult(r))
            .appendTo($results);

            $d.click(function(e) {
              $("div.doc, pre.source", this).toggle();
            });

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
