function onResponse(response) {
    var divNotification = document.querySelector("#notice");
    var divNotificationMessage = document.querySelector("#notice-message");
    var className = (response.status !== 400) ? "alert-success" : "alert-danger";
    divNotificationMessage.innerHTML = response.data;

    divNotification.classList.add(className);

    document.getElementById("notice").style.display = "block";
    setTimeout(function () {
        document.getElementById("notice").style.display = "none";
        divNotification.classList.remove(className);
    },3000);
}

function senddata() {
    var file = document.getElementById("inputFile").files[0];
    var tool = document.getElementById("ChooseToolReport").value

    if (!file) {
        return;
    }
    upload(file,tool);

    function upload(file,tool) {
        var formData = new FormData();
        formData.append('file', file);
        formData.append('tool', tool);

        post('/apps/uploadReport', formData)
            .then(onResponse)
            .catch(onResponse)
    }
}(document, axios);

function addapp() {
    var name = document.getElementById("name").value;
    var url = document.getElementById("url").value;
    var language = document.getElementById("language").value
    var framework = document.getElementById("framework").value

    uploadAppData(name, url, language, framework);

    name.reset();
    url.reset();
    language.reset();
    framework.reset();

    function uploadAppData(name,url,language,framework) {
        var formData = new FormData();
        formData.append('name', name);
        formData.append('url', url);
        formData.append('language', language);
        formData.append('framework', framework);

        post('/apps/addApp', formData)
            .then(onResponse)
            .catch(onResponse)
    }
}(document, axios);


function startscan(idApp) {
    var formData = new FormData();
    if (document.getElementById("zaproxy").checked === true) {
        formData.append('zaproxy', zaproxy);
    }
    if (document.getElementById("arachni").checked === true) {
        formData.append('arachni', arachni);
    }

    formData.append('idApp', idApp);

    post('/apps/startScan', formData)
        .then(onResponse)
        .catch(onResponse)

    document.getElementById("zaproxy").checked = false;
    document.getElementById("arachni").checked = false;

}(document, axios);
