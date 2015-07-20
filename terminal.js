jQuery(function($, undefined) {
  $('#main-term').terminal("/repl/", {
    greetings: 'Javascript Interpreter',
    name: 'main-term',
    prompt: 'js> '});
});
