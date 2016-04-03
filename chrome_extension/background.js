// Set up context menu at install time.
chrome.runtime.onInstalled.addListener(function() {
    var context = "link";
    var title = "Send to Server";
    var id = chrome.contextMenus.create({"title": title, "contexts":[context],
                                         "id": "context" + context});
});

// add click event
chrome.contextMenus.onClicked.addListener(onClickHandler);

// The onClicked callback function.
function onClickHandler(info, tab) {
    test();
    // var sText = info.selectionText;
    // var url = "https://www.google.com/search?q=" + encodeURIComponent(sText);
    // window.open(url, '_blank');
};

command = function(){
    fetch('http://localhost:17901/api/', {
        method: 'POST',
        headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        apiKey: "23462",
        action: "AHHH"
    })}).then(function(response) {
        return response.json()
    }).then(function(json) {
        console.log('parsed json', json)
    }).catch(function(ex) {
        console.log('parsing failed', ex)
    })
}
