$(function() {
  var _KEYUP = 38;
  var _KEYDOWN = 40;
  var _KEYESC = 27;
  var _KEYENTER = 13;
  var _KEYSPACE = 32;

  var submitTimer = null;
  var prevQ = null;
  
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
              "<div style=\"display: none;\" class=\"path\">" + escapeHtml(r.FilePath) + "</div>" +
            "</li>";
  };

  var deselect = function() {
    $("li.selected", $results).removeClass("selected");
  };
  var select = function($li) {
    if ($li === undefined) {
      $li = $("li:first", $results);
    }

    $li.addClass("selected");

    if ($li.length > 0) {
      $li[0].scrollIntoView(false);
    }
  };
  var selectNext = function() {
    var $s = $("li.selected", $results).removeClass("selected");
    
    var $n = $s.next("li");

    if ($n.length === 0) {
      select();
    } else {
      select($n);
    }
  };
  var selectPrev = function() {
    var $s = $("li.selected", $results).removeClass("selected");
    var $p = $s.prev();

    if ($p.length === 0) {
      $p = $("li:last", $results);
    }

    select($p);
  };

  var reset = function() {
    deselect();
    $("li", $results).remove();
    $input.val("");
    $input.focus();
  };

  var expand = function($li) {
    if ($li === undefined) {
      $li = $("li.selected", $results);
    }

    $("div.doc, pre.source, div.path", $li).toggle();
  };

  var submitQuery = function(q) {
   
    var callback = function(data, status, xhr) {

      $("li", $results).remove();

      if (data && data.length) {
        for (var i = 0; i < data.length; i++) {
          var r = data[i];

          var $d = $(renderResult(r))
            .appendTo($results);

          $d.click(function(e) { deselect(); expand($(this)); });
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

    e.stopPropagation();

    if (submitTimer != null) {
      clearTimeout(submitTimer);
    }

    switch (e.which) {
      case _KEYDOWN:
        e.preventDefault();
        $(self).blur();
        select();
        return;
      case _KEYESC:
        e.preventDefault();
        reset();
        return;
    }

    submitTimer = setTimeout(function() {
      var q = $(self).val();

      if (prevQ === q) {
        return;
      }

      prevQ = q;
      submitQuery($(self).val());
    }, 50);
  });
  
  $input.keypress(function(e) {
    e.stopPropagation();
  });


  $(window).keyup(function(e) {
    switch (e.which) {
      case _KEYDOWN:
        selectNext();
        e.preventDefault();
        return;
      case _KEYUP:
        selectPrev();
        e.preventDefault();
        return;
      case _KEYESC:
        reset();
        e.preventDefault();
        return;
      case _KEYENTER:
      case _KEYSPACE:
        expand();
        e.preventDefault();
        return;
    }
  });
  $(window).keydown(function(e) {
    if (e.which === _KEYDOWN || e.which === _KEYUP) {
      e.preventDefault();
    }
  });
});
