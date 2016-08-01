var page = require('webpage').create();
page.viewportSize = { width: 1024, height: 768 };

var system = require('system');
var args = system.args;

var url = args[1];
var outputFileName = args[2];

/*
 *page.onError = function(msg, trace) {
 *  console.log(msg);
 *  console.log(trace);
 *});
 */

page.open(url, function() {
  page.render(outputFileName,  {format: 'png', quality: '100'});
  phantom.exit();
});
