String.prototype.format = function() {
    var theString = this;
    for (var i = 0; i < arguments.length; i++) {
        var regEx = new RegExp("\\{" + i + "\\}", "gm");
        theString = theString.replace(regEx, arguments[i]);
    }
    return theString;
}

jQuery(function($, undefined) {
  id = 1;
  $('#term-now').keydown(function(e) {
    //console.log(e);
    if (e.shiftKey == true && e.keyCode == 13) {
      e.preventDefault();
      var item = $('#term-now');
      item.before("<div class=\"term-input\">{0}</div>".format(item[0].innerText.replace(/\n/g, "<br>")));
      $.post("/repl/", {"jsonrpc":"2.0", "id":1, "method":"cmd", "params":item[0].innerText}, function(d) {
        r = JSON.parse(d)
        item.before("<div class=\"term-output\">{0}</div>".format(r.result.replace(/\n/g, "<br>")));
      })
      item.empty();
      return false;
    }
    return true;
  })
});
