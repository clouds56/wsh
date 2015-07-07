jQuery(function($, undefined) {
  $('#main-term').terminal(function(command, term) {
    if (command !== '') {
      var result = window.eval(command);
      if (result != undefined) {
        term.echo(String(result));
      }
    }
  }, {
    greetings: 'Javascript Interpreter',
    name: 'main-term',
    prompt: 'js> '});
});
